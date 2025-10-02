package cmd

import (
	"fmt"

	"github.com/ESilva15/BeamNGMockOg/mockserver"

	"github.com/spf13/cobra"
)

func recordAction(cmd *cobra.Command, args []string) {
	outputFile, _ := cmd.Flags().GetString("output")
	address, _ := cmd.Flags().GetString("address")
	port, _ := cmd.Flags().GetInt("port")

	if err := mockserver.Record(address, port, outputFile); err != nil {
		fmt.Printf("Something went wrong while recording the file: %v", err)
	}
}

// shelf
var recordCmd = &cobra.Command{
	Use:   "record",
	Short: "record -o <path-to-bin-file>",
	Long: `record will store the data from the given UDP server to the 
	filepath given by -i`,
	Args: nil,
	Run:  recordAction,
}

func init() {
	rootCmd.AddCommand(recordCmd)

	recordCmd.PersistentFlags().StringP("output", "o", "output.bin", "output file for recording")
}
