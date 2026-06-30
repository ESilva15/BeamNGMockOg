package mockserver

import (
	"context"
	"encoding/binary"
	"fmt"
	"os"
	"strings"
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
			r.mut.RLock()
			err := binary.Write(r.OutputFile, binary.LittleEndian, data)
			r.TotalBytes += len(data)
			r.mut.RUnlock()

			if err != nil {
				// NOTE: find a way of logging this somehow
			}
		}
	}
}

func (r *Recorder) view(ctx context.Context) {
	var s strings.Builder
	var nBytes int

	for {
		select {
		case <-ctx.Done():
			return
		case viewData := <-r.viewCh:
			s.Reset()
			fmt.Fprintf(os.Stdout, "\x1b[2J\x1b[H")

			stringifyRecordingProgress(&s, nBytes)
			fmt.Fprintf(&s, "\n\n")

			r.viewDataMut.RLock()
			nBytes = viewData.TotalBytes
			stringifyOutgaugeData(&s, &r.SDK)
			r.viewDataMut.RUnlock()

			fmt.Fprint(os.Stdout, s.String())
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

			// Send the data to the UDP socket
			select {
			case r.recorderCh <- r.SDK.Buffer:
				// Sent the data
			default:
				// Dropped the frame!
			}
		}
	}
}
