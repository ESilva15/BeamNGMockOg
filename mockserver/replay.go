// Package mockserver is the core of this program and it will have the API
// to record binary data from the UDP server and then be able to mock it
package mockserver

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
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
	var s strings.Builder
	var bytesReader bytes.Reader
	og := bngsdk.Outgauge{}

	for {
		select {
		case <-ctx.Done():
			return
		case data := <-r.dataViewCh:
			// Reset to the start of the terminal
			s.Reset()
			fmt.Fprintf(os.Stdout, "\x1b[2J\x1b[H")

			bytesReader.Reset(data.Data)
			err := binary.Read(&bytesReader, binary.LittleEndian, &og)
			if err != nil {
				fmt.Fprintf(&s, "FAILED TO PARSE DATA\nError: %+v", err)
				continue
			}

			percent := int(float64(data.SizeRead) / float64(fileInfo.Size()) * 100)
			fmt.Fprintf(&s, "Replayed: %d%%\n", percent)
			// NOTE: write a string serialization function on the SDK itself
			fmt.Fprint(&s, "Outgauge {\n")
			fmt.Fprintf(&s, "  Time:        %d ms\n", og.Time)
			fmt.Fprintf(&s, "  Car:         %s\n", og.Car)
			fmt.Fprintf(&s, "  Flags:       %b\n", og.Flags)
			fmt.Fprintf(&s, "  Gear:        %d\n", og.Gear)
			fmt.Fprintf(&s, "  Plid:        %d\n", og.Plid)
			fmt.Fprintf(&s, "  Speed:       %f m/s\n", og.Speed)
			fmt.Fprintf(&s, "  RPM:         %f RPM\n", og.RPM)
			fmt.Fprintf(&s, "  Turbo:       %f Bar\n", og.Turbo)
			fmt.Fprintf(&s, "  EngTemp:     %f °C\n", og.EngTemp)
			fmt.Fprintf(&s, "  Fuel:        %f\n", og.Fuel)
			fmt.Fprintf(&s, "  OilPressure: %f Bar\n", og.OilPressure)
			fmt.Fprintf(&s, "  OilTemp:     %f °C\n", og.OilTemp)
			fmt.Fprintf(&s, "  DashLights:  %b\n", og.DashLights)
			fmt.Fprintf(&s, "  ShowLights:  %b\n", og.ShowLights)
			fmt.Fprintf(&s, "  Throttle:    %f\n", og.Throttle)
			fmt.Fprintf(&s, "  Brakes:      %f\n", og.Brake)
			fmt.Fprintf(&s, "  Clutch:      %f\n", og.Clutch)
			fmt.Fprintf(&s, "  Display1:    %s\n", og.Display1)
			fmt.Fprintf(&s, "  Display2:    %s\n", og.Display2)
			fmt.Fprintf(&s, "  ID:          %d\n", og.Display2)
			fmt.Fprint(&s, "}\n\n")

			fmt.Fprint(&s, "DashLights {\n")
			fmt.Fprintf(&s, "  DL_SHIFT:      %t\n", og.DashLights&bngsdk.DL_SHIFT != 0)
			fmt.Fprintf(&s, "  DL_FULLBEAM:   %t\n", og.DashLights&bngsdk.DL_FULLBEAM != 0)
			fmt.Fprintf(&s, "  DL_HANDBRAKE:  %t\n", og.DashLights&bngsdk.DL_HANDBRAKE != 0)
			fmt.Fprintf(&s, "  DL_PITSPEED:   %t\n", og.DashLights&bngsdk.DL_PITSPEED != 0)
			fmt.Fprintf(&s, "  DL_TC:         %t\n", og.DashLights&bngsdk.DL_TC != 0)
			fmt.Fprintf(&s, "  DL_SIGNAL_L:   %t\n", og.DashLights&bngsdk.DL_SIGNAL_L != 0)
			fmt.Fprintf(&s, "  DL_SIGNAL_R:   %t\n", og.DashLights&bngsdk.DL_SIGNAL_R != 0)
			fmt.Fprintf(&s, "  DL_SIGNAL_ANY: %t\n", og.DashLights&bngsdk.DL_SIGNAL_ANY != 0)
			fmt.Fprintf(&s, "  DL_OILWARN:    %t\n", og.DashLights&bngsdk.DL_OILWARN != 0)
			fmt.Fprintf(&s, "  DL_BATTERY:    %t\n", og.DashLights&bngsdk.DL_BATTERY != 0)
			fmt.Fprintf(&s, "  DL_ABS:        %t\n", og.DashLights&bngsdk.DL_ABS != 0)
			fmt.Fprintf(&s, "  DL_SPARE:      %t\n", og.DashLights&bngsdk.DL_SPARE != 0)
			fmt.Fprint(&s, "}\n\n")

			fmt.Fprint(&s, "ShowLights {\n") // Fixed typo "ShowLigths"
			fmt.Fprintf(&s, "  DL_SHIFT:      %t\n", og.ShowLights&bngsdk.DL_SHIFT != 0)
			fmt.Fprintf(&s, "  DL_FULLBEAM:   %t\n", og.ShowLights&bngsdk.DL_FULLBEAM != 0)
			fmt.Fprintf(&s, "  DL_HANDBRAKE:  %t\n", og.ShowLights&bngsdk.DL_HANDBRAKE != 0)
			fmt.Fprintf(&s, "  DL_PITSPEED:   %t\n", og.ShowLights&bngsdk.DL_PITSPEED != 0)
			fmt.Fprintf(&s, "  DL_TC:         %t\n", og.ShowLights&bngsdk.DL_TC != 0)
			fmt.Fprintf(&s, "  DL_SIGNAL_L:   %t\n", og.ShowLights&bngsdk.DL_SIGNAL_L != 0)
			fmt.Fprintf(&s, "  DL_SIGNAL_R:   %t\n", og.ShowLights&bngsdk.DL_SIGNAL_R != 0)
			fmt.Fprintf(&s, "  DL_SIGNAL_ANY: %t\n", og.ShowLights&bngsdk.DL_SIGNAL_ANY != 0)
			fmt.Fprintf(&s, "  DL_OILWARN:    %t\n", og.ShowLights&bngsdk.DL_OILWARN != 0)
			fmt.Fprintf(&s, "  DL_BATTERY:    %t\n", og.ShowLights&bngsdk.DL_BATTERY != 0)
			fmt.Fprintf(&s, "  DL_ABS:        %t\n", og.ShowLights&bngsdk.DL_ABS != 0)
			fmt.Fprintf(&s, "  DL_SPARE:      %t\n", og.ShowLights&bngsdk.DL_SPARE != 0)
			fmt.Fprint(&s, "}\n\n")

			fmt.Fprint(&s, "Flags {\n")
			fmt.Fprintf(&s, "  OG_TURBO (Has Turbo): %t\n", og.Flags&bngsdk.OG_TURBO != 0)
			fmt.Fprintf(&s, "  OG_KM (Is Metric):   %t\n", og.Flags&bngsdk.OG_KM != 0)
			fmt.Fprintf(&s, "  OG_BAR (Pressure):   %t\n", og.Flags&bngsdk.OG_BAR != 0)
			fmt.Fprint(&s, "}")

			fmt.Fprint(os.Stdout, s.String())
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
