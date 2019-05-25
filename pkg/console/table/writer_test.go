package table

import (
	"os"
	"testing"

	"strings"

	"github.com/fatih/color"
	"github.com/kyokomi/emoji"
)

func TestTableWriter(t *testing.T) {
	w := TableWriter{
		Writer:          os.Stdout,
		Header:          []string{color.RedString("Name"), "Value", "Author"},
		HeaderDecorator: color.New(color.FgHiGreen, color.Bold).SprintFunc(),
		ColPadding:      4,
	}

	w.AddRow([]string{"A", color.BlueString("BBBBBBB"), "CCC"})
	w.AddRow([]string{"AAAAAA", "B", "CCCCCCC"})
	w.AddRow([]string{"AAAA", strings.Repeat(emoji.Sprint(":bell:"), 3), "CCCCCCC"})
	w.AddRow([]string{"AAAA", color.RedString("BBBBBBBBBBBBBBBBBBBBBBBB"), "CCC"})
	w.Write()
}
