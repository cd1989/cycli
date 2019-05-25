package cmd

import (
	"github.com/spf13/cobra"

	"github.com/cd1989/cycli/pkg/common"
	"github.com/cd1989/cycli/pkg/context"
)

func init() {
	rootCmd.AddCommand(useCmd)
}

var useCmd = &cobra.Command{
	Use:   "use [tenant|project] <value>",
	Short: "Set default tenant, project in context",
	Long:  "Set default tenant, project in context",
	ValidArgs: []string{
		common.ContextTenant,
		common.ContextProject},
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			return err
		}
		if err := cobra.OnlyValidArgs(cmd, args[:1]); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case common.ContextTenant:
			context.SetTenant(args[1:]...)
		case common.ContextProject:
			context.SetProject(args[1:]...)
		}
	},
}
