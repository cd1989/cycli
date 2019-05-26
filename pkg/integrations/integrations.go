package integrations

import (
	"fmt"

	"github.com/caicloud/cyclone/pkg/apis/cyclone/v1alpha1"
	"github.com/caicloud/cyclone/pkg/meta"
	api "github.com/caicloud/cyclone/pkg/server/apis/v1alpha1"
	"github.com/fatih/color"
	"github.com/kyokomi/emoji"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"strings"

	"github.com/caicloud/cyclone/pkg/server/biz/integration"
	"github.com/cd1989/cycli/pkg/client"
	"github.com/cd1989/cycli/pkg/common"
	"github.com/cd1989/cycli/pkg/console"
	"github.com/cd1989/cycli/pkg/context"
)

var itgTypeMap = map[string]string{
	"scm":            "SCM",
	"sonarqube":      "SonarQube",
	"dockerregistry": "DockerRegistry",
	"cluster":        "Cluster",
	"general":        "General",
}

func determineType(itgType string) string {
	if itgType == "" {
		return ""
	}

	var matchedCnt int
	var matched string
	for k, v := range itgTypeMap {
		if strings.HasPrefix(k, strings.ToLower(itgType)) {
			matched = v
			matchedCnt++
		}
	}

	if matchedCnt == 1 {
		return matched
	}

	console.Error("Unknown integration type: ", itgType)
	return ""
}

func listOptions(cmd *cobra.Command) *metav1.ListOptions {
	labelSelector := ""
	itgType := determineType(common.GetFlagValue(cmd, "type"))
	if itgType != "" {
		labelSelector = fmt.Sprintf("%s=%s", meta.LabelIntegrationType, itgType)
	} else {
		labelSelector = meta.LabelIntegrationType
	}

	return &metav1.ListOptions{
		LabelSelector: labelSelector,
	}
}

func Get(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		secret, err := client.K8sClient.CoreV1().Secrets(common.MetaNamespace(context.GetTenant())).Get(args[0], metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				fmt.Printf("Integration %s %s %s\n", args[0], color.RedString("NOT FOUND"), emoji.Sprint(":beer:"))
			} else {
				console.Error("Get integration error: ", err)
			}
			return
		}

		itg, err := integration.FromSecret(secret)
		if err != nil {
			console.Error(fmt.Sprintf("Convert secret %s to integration error: %v", secret.Name, err))
			return
		}

		RenderIntegrationItems([]*api.Integration{itg})
		return
	}

	secrets, err := client.K8sClient.CoreV1().Secrets(common.MetaNamespace(context.GetTenant())).List(*listOptions(cmd))
	if err != nil {
		console.Error("Get integration secrets error: ", err)
		return
	}

	var itgs []*api.Integration
	for _, secret := range secrets.Items {
		itg, err := integration.FromSecret(&secret)
		if err != nil {
			console.Error(fmt.Sprintf("Convert secret %s to integration error: %v", secret.Name, err))
			return
		}

		itgs = append(itgs, itg)
	}

	RenderIntegrationItems(itgs)
}

func getLabel(stg *v1alpha1.Resource, label string) string {
	if stg.Labels != nil && stg.Labels[label] != "" {
		return stg.Labels[label]
	}
	return "--"
}
