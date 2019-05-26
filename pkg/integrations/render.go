package integrations

import (
	"os"
	"time"

	api "github.com/caicloud/cyclone/pkg/server/apis/v1alpha1"
	"github.com/fatih/color"
	"github.com/kyokomi/emoji"

	"github.com/cd1989/cycli/pkg/console/table"
)

func getURL(itg *api.Integration) string {
	switch itg.Spec.Type {
	case api.SCM:
		return itg.Spec.SCM.Server
	case api.DockerRegistry:
		return itg.Spec.DockerRegistry.Server
	case api.SonarQube:
		return itg.Spec.SonarQube.Server
	case api.Cluster:
		return itg.Spec.Cluster.Credential.Server
	}
	return ""
}

func RenderIntegrationItems(itgs []*api.Integration) {
	var rows [][]string
	for _, itg := range itgs {
		rows = append(rows, []string{
			itg.Name,
			string(itg.Spec.Type),
			itg.CreationTimestamp.Format(time.RFC3339),
			getURL(itg)})
	}

	tableWriter := &table.TableWriter{
		Writer:          os.Stdout,
		ColPadding:      4,
		Header:          []string{"NAME", "TYPE", "CREATED", "URL"},
		HeaderDecorator: color.New(color.FgBlue, color.Bold).SprintFunc(),
		Rows:            rows,
		RowEmoji:        emoji.Sprint(":beer:"),
	}

	tableWriter.Write()
}
