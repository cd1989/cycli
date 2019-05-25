package cmd

import (
	"github.com/spf13/cobra"

	"github.com/cd1989/cycli/pkg/common"
	"github.com/cd1989/cycli/pkg/resources"
	"github.com/cd1989/cycli/pkg/stages"
	"github.com/cd1989/cycli/pkg/workflowruns"
	"github.com/cd1989/cycli/pkg/workflows"
	"github.com/cd1989/cycli/pkg/workflowtriggers"
)

func init() {
	getCmd.PersistentFlags().StringP("project", "p", "", "Project of the resources")
	getCmd.PersistentFlags().StringP("type", "t", "", "Type of the resources")
	rootCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{
	Use:   "get [resource type]",
	Short: "Get resources(resource, stage, workflow, etc) in Cyclone",
	Long:  "Get resources(resource, stage, workflow, etc) in Cyclone",
	ValidArgs: []string{
		common.TypeResource,
		common.TypeResourceShort,
		common.TypeStage,
		common.TypeStageShort,
		common.TypeWorkflow,
		common.TypeWorkflowShort,
		common.TypeWorkflowRun,
		common.TypeWorkflowRunShort,
		common.TypeWorkflowTrigger,
		common.TypeWorkflowTrigerShort,
		common.TypeStgTemplate,
		common.TypeStgTemplateShort},
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
		case common.TypeResource, common.TypeResourceShort:
			resources.Get(cmd, args[1:])
		case common.TypeStage, common.TypeStageShort:
			stages.Get(args[1:], false)
		case common.TypeStgTemplate, common.TypeStgTemplateShort:
			stages.Get(args[1:], true)
		case common.TypeWorkflow, common.TypeWorkflowShort:
			workflows.Get(args[1:])
		case common.TypeWorkflowRun, common.TypeWorkflowRunShort:
			workflowruns.Get(args[1:])
		case common.TypeWorkflowTrigger, common.TypeWorkflowTrigerShort:
			workflowtriggers.Get(args[1:])
		}
	},
}
