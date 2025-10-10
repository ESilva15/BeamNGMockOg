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

// Replay replays a given file <fp> in a UDP server <addr>:<port>
func Replay(ctx context.Context, address string, port int, loop bool, fp string) error {
	fileInfo, err := os.Stat(fp)
	if err != nil {
		return fmt.Errorf("error stating file: %v", err)
	}

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

			_, err = udp.Send(data)
			if err != nil {
				return err
			}

			percent := int(float64(reader.TotalRead) / float64(fileInfo.Size()) * 100)
			fmt.Printf("\rReplayed: %d%%", percent)
		}
	}
}
