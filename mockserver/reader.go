package mockserver

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"io"

	sdk "github.com/ESilva15/gobngsdk"
)

type GobReader struct {
	TotalRead int
	File      io.ReadSeeker
	Dec       *gob.Decoder
}

func NewGobReader(r io.ReadSeeker) *GobReader {
	return &GobReader{
		TotalRead: 0,
		File:      r,
		Dec:       gob.NewDecoder(r),
	}
}

func (g *GobReader) Reset() error {
	_, err := g.File.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	g.Dec = gob.NewDecoder(g.File)
	g.TotalRead = 0

	return nil
}

func (g *GobReader) Next() ([]byte, error) {
	var og sdk.Outgauge

	err := g.Dec.Decode(&og)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.LittleEndian, &og)
	if err != nil {
		return nil, err
	}

	g.TotalRead += len(buf.Bytes())

	return buf.Bytes(), nil
}
