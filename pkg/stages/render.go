package stages

import (
	"os"
	"time"

	"github.com/caicloud/cyclone/pkg/apis/cyclone/v1alpha1"
	"github.com/caicloud/cyclone/pkg/meta"
	"github.com/fatih/color"

	"github.com/cd1989/cycli/pkg/console/table"
	"github.com/kyokomi/emoji"
)

func RenderTplItems(tpls []v1alpha1.Stage) {
	var rows [][]string
	for _, tpl := range tpls {
		rows = append(rows, []string{
			tpl.Name,
			getScope(&tpl),
			getType(&tpl),
			tpl.CreationTimestamp.Format(time.RFC3339),
			getAnnotation(&tpl, meta.AnnotationDescription)})
	}

	tableWriter := &table.TableWriter{
		Writer:          os.Stdout,
		ColPadding:      4,
		Header:          []string{"NAME", "SCOPE", "TYPE", "CREATED", "DESC"},
		HeaderDecorator: color.New(color.FgBlue, color.Bold).SprintFunc(),
		Rows:            rows,
	}

	tableWriter.Write()
}

func RenderStgItems(stgs []v1alpha1.Stage) {
	var rows [][]string
	for _, stg := range stgs {
		rows = append(rows, []string{
			stg.Name,
			stg.Namespace,
			getLabel(&stg, meta.LabelProjectName),
			stg.CreationTimestamp.Format(time.RFC3339),
			getAnnotation(&stg, meta.AnnotationDescription)})
	}

	tableWriter := &table.TableWriter{
		Writer:          os.Stdout,
		ColPadding:      4,
		Header:          []string{"NAME", "NAMESPACE", "PROJECT", "CREATED", "DESC"},
		HeaderDecorator: color.New(color.FgBlue, color.Bold).SprintFunc(),
		Rows:            rows,
		RowEmoji:        emoji.Sprint(":cake:"),
	}

	tableWriter.Write()
}
