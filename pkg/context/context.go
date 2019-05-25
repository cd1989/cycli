package context

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/kyokomi/emoji"

	"github.com/cd1989/cycli/pkg/console"
)

func SetTenant(tenant ...string) {
	value := ""
	if len(tenant) != 0 {
		value = tenant[0]
	}

	Ctx.Tenant = value
	if err := Save(); err != nil {
		console.Error("Save context error: ", err)
		os.Exit(1)
	}

	if value == "" {
		fmt.Println("Tenant cleaned in context successfully,", emoji.Sprint(":tada:"))
	} else {
		fmt.Printf("Tenant set to %s in context successfully, %s\n", color.GreenString(value), emoji.Sprint(":tada:"))
	}
}

func SetProject(project ...string) {
	value := ""
	if len(project) != 0 {
		value = project[0]
	}

	Ctx.Project = value
	if err := Save(); err != nil {
		console.Error("Save context error: ", err)
		os.Exit(1)
	}

	if value == "" {
		fmt.Println("Project cleaned in context,", emoji.Sprint(":tada:"))
	} else {
		fmt.Printf("Project set to %s in context, %s", color.GreenString(value), emoji.Sprint(":tada:"))
	}
}

func GetTenant() string {
	return Ctx.Tenant
}

func GetProject() string {
	return Ctx.Project
}
