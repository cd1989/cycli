package console

import (
	"github.com/fatih/color"
	"github.com/kyokomi/emoji"
)

var (
	ErrorColor *color.Color
	InfoColor  *color.Color
)

func init() {
	ErrorColor = color.New(color.FgRed)
	InfoColor = color.New(color.FgHiCyan)
}
func Error(msg ...interface{}) {
	ErrorColor.Print(msg...)
	ErrorColor.Println(emoji.Sprint(" :confused:"))
}

func Info(msg ...interface{}) {
	InfoColor.Println(msg...)
}
