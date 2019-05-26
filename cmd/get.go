package cmd

import (
	"github.com/spf13/cobra"

	"github.com/cd1989/cycli/pkg/common"
	"github.com/cd1989/cycli/pkg/integrations"
	"github.com/cd1989/cycli/pkg/resources"
	"github.com/cd1989/cycli/pkg/stages"
	"github.com/cd1989/cycli/pkg/workflowruns"
	"github.com/cd1989/cycli/pkg/workflows"
	"github.com/cd1989/cycli/pkg/workflowtriggers"
)

func init() {
	getCmd.Flags().StringP("project", "p", "", "Project of the resources")
	getCmd.Flags().StringP("type", "t", "", "Type of the resources")
	getCmd.Flags().StringP("status", "s", "", "Status of WorkflowRun")
	getCmd.Flags().BoolP("cleaned", "c", false, "Whether GC is performed")
	getCmd.Flags().StringP("wf", "w", "", "Workflow")
	getCmd.Flags().BoolP("active", "a", true, "Whether Workflow is active")
	rootCmd.AddCommand(getCmd)
}

var getCmd = &cobra.Command{
	Use:   "get [rsc|stg|wf|wfr|wft|itg]",
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
		common.TypeStgTemplateShort,
		common.TypeIntegration,
		common.TypeIntegrationShort},
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
			stages.Get(cmd, args[1:], false)
		case common.TypeStgTemplate, common.TypeStgTemplateShort:
			stages.Get(cmd, args[1:], true)
		case common.TypeWorkflow, common.TypeWorkflowShort:
			workflows.Get(cmd, args[1:])
		case common.TypeWorkflowRun, common.TypeWorkflowRunShort:
			workflowruns.Get(cmd, args[1:])
		case common.TypeWorkflowTrigger, common.TypeWorkflowTrigerShort:
			workflowtriggers.Get(args[1:])
		case common.TypeIntegrationShort, common.TypeIntegration:
			integrations.Get(cmd, args[1:])
		}
	},
}
