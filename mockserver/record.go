package mockserver

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"os"
	"sync"
	"time"

	bngsdk "github.com/ESilva15/gobngsdk"
)

type recorderViewData struct {
	TotalBytes int
	SDK        *bngsdk.BeamNGSDK
}

type Recorder struct {
	SDK        bngsdk.BeamNGSDK
	OutputFile *os.File
	TotalBytes int
	// Views
	mut         sync.RWMutex
	viewDataMut sync.RWMutex
	viewData    recorderViewData
	viewCh      chan *recorderViewData
	recorderCh  chan []byte
}

func NewRecorder(fp string, address string, port int) (*Recorder, error) {
	var recorder Recorder
	var err error

	recorder.SDK, err = bngsdk.Init(address, port)
	if err != nil {
		return &Recorder{}, err
	}

	recorder.OutputFile, err = os.Create(fp)
	if err != nil {
		return &Recorder{}, err
	}

	recorder.viewData = recorderViewData{}
	recorder.viewCh = make(chan *recorderViewData, 1)
	recorder.recorderCh = make(chan []byte, 1)

	return &recorder, nil
}

func (r *Recorder) Close() {
	r.SDK.Close()

	if r.OutputFile != nil {
		r.OutputFile.Close()
	}
}

func (r *Recorder) record(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case data := <-r.recorderCh:
			r.mut.Lock()
			err := binary.Write(r.OutputFile, binary.LittleEndian, r.SDK.Data)
			r.TotalBytes += len(data)
			r.mut.Unlock()

			if err != nil {
				// NOTE: find a way of logging this somehow
			}
		}
	}
}

func (r *Recorder) view(ctx context.Context) {
	var buf bytes.Buffer
	var nBytes int

	buf.Grow(2048)

	for {
		select {
		case <-ctx.Done():
			return
		case viewData := <-r.viewCh:
			buf.Reset()
			buf.WriteString("\x1b[2J\x1b[H")

			stringifyRecordingProgress(&buf, nBytes)
			fmt.Fprintf(&buf, "\n\n")

			r.viewDataMut.RLock()
			nBytes = viewData.TotalBytes
			stringifyOutgaugeData(&buf, &r.SDK)
			r.viewDataMut.RUnlock()

			_, _ = buf.WriteTo(os.Stdout)
		}
	}
}

// Record records data from the UDP connection created by address and port
func (r *Recorder) Record(ctx context.Context) error {
	ticker := time.NewTicker(time.Second / 60)
	defer ticker.Stop()

	go r.record(ctx)
	go r.view(ctx)

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			err := r.SDK.ReadData()
			if err != nil {
				return err
			}

			r.viewData.TotalBytes = r.TotalBytes

			// Send the data to the view
			select {
			case r.viewCh <- &r.viewData:
				// Sent the data
			default:
				// Dropped the frame!
			}

			// Write the data to the file
			select {
			case r.recorderCh <- r.SDK.Buffer:
				// Sent the data
			default:
				// Dropped the frame!
			}
		}
	}
}
