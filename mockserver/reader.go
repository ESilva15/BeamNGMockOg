package mockserver

import (
	"io"
	"unsafe"

	sdk "github.com/ESilva15/gobngsdk"
)

type GobReader struct {
	TotalRead int64
	File      io.ReadSeeker
	Buf       []byte
}

func NewGobReader(r io.ReadSeeker) *GobReader {
	return &GobReader{
		TotalRead: 0,
		File:      r,
		Buf:       make([]byte, unsafe.Sizeof(sdk.Outgauge{})),
	}
}

func (g *GobReader) Reset() error {
	_, err := g.File.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	g.TotalRead = 0

	return nil
}

func (g *GobReader) Next(buffer []byte) error {
	_, err := io.ReadFull(g.File, buffer)
	if err != nil {
		return err
	}

	pos, _ := g.File.Seek(0, io.SeekCurrent)
	g.TotalRead = pos

	return nil
}
