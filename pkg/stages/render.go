package stages

import (
	"os"

	"github.com/caicloud/cyclone/pkg/apis/cyclone/v1alpha1"
	"github.com/kyokomi/emoji"

	"github.com/cd1989/cycli/pkg/console"
	"github.com/cd1989/cycli/pkg/console/table"
	"github.com/fatih/color"
)

func RenderItem(stg *v1alpha1.Stage) {
	console.Info(emoji.Sprint(":dolphin:"), stg.Name, emoji.Sprint(":smiley:"))
}

func RenderTplItems(tpls []v1alpha1.Stage) {
	var rows [][]string
	for _, tpl := range tpls {
		rows = append(rows, []string{tpl.Name, getScope(&tpl), getType(&tpl), getCreated(&tpl), getDesc(&tpl)})
	}

	tableWriter := &table.TableWriter{
		Writer:          os.Stdout,
		ColPadding:      4,
		Header:          []string{"Name", "Scope", "Type", "Created", "Desc"},
		HeaderDecorator: color.New(color.BgBlue, color.FgWhite, color.Bold).SprintFunc(),
		Rows:            rows,
	}

	tableWriter.Write()
}
