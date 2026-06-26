package mockserver

import (
	"context"
	"net"
	"os"
	"testing"
)

func BenchmarkReplayAsync(b *testing.B) {
	// Dummy socket
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		b.Fatalf("failed to resolve UDP address: %v", err)
	}

	listener, err := net.ListenUDP("udp", addr)
	if err != nil {
		b.Fatalf("failed to start background UDP listener: %v", err)
	}
	defer listener.Close()

	// Drain the socket
	go func() {
		buf := make([]byte, 65535)
		for {
			_, _, err := listener.ReadFrom(buf)
			if err != nil {
				return
			}
		}
	}()

	// Extract the random port assigned by the OS
	assignedAddr := listener.LocalAddr().(*net.UDPAddr)

	// Get the file path
	filePath := "sunburstManual.bin"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		filePath = "../" + filePath
	}

	// Reset the timer
	b.ResetTimer()

	// Benchmark loop
	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithCancel(context.Background())

		replayer, err := NewReplayer(assignedAddr.IP.String(), assignedAddr.Port, filePath)
		if err != nil {
			cancel()
			b.Fatalf("Error setting up replayer: %+v", err)
		}

		err = replayer.Replay(ctx, false)
		if err != nil {
			cancel()
			b.Fatalf("Replay failed during benchmark run: %+v", err)
		}

		cancel()
	}
}
