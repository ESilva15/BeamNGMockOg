package mockserver

import (
	"encoding/binary"
	"log"
	"os"
	"time"

	bngsdk "github.com/ESilva15/gobngsdk"
)

// NOTE: add some visual feedback of whats happening.
// Maybe reuse the replay view function

// Record records data from the UDP connection created by address and port
func Record(address string, port int, filePath string) error {
	ticker := time.NewTicker(time.Second / 60)
	defer ticker.Stop()

	// Create the output file
	bin, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer bin.Close()

	// Create the BeamNGSDK instance
	beam, err := bngsdk.Init(address, port)
	if err != nil {
		return err
	}
	defer beam.Close()

	for {
		err := beam.ReadData()
		if err != nil {
			return err
		}

		err = binary.Write(bin, binary.LittleEndian, beam.Data)
		if err != nil {
			log.Fatal(err)
		}

		<-ticker.C
	}
}
