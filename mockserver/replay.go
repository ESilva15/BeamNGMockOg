// Package mockserver is the core of this program and it will have the API
// to record binary data from the UDP server and then be able to mock it
package mockserver

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
	"unsafe"

	bngsdk "github.com/ESilva15/gobngsdk"
)

type ViewData struct {
	Data     []byte
	SizeRead int64
}

// Replayer does the replaying
// Should we make a "player" struct that can record and replay?
type Replayer struct {
	DataSourcePath string
	Socket         *UDPTransport

	// Streams
	dataViewCh chan ViewData
	socketCh   chan []byte

	// Mut
	mut  sync.RWMutex
	data []byte

	// View
	viewData ViewData
}

func NewReplayer(address string, port int, fp string) (*Replayer, error) {
	udp, err := NewUDPTransport(address, port)
	if err != nil {
		return nil, err
	}

	replayer := &Replayer{
		DataSourcePath: fp,
		Socket:         udp,
		data:           make([]byte, unsafe.Sizeof(bngsdk.Outgauge{})),
		viewData:       ViewData{},
		dataViewCh:     make(chan ViewData, 1),
		socketCh:       make(chan []byte, 1),
	}

	return replayer, nil
}

// renderToTerminal will render the data for the users viewing pleasure
func (r *Replayer) renderToTerminal(ctx context.Context) {
	fileInfo, err := os.Stat(r.DataSourcePath)
	if err != nil {
		// NOTE: learn how to handle this error
		// return fmt.Errorf("error stating file: %v", err)
	}

	// NOTE: temporary until I make a better view
	lastPercent := -1

	for {
		select {
		case <-ctx.Done():
			return
		case data := <-r.dataViewCh:
			percent := int(float64(data.SizeRead) / float64(fileInfo.Size()) * 100)
			if percent != lastPercent {
				fmt.Printf("\rReplayed: %d%%", percent)
				lastPercent = percent
			}
		}
	}
}

// writeToUDPSocket will write the telemetry data to the UDP socket
func (r *Replayer) writeToUDPSocket(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case data := <-r.socketCh:
			r.mut.RLock()
			_, err := r.Socket.Send(data)
			r.mut.RUnlock()
			if err != nil {
				panic(fmt.Sprintf("error writing buffer to socket: %+v", err))
				continue
				// NOTE: log the error somewhere maybe
				// return err
			}
		}
	}
}

// Replay replays a given file <fp> in a UDP server <addr>:<port>
func (r *Replayer) Replay(ctx context.Context, loop bool) error {
	bin, err := os.Open(r.DataSourcePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}

	reader := NewGobReader(bin)

	go r.renderToTerminal(ctx)
	go r.writeToUDPSocket(ctx)

	ticker := time.NewTicker(time.Second / 60)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-ticker.C:
			r.mut.Lock()
			err := reader.Next(r.data)
			r.mut.Unlock()

			if err == io.EOF {
				if !loop {
					return r.Socket.Close()
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

			r.viewData.Data = r.data
			r.viewData.SizeRead = reader.TotalRead

			// Send the data to the view
			select {
			case r.dataViewCh <- r.viewData:
				// Sent the data
			default:
				// Dropped the frame!
			}

			// Send the data to the UDP socket
			select {
			case r.socketCh <- r.data:
				// Sent the data
			default:
				// Dropped the frame!
			}

		}
	}
}
