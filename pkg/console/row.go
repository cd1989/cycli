package console

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

type Col struct {
	Text   string
	Render func(format string, a ...interface{}) string
}

func (c Col) String() string {
	return c.Render(c.Text)
}

func Rows(rows [][]Col) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 8, 4, 1, ' ', 0)
	defer w.Flush()

	for _, row := range rows {
		fmt.Fprintln(w, convert(row))
	}
}

func convert(cols []Col) string {
	var outputs []string
	for _, c := range cols {
		outputs = append(outputs, c.String())
	}
	return strings.Join(outputs, "\t") + "\t"
}
