// Package mockserver is the core of this program and it will have the API
// to record binary data from the UDP server and then be able to mock it
package mockserver

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

type ViewData struct {
	Data     []byte
	SizeRead int
}

// renderToTerminal will render the data for the users viewing pleasure
func renderToTerminal(ctx context.Context, stream chan ViewData, fp string) {
	fileInfo, err := os.Stat(fp)
	if err != nil {
		// NOTE: learn how to handle this error
		// return fmt.Errorf("error stating file: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case data := <-stream:
			percent := int(float64(data.SizeRead) / float64(fileInfo.Size()) * 100)
			fmt.Printf("\rReplayed: %d%%", percent)
		}
	}
}

// writeToUDPSocket will write the telemetry data to the UDP socket
func writeToUDPSocket(ctx context.Context, stream chan []byte, socket *UDPTransport) {
	for {
		select {
		case <-ctx.Done():
			return
		case data := <-stream:
			_, err := socket.Send(data)
			if err != nil {
				continue
				// NOTE: log the error somewhere maybe
				// return err
			}
		}
	}
}

// Replay replays a given file <fp> in a UDP server <addr>:<port>
func Replay(ctx context.Context, address string, port int, loop bool, fp string) error {
	bin, err := os.Open(fp)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}

	reader := NewGobReader(bin)
	udp, err := NewUDPTransport(address, port)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(time.Second / 60)
	defer ticker.Stop()

	uiChannel := make(chan ViewData, 1)
	socketChannel := make(chan []byte, 1)

	go renderToTerminal(ctx, uiChannel, fp)
	go writeToUDPSocket(ctx, socketChannel, udp)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-ticker.C:
			data, err := reader.Next()
			if err == io.EOF {
				if !loop {
					return udp.Close()
				}

				err = reader.Reset()
				if err != nil {
					return err
				}

				continue
			}

			if err != nil {
				return err
			}

			// Send the data to the view
			select {
			case uiChannel <- ViewData{Data: data, SizeRead: reader.TotalRead}:
				// Sent the data
			default:
				// Dropped the frame!
			}

			// Send the data to the UDP socket
			select {
			case socketChannel <- data:
				// Sent the data
			default:
				// Dropped the frame!
			}

		}
	}
}
