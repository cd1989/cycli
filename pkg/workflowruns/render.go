package workflowruns

import (
	"os"
	"strconv"
	"time"

	"github.com/caicloud/cyclone/pkg/apis/cyclone/v1alpha1"
	"github.com/caicloud/cyclone/pkg/meta"
	"github.com/fatih/color"

	"github.com/cd1989/cycli/pkg/common"
	"github.com/cd1989/cycli/pkg/console/table"
	"github.com/cd1989/cycli/pkg/dag"
	"github.com/kyokomi/emoji"
)

func elapseTime(wfr *v1alpha1.WorkflowRun) string {
	t := wfr.Status.Overall.LastTransitionTime.Unix() - wfr.Status.Overall.StartTime.Unix()
	if t < 0 {
		return "--"
	}

	return common.ReadableDuration(int(t))
}

func determineNodeElement(wfr *v1alpha1.WorkflowRun, stg string) string {
	if s, ok := wfr.Status.Stages[stg]; ok {
		switch s.Status.Phase {
		case v1alpha1.StatusSucceeded:
			return emoji.Sprint(":white_check_mark:")
		case v1alpha1.StatusFailed:
			return emoji.Sprint(":x:")
		case v1alpha1.StatusRunning:
			return emoji.Sprint(":cyclone:")
		}
	}

	return ""
}

func determineEdgeDecorator(wfr *v1alpha1.WorkflowRun, from, to string) func(format string, a ...interface{}) string {
	if s, ok := wfr.Status.Stages[from]; !ok {
		return nil
	} else {
		if s.Status.Phase != v1alpha1.StatusSucceeded && s.Status.Phase != v1alpha1.StatusFailed {
			return nil
		}
	}

	if s, ok := wfr.Status.Stages[to]; !ok || s.Status.Phase == v1alpha1.StatusCancelled {
		return nil
	}

	return color.GreenString
}

func RenderWorkflowRun(wfr *v1alpha1.WorkflowRun, wf *v1alpha1.Workflow) {
	RenderWorkflowRunItems([]v1alpha1.WorkflowRun{*wfr})

	color.New(color.FgCyan, color.Bold).Println("\n[WORKFLOW DAG]")
	dagRender := dag.NewAsciDAGRender(emoji.Sprint(":white_circle:"))
	for _, stg := range wf.Spec.Stages {
		dagRender.AddNode(&dag.Node{Name: stg.Name, Element: determineNodeElement(wfr, stg.Name)})
		for _, depend := range stg.Depends {
			dagRender.AddEdge(&dag.Edge{
				From:      depend,
				To:        stg.Name,
				Decorator: determineEdgeDecorator(wfr, depend, stg.Name),
			})
		}
	}
	dagRender.Render()
}

func RenderWorkflowRunItems(wfrs []v1alpha1.WorkflowRun) {
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
