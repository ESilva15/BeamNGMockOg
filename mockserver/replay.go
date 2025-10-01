// Package mockserver is the core of this program and it will have the API
// to record binary data from the UDP server and then be able to mock it
package mockserver

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"log"
	"os"
	"time"

	sdk "github.com/ESilva15/gobngsdk"
)

// Replay replays a given file <fp> in a UDP server <addr>:<port>
func Replay(address string, port int, fp string) error {
	bin, err := os.Open(fp)
	if err != nil {
		log.Fatal("error opening file:", err)
	}

	addr, conn, err := openUDPServer(address, port)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(time.Second / 60)
	defer ticker.Stop()

	dec := gob.NewDecoder(bin)
	for {
		var og sdk.Outgauge
		if err := dec.Decode(&og); err != nil {
			log.Fatal("failed to read more data:", err)
			break
		}

		buf := new(bytes.Buffer)
		if err := binary.Write(buf, binary.LittleEndian, &og); err != nil {
			log.Fatal("serialization failed:", err)
			break
		}

		if _, err := conn.WriteToUDP(buf.Bytes(), addr); err != nil {
			log.Fatal("send failed:", err)
			break
		}

		<-ticker.C
	}

	return conn.Close()
}
