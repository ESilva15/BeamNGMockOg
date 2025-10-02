package cmd

import (
	"fmt"

	"github.com/ESilva15/BeamNGMockOg/mockserver"

	"github.com/spf13/cobra"
)

func replayAction(cmd *cobra.Command, args []string) {
	inputFile, _ := cmd.Flags().GetString("input")
	address, _ := cmd.Flags().GetString("address")
	port, _ := cmd.Flags().GetInt("port")

	if err := mockserver.Replay(address, port, inputFile); err != nil {
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
}
