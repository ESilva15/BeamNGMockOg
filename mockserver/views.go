package mockserver

import (
	"fmt"
	"strings"

	bngsdk "github.com/ESilva15/gobngsdk"
)

const (
	_   = iota
	KiB = 1 << (10 * iota)
	MiB
	GiB
)

func stringifyRecordingProgress(s *strings.Builder, nBytes int) {
	if nBytes < KiB {
		fmt.Fprintf(s, "%d B", nBytes)
	} else if nBytes < MiB {
		fmt.Fprintf(s, "%.2f KiB", float64(nBytes)/float64(KiB))
	} else if nBytes < GiB {
		fmt.Fprintf(s, "%.2f MiB", float64(nBytes)/float64(MiB))
	} else {
		fmt.Fprintf(s, "%.2f GiB", float64(nBytes)/float64(GiB))
	}
}

func stringifyOutgaugeData(s *strings.Builder, og *bngsdk.Outgauge) {
	// NOTE: write a string serialization function on the SDK itself
	fmt.Fprint(s, "Outgauge {\n")
	fmt.Fprintf(s, "  Time:        %d ms\n", og.Time)
	fmt.Fprintf(s, "  Car:         %s\n", og.Car)
	fmt.Fprintf(s, "  Flags:       %b\n", og.Flags)
	fmt.Fprintf(s, "  Gear:        %d\n", og.Gear)
	fmt.Fprintf(s, "  Plid:        %d\n", og.Plid)
	fmt.Fprintf(s, "  Speed:       %f m/s\n", og.Speed)
	fmt.Fprintf(s, "  RPM:         %f RPM\n", og.RPM)
	fmt.Fprintf(s, "  Turbo:       %f Bar\n", og.Turbo)
	fmt.Fprintf(s, "  EngTemp:     %f °C\n", og.EngTemp)
	fmt.Fprintf(s, "  Fuel:        %f\n", og.Fuel)
	fmt.Fprintf(s, "  OilPressure: %f Bar\n", og.OilPressure)
	fmt.Fprintf(s, "  OilTemp:     %f °C\n", og.OilTemp)
	fmt.Fprintf(s, "  DashLights:  %b\n", og.DashLights)
	fmt.Fprintf(s, "  ShowLights:  %b\n", og.ShowLights)
	fmt.Fprintf(s, "  Throttle:    %f\n", og.Throttle)
	fmt.Fprintf(s, "  Brakes:      %f\n", og.Brake)
	fmt.Fprintf(s, "  Clutch:      %f\n", og.Clutch)
	fmt.Fprintf(s, "  Display1:    %s\n", og.Display1)
	fmt.Fprintf(s, "  Display2:    %s\n", og.Display2)
	fmt.Fprintf(s, "  ID:          %d\n", og.Display2)
	fmt.Fprint(s, "}\n\n")

	fmt.Fprint(s, "DashLights {\n")
	fmt.Fprintf(s, "  DL_SHIFT:      %t\n", og.DashLights&bngsdk.DL_SHIFT != 0)
	fmt.Fprintf(s, "  DL_FULLBEAM:   %t\n", og.DashLights&bngsdk.DL_FULLBEAM != 0)
	fmt.Fprintf(s, "  DL_HANDBRAKE:  %t\n", og.DashLights&bngsdk.DL_HANDBRAKE != 0)
	fmt.Fprintf(s, "  DL_PITSPEED:   %t\n", og.DashLights&bngsdk.DL_PITSPEED != 0)
	fmt.Fprintf(s, "  DL_TC:         %t\n", og.DashLights&bngsdk.DL_TC != 0)
	fmt.Fprintf(s, "  DL_SIGNAL_L:   %t\n", og.DashLights&bngsdk.DL_SIGNAL_L != 0)
	fmt.Fprintf(s, "  DL_SIGNAL_R:   %t\n", og.DashLights&bngsdk.DL_SIGNAL_R != 0)
	fmt.Fprintf(s, "  DL_SIGNAL_ANY: %t\n", og.DashLights&bngsdk.DL_SIGNAL_ANY != 0)
	fmt.Fprintf(s, "  DL_OILWARN:    %t\n", og.DashLights&bngsdk.DL_OILWARN != 0)
	fmt.Fprintf(s, "  DL_BATTERY:    %t\n", og.DashLights&bngsdk.DL_BATTERY != 0)
	fmt.Fprintf(s, "  DL_ABS:        %t\n", og.DashLights&bngsdk.DL_ABS != 0)
	fmt.Fprintf(s, "  DL_SPARE:      %t\n", og.DashLights&bngsdk.DL_SPARE != 0)
	fmt.Fprint(s, "}\n\n")

	fmt.Fprint(s, "ShowLights {\n") // Fixed typo "ShowLigths"
	fmt.Fprintf(s, "  DL_SHIFT:      %t\n", og.ShowLights&bngsdk.DL_SHIFT != 0)
	fmt.Fprintf(s, "  DL_FULLBEAM:   %t\n", og.ShowLights&bngsdk.DL_FULLBEAM != 0)
	fmt.Fprintf(s, "  DL_HANDBRAKE:  %t\n", og.ShowLights&bngsdk.DL_HANDBRAKE != 0)
	fmt.Fprintf(s, "  DL_PITSPEED:   %t\n", og.ShowLights&bngsdk.DL_PITSPEED != 0)
	fmt.Fprintf(s, "  DL_TC:         %t\n", og.ShowLights&bngsdk.DL_TC != 0)
	fmt.Fprintf(s, "  DL_SIGNAL_L:   %t\n", og.ShowLights&bngsdk.DL_SIGNAL_L != 0)
	fmt.Fprintf(s, "  DL_SIGNAL_R:   %t\n", og.ShowLights&bngsdk.DL_SIGNAL_R != 0)
	fmt.Fprintf(s, "  DL_SIGNAL_ANY: %t\n", og.ShowLights&bngsdk.DL_SIGNAL_ANY != 0)
	fmt.Fprintf(s, "  DL_OILWARN:    %t\n", og.ShowLights&bngsdk.DL_OILWARN != 0)
	fmt.Fprintf(s, "  DL_BATTERY:    %t\n", og.ShowLights&bngsdk.DL_BATTERY != 0)
	fmt.Fprintf(s, "  DL_ABS:        %t\n", og.ShowLights&bngsdk.DL_ABS != 0)
	fmt.Fprintf(s, "  DL_SPARE:      %t\n", og.ShowLights&bngsdk.DL_SPARE != 0)
	fmt.Fprint(s, "}\n\n")

	fmt.Fprint(s, "Flags {\n")
	fmt.Fprintf(s, "  OG_TURBO (Has Turbo): %t\n", og.Flags&bngsdk.OG_TURBO != 0)
	fmt.Fprintf(s, "  OG_KM (Is Metric):    %t\n", og.Flags&bngsdk.OG_KM != 0)
	fmt.Fprintf(s, "  OG_BAR (Pressure):    %t\n", og.Flags&bngsdk.OG_BAR != 0)
	fmt.Fprint(s, "}")
}
