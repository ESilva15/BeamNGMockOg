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

func stringifyOutgaugeData(s *strings.Builder, sdk *bngsdk.BeamNGSDK) {
	// NOTE: write a string serialization function on the SDK itself
	fmt.Fprint(s, "Outgauge {\n")
	fmt.Fprintf(s, "  Time:        %d ms\n", sdk.Data.Time)
	fmt.Fprintf(s, "  Car:         %s\n", sdk.Data.Car)
	fmt.Fprintf(s, "  Flags:       %b\n", sdk.Data.Flags)
	fmt.Fprintf(s, "  Gear:        %d\n", sdk.Data.Gear)
	fmt.Fprintf(s, "  Plid:        %d\n", sdk.Data.Plid)
	fmt.Fprintf(s, "  Speed:       %f m/s\n", sdk.Data.Speed)
	fmt.Fprintf(s, "  RPM:         %f RPM\n", sdk.Data.RPM)
	fmt.Fprintf(s, "  Turbo:       %f Bar\n", sdk.Data.Turbo)
	fmt.Fprintf(s, "  EngTemp:     %f °C\n", sdk.Data.EngTemp)
	fmt.Fprintf(s, "  Fuel:        %f\n", sdk.Data.Fuel)
	fmt.Fprintf(s, "  OilPressure: %f Bar\n", sdk.Data.OilPressure)
	fmt.Fprintf(s, "  OilTemp:     %f °C\n", sdk.Data.OilTemp)
	fmt.Fprintf(s, "  DashLights:  %b\n", sdk.Data.DashLights)
	fmt.Fprintf(s, "  ShowLights:  %b\n", sdk.Data.ShowLights)
	fmt.Fprintf(s, "  Throttle:    %f\n", sdk.Data.Throttle)
	fmt.Fprintf(s, "  Brakes:      %f\n", sdk.Data.Brake)
	fmt.Fprintf(s, "  Clutch:      %f\n", sdk.Data.Clutch)
	fmt.Fprintf(s, "  Display1:    %s\n", sdk.Data.Display1)
	fmt.Fprintf(s, "  Display2:    %s\n", sdk.Data.Display2)
	fmt.Fprintf(s, "  ID:          %d\n", sdk.Data.ID)
	fmt.Fprint(s, "}\n\n")

	fmt.Fprint(s, "DashLights {\n")
	fmt.Fprintf(s, "  DL_SHIFT:      %t\n", sdk.HasShiftLight())
	fmt.Fprintf(s, "  DL_FULLBEAM:   %t\n", sdk.HasHighBeamLight())
	fmt.Fprintf(s, "  DL_HANDBRAKE:  %t\n", sdk.HasHandbrakeLight())
	fmt.Fprintf(s, "  DL_PITSPEED:   %t\n", sdk.HasPitspeed())
	fmt.Fprintf(s, "  DL_TC:         %t\n", sdk.HasTractionControlLight())
	fmt.Fprintf(s, "  DL_SIGNAL_L:   %t\n", sdk.HasLeftIndicatorLight())
	fmt.Fprintf(s, "  DL_SIGNAL_R:   %t\n", sdk.HasRightIndicatorLight())
	fmt.Fprintf(s, "  DL_SIGNAL_ANY: %t\n", sdk.HasAnyIndicatorLight())
	fmt.Fprintf(s, "  DL_OILWARN:    %t\n", sdk.HasOilLight())
	fmt.Fprintf(s, "  DL_BATTERY:    %t\n", sdk.HasBatteryLight())
	fmt.Fprintf(s, "  DL_ABS:        %t\n", sdk.HasABSLight())
	fmt.Fprintf(s, "  DL_SPARE:      %t\n", sdk.Data.DashLights&bngsdk.DL_SPARE != 0)
	fmt.Fprint(s, "}\n\n")

	fmt.Fprint(s, "ShowLights {\n") // Fixed typo "ShowLigths"
	fmt.Fprintf(s, "  DL_SHIFT:      %t\n", sdk.ShiftLight())
	fmt.Fprintf(s, "  DL_FULLBEAM:   %t\n", sdk.HighBeam())
	fmt.Fprintf(s, "  DL_HANDBRAKE:  %t\n", sdk.Handbrake())
	fmt.Fprintf(s, "  DL_PITSPEED:   %t\n", sdk.Pitspeed())
	fmt.Fprintf(s, "  DL_TC:         %t\n", sdk.TractionControl())
	fmt.Fprintf(s, "  DL_SIGNAL_L:   %t\n", sdk.LeftIndicator())
	fmt.Fprintf(s, "  DL_SIGNAL_R:   %t\n", sdk.RightIndicator())
	fmt.Fprintf(s, "  DL_SIGNAL_ANY: %t\n", sdk.AnyIndicator())
	fmt.Fprintf(s, "  DL_OILWARN:    %t\n", sdk.OilLight())
	fmt.Fprintf(s, "  DL_BATTERY:    %t\n", sdk.BatteryLight())
	fmt.Fprintf(s, "  DL_ABS:        %t\n", sdk.ABS())
	fmt.Fprintf(s, "  DL_SPARE:      %t\n", sdk.Data.ShowLights&bngsdk.DL_SPARE != 0)
	fmt.Fprint(s, "}\n\n")

	fmt.Fprint(s, "Flags {\n")
	fmt.Fprintf(s, "  OG_TURBO (Has Turbo): %t\n", sdk.HasTurbo())
	fmt.Fprintf(s, "  OG_KM (Is Metric):    %t\n", sdk.PrefersKm())
	fmt.Fprintf(s, "  OG_BAR (Pressure):    %t\n", sdk.PrefersBAR())
	fmt.Fprint(s, "}")
}
