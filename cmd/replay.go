package cmd

import (
	"context"
	"fmt"

	"github.com/ESilva15/BeamNGMockOg/mockserver"

	"github.com/spf13/cobra"
)

func replayAction(cmd *cobra.Command, args []string) {
	inputFile, _ := cmd.Flags().GetString("input")
	address, _ := cmd.Flags().GetString("address")
	port, _ := cmd.Flags().GetInt("port")
	loop, _ := cmd.Flags().GetBool("loop")

	replayer, err := mockserver.NewReplayer(address, port, inputFile)
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

// shelf
var replayCmd = &cobra.Command{
	Use:   "replay",
	Short: "replay -i <path-to-bin-file>",
	Long:  `replay will replay the data on the given filepath on a UDP server`,
	Args:  nil,
	Run:   replayAction,
}

func init() {
	rootCmd.AddCommand(replayCmd)

	replayCmd.PersistentFlags().StringP("input", "i", "nofile.bin", "input file for serving")
	replayCmd.PersistentFlags().BoolP("loop", "l", false, "loop replay after finishing")
}
