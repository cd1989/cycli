package resources

import (
	"os"
	"time"

	"github.com/caicloud/cyclone/pkg/apis/cyclone/v1alpha1"
	"github.com/caicloud/cyclone/pkg/meta"
	"github.com/fatih/color"
	"github.com/kyokomi/emoji"

	"github.com/cd1989/cycli/pkg/console/table"
)

func RenderResourceItems(rscs []v1alpha1.Resource) {
	var rows [][]string
	for _, rsc := range rscs {
		rows = append(rows, []string{
			rsc.Name,
			getLabel(&rsc, meta.LabelProjectName),
			string(rsc.Spec.Type),
			rsc.CreationTimestamp.Format(time.RFC3339)})
	}

	tableWriter := &table.TableWriter{
		Writer:          os.Stdout,
		ColPadding:      4,
		Header:          []string{"NAME", "PROJECT", "TYPE", "CREATED"},
		HeaderDecorator: color.New(color.FgBlue, color.Bold).SprintFunc(),
		Rows:            rows,
		RowEmoji:        emoji.Sprint(":beer:"),
	}

	tableWriter.Write()
}
