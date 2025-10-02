package mockserver

import (
	"encoding/gob"
	"log"
	"os"
	"time"

	bngsdk "github.com/ESilva15/gobngsdk"
)

// Record records data from the UDP connection created by address and port
func Record(address string, port int, filePath string) error {
	ticker := time.NewTicker(time.Second / 60)
	defer ticker.Stop()

	// Create the output file
	bin, err := os.Create(filePath)
	if err != nil {
		return err
	}

	// Create the BeamNGSDK instance
	beam, err := bngsdk.Init(address, port)
	if err != nil {
		return err
	}
	defer beam.Close()

	enc := gob.NewEncoder(bin)
	for {
		err := beam.ReadData()
		if err != nil {
			return err
		}
		if err := enc.Encode(beam.Data); err != nil {
			log.Fatal(err)
		}
		<-ticker.C
	}
}
