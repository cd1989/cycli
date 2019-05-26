package workflows

import (
	"os"
	"time"

	"github.com/caicloud/cyclone/pkg/apis/cyclone/v1alpha1"
	"github.com/caicloud/cyclone/pkg/meta"
	"github.com/fatih/color"
	"github.com/kyokomi/emoji"

	"github.com/cd1989/cycli/pkg/console/table"
	"github.com/cd1989/cycli/pkg/dag"
	"github.com/cd1989/cycli/pkg/workflowruns"
)

func RenderWorkflow(wf *v1alpha1.Workflow, wfrs []v1alpha1.WorkflowRun, stats *wfStats) {
	RenderWorkflowItems([]v1alpha1.Workflow{*wf}, map[string]*wfStats{wf.Name: stats})

	color.New(color.FgCyan, color.Bold).Println("\n[WORKFLOW DAG]")
	dagRender := dag.NewAsciDAGRender(emoji.Sprint(":large_blue_circle:"))
	for _, stg := range wf.Spec.Stages {
		dagRender.AddNode(&dag.Node{Name: stg.Name})
		for _, depend := range stg.Depends {
			dagRender.AddEdge(&dag.Edge{
				From: depend,
				To:   stg.Name,
			})
		}
	}
	dagRender.Render()

	color.New(color.FgCyan, color.Bold).Println("\n[EXECUTION RECORDS]")
	workflowruns.RenderWorkflowRunItems(wfrs)
}

func RenderWorkflowItems(wfs []v1alpha1.Workflow, statsMap map[string]*wfStats) {
	var rows [][]string
	for _, wf := range wfs {
		stats := "--"
		percent := ""
		if s, ok := statsMap[wf.Name]; ok {
			stats = s.String()
			percent = s.Percent()
		}

		rows = append(rows, []string{
			wf.Name,
			getLabel(&wf, meta.LabelProjectName),
			wf.CreationTimestamp.Format(time.RFC3339),
			stats,
			percent})
	}

	tableWriter := &table.TableWriter{
		Writer:          os.Stdout,
		ColPadding:      4,
		Header:          []string{"NAME", "PROJECT", "CREATED", "---% [S--F--R--P]", "PERCENT"},
		HeaderDecorator: color.New(color.FgBlue, color.Bold).SprintFunc(),
		Rows:            rows,
		RowEmoji:        emoji.Sprint(":cyclone:"),
	}

	tableWriter.Write()
}
