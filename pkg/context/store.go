package context

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/cd1989/cycli/pkg/console"
)

type Context struct {
	Tenant  string `yaml:"tenant"`
	Project string `yaml:"project"`
}

var Ctx *Context
var ContextFilePath string

func Init() {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		console.Error("Environment variable $HOME not set")
		os.Exit(1)
	}

	contextDir := fmt.Sprintf("%s/.cyclone", homeDir)
	if err := os.MkdirAll(contextDir, os.ModePerm); err != nil {
		console.Error("Create context dir error: ", err)
		os.Exit(1)
	}

	ContextFilePath = fmt.Sprintf("%s/context.yaml", contextDir)
	if _, err := os.Stat(ContextFilePath); err != nil {
		if os.IsNotExist(err) {
			Ctx = &Context{}
		} else {
			console.Error("Check context file error: ", err)
			os.Exit(1)
		}
	} else {
		Load()
	}
}

func Load() {
	b, err := ioutil.ReadFile(ContextFilePath)
	if err != nil {
		console.Error("Read context file error: ", err)
		os.Exit(1)
	}

	Ctx = &Context{}
	err = yaml.Unmarshal(b, Ctx)
	if err != nil {
		console.Error("Unmarshal context error: ", err)
		os.Exit(1)
	}
}

func Save() error {
	b, err := yaml.Marshal(Ctx)
	if err != nil {
		console.Error("Marshal context error: ", err)
		return err
	}

	err = ioutil.WriteFile(ContextFilePath, b, 0644)
	if err != nil {
		console.Error("Save context error: ", err)
		return err
	}

	return nil
}
