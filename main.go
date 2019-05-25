package main

import (
	"github.com/cd1989/cycli/cmd"
	"github.com/cd1989/cycli/pkg/context"
)

func main() {
	context.Init()
	cmd.Execute()
}
