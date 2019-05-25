package workflowruns

import (
	"os"
	"strconv"
	"time"

	"github.com/caicloud/cyclone/pkg/apis/cyclone/v1alpha1"
	"github.com/caicloud/cyclone/pkg/meta"
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/cd1989/cycli/pkg/common"
	"github.com/cd1989/cycli/pkg/console/table"
)

func elapseTime(wfr *v1alpha1.WorkflowRun) string {
	t := wfr.Status.Overall.LastTransitionTime.Unix() - wfr.Status.Overall.StartTime.Unix()
	if t < 0 {
		return "--"
	}

	return common.ReadableDuration(int(t))
}

func RenderWorkflowRunItems(cmd *cobra.Command, wfrs []v1alpha1.WorkflowRun) {
	var rows [][]string
	for _, wfr := range wfrs {
		rows = append(rows, []string{
			wfr.Name,
			getLabel(&wfr, meta.LabelProjectName),
			wfr.Spec.WorkflowRef.Name,
			wfr.CreationTimestamp.Format(time.RFC3339),
			string(wfr.Status.Overall.Phase),
			strconv.FormatBool(wfr.Status.Cleaned),
			elapseTime(&wfr)})
	}

	tableWriter := &table.TableWriter{
		Writer:          os.Stdout,
		ColPadding:      2,
		Header:          []string{"NAME", "PROJECT", "WORKFLOW", "CREATED", "STATUS", "CLEANED", "TIME"},
		HeaderDecorator: color.New(color.FgBlue, color.Bold).SprintFunc(),
		Rows:            rows,
	}

	tableWriter.Write()
}
