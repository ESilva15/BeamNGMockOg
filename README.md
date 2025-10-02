# BeamNGMockOG
Utility to both record (with the caveat that it doesn't record yet) and replay
binary data relating to the Outgauge struct of beamng.

This will help me develop some stuff I'm doing for `ESDI` and my peripherals
without having to be playing the game while doing it (fans are noisy).


# Installation
`go install github.com/ESilva15/BeamNGMockOg@latest`


# Usage
To start replaying:
`BeamNGMockOg replay -a 127.0.0.1 -p 4443 -i sunburstManual.bin`

To start recording:
`BeamNGMockOg record -a 127.0.0.1 -p 4443 -o sunburstDCT.bin`
