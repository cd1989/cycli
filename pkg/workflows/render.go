package workflows

import (
	"os"
	"time"

	"github.com/caicloud/cyclone/pkg/apis/cyclone/v1alpha1"
	"github.com/caicloud/cyclone/pkg/meta"
	"github.com/fatih/color"
	"github.com/kyokomi/emoji"

	"github.com/cd1989/cycli/pkg/console/table"
)

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
