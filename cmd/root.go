package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "cycli",
	Version: "v0.0.1",
	Short:   "CLI to interact with Cyclone",
	Long:    "CLI to interact with Cyclone",
}

func init() {
	rootCmd.SetVersionTemplate(`Cyclone CLI interface, {{printf "Version %s" .Version}}, Author: De Chen`)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
