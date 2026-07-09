package cmd

import (
	"context"
	"fmt"

	beamng "github.com/ESilva15/TelemetryMockserver/internal/mockservers/beamng"

	"github.com/spf13/cobra"
)

// beamNGCmd is the parent command for the BeamNG agent actions
var beamNGCmd = &cobra.Command{
	Use:   "beamng",
	Short: "BeamNG telemetry utilities",
	Args:  nil,
}

var beamNGRecordCMD = &cobra.Command{
	Use:   "record",
	Short: "record -o <path-to-output-file>",
	Long: `record will store the data from the given UDP server to the 
	filepath given by -o`,
	Args: nil,
	Run:  recordAction,
}

var beamNGReplayCMD = &cobra.Command{
	Use:   "replay",
	Short: "replay -i <path-to-input-file>",
	Long:  "replay will replay the data on the given filepath on a UDP server",
	Args:  nil,
	Run:   replayAction,
}

func recordAction(cmd *cobra.Command, args []string) {
	outputFile, _ := cmd.Flags().GetString("output")
	address, _ := cmd.Flags().GetString("address")
	port, _ := cmd.Flags().GetInt("port")

	recorder, err := beamng.NewRecorder(outputFile, address, port)
	if err != nil {
		panic(fmt.Sprintf("failed to create new BeamNG recorder: %+v", err))
	}
	defer recorder.Close()

	// NOTE: is this doing anything at all??
	ctx := context.Background()
	if err := recorder.Record(ctx); err != nil {
		fmt.Printf("Something went wrong while recording the file: %v", err)
	}
}

func replayAction(cmd *cobra.Command, args []string) {
	inputFile, _ := cmd.Flags().GetString("input")
	address, _ := cmd.Flags().GetString("address")
	port, _ := cmd.Flags().GetInt("port")
	loop, _ := cmd.Flags().GetBool("loop")

	replayer, err := beamng.NewReplayer(address, port, inputFile)
	if err != nil {
		fmt.Printf("Something went wrong setting up the player: %+v", err)
		return
	}

	// NOTE: is this doing anything at all??
	ctx := context.Background()
	if err := replayer.Replay(ctx, loop); err != nil {
		fmt.Printf("Something went wrong while playing the file: %v", err)
	}
}

func init() {
	// Beamng commnad flags
	beamNGCmd.PersistentFlags().StringP("address", "a", "127.0.0.1", "Address for the UDP server")
	beamNGCmd.MarkFlagRequired("address")
	beamNGCmd.PersistentFlags().IntP("port", "p", 4444, "Port for the UDP server")
	beamNGCmd.MarkFlagRequired("port")

	// Record command flags
	beamNGRecordCMD.Flags().StringP("output", "o", "output.bin", "output file for recording")
	beamNGRecordCMD.MarkFlagRequired("output")

	// Replay command flags
	beamNGReplayCMD.Flags().BoolP("loop", "l", false, "whether to loop the recording")
	beamNGReplayCMD.Flags().StringP("input", "i", "input.bin", "input file for reading")
	beamNGReplayCMD.MarkFlagRequired("input")

	rootCmd.AddCommand(beamNGCmd)
	beamNGCmd.AddCommand(beamNGRecordCMD, beamNGReplayCMD)
}
