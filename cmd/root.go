package cmd

import (
	"fmt"
	"os"

	"github.com/kyokomi/emoji"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "cycli",
	Short: "CLI to interact with Cyclone",
	Long:  "CLI to interact with Cyclone",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(emoji.Sprintf(":confused: : %v", err))
		os.Exit(1)
	}
}
