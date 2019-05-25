package table

import (
	"io"

	"fmt"

	"strings"

	"github.com/kyokomi/emoji"
	"github.com/mattn/go-runewidth"
)

type Decorator func(a ...interface{}) string

var DefaultDecorator = func(a ...interface{}) string {
	return fmt.Sprint(a...)
}

type TableWriter struct {
	Writer          io.Writer
	Header          []string
	HeaderDecorator Decorator
	ColPadding      int
	Rows            [][]string
	RowEmoji        string
}

func (w *TableWriter) SetHeader(header []string) {
	w.Header = header
}

func (w *TableWriter) AddRow(row []string) {
	w.Rows = append(w.Rows, row)
}

func (w *TableWriter) Write() {
	colsSize := make(map[int]int)
	for i, h := range w.Header {
		colsSize[i] = runewidth.StringWidth(h)
		for _, r := range w.Rows {
			l := runewidth.StringWidth(r[i])
			if l > colsSize[i] {
				colsSize[i] = l
			}
		}
	}

	w.writeRow(w.Header, colsSize, w.HeaderDecorator, true)
	for _, r := range w.Rows {
		w.writeRow(r, colsSize, nil, false)
	}
}

func padding(s string, length int) string {
	l := runewidth.StringWidth(s)
	if l >= length {
		return s
	}

	return s + strings.Repeat(" ", length-l)
}

func (w *TableWriter) writeRow(row []string, colsSize map[int]int, decorator Decorator, header bool) {
	if decorator == nil {
		decorator = DefaultDecorator
	}

	line := "   "
	if !header {
		line = w.RowEmoji
		if w.RowEmoji == "" {
			line = emoji.Sprint(":small_blue_diamond:")
		}
		line = line + " "
	}
	for i, c := range row {
		line += padding(c, colsSize[i]+w.ColPadding)
	}
	fmt.Fprintln(w.Writer, decorator(line))
}
