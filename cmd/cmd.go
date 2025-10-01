// Package cmd will be the CLI of our BeamNG OG mockserver
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bngmock",
	Short: "CLI for a BeamNG OG mockserver",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("address", "a", "127.0.0.1", "Address for the UDP server")
	rootCmd.PersistentFlags().IntP("port", "p", 4444, "Port for the UDP server")
}
