# BeamNGMockOG
Utility to both record (with the caveat that it doesn't record yet) and replay
binary data relating to the Outgauge struct of beamng.

This will help me develop some stuff I'm doing for `ESDI` and my peripherals
without having to be playing the game while doing it (fans are noisy).


# Installation
`go install github.com/ESilva15/BeamNGMockOg@latest`


# Usage
To start replaying:
`telemetrymockserver beamng replay [--loop] -a 127.0.0.1 -p 4443 -i sunburstManual.bin`
- `loop` allows the replay functionality to keep replaying the same data file


To start recording:
`BeamNGMockOg beamng record -a 127.0.0.1 -p 4443 -o sunburstDCT.bin`

## Development
Use `tcpdump` to listen to the socket and check if data is coming through:
`tcpdump -i any udp port <port> -X`

### TODO:
- [ ] Create the auto completion file and add it to `$FPATH`
- [ ] Create a shortcut so we don't have to type out `telemetrymockserver` everytime


### Benchmark
Run with:
`go test ./mockserver -bench=BenchmarkReplayAsync -benchmem -benchtime=1x -memprofile=mem.pprof`

Analyze the output with:
`go tool pprof -sample_index=alloc_objects mem.pprof` -> `top`
- `ignore=net` will ignore the net package for example
- `focus=mockserver` will show only the results from this code we are testing
or
`go tool pprof -http=:8080 mem.pprof`


#### Previous results:
Note: this were made before introducing the real time visualizer
```
goos: linux
goarch: amd64
pkg: github.com/ESilva15/BeamNGMockOg/mockserver
cpu: AMD Ryzen 5 5600G with Radeon Graphics
BenchmarkReplayAsync-12    	      1	109433931841 ns/op	 356088 B/op	  13272 allocs/op
PASS
ok  	github.com/ESilva15/BeamNGMockOg/mockserver	109.440s
```


#### Current results
```
goos: linux
goarch: amd64
pkg: github.com/ESilva15/BeamNGMockOg/mockserver
cpu: AMD Ryzen 7 5800X3D 8-Core Processor
BenchmarkReplayAsync-16    	      1	109433890158 ns/op	1673128 B/op	 112986 allocs/op
PASS
ok  	github.com/ESilva15/BeamNGMockOg/mockserver	109.439s
```
