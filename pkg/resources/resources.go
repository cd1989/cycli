package resources

import (
	"fmt"

	"github.com/caicloud/cyclone/pkg/apis/cyclone/v1alpha1"
	"github.com/caicloud/cyclone/pkg/meta"
	"github.com/fatih/color"
	"github.com/kyokomi/emoji"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cd1989/cycli/pkg/client"
	"github.com/cd1989/cycli/pkg/common"
	"github.com/cd1989/cycli/pkg/console"
	"github.com/cd1989/cycli/pkg/context"
)

func listOptions(cmd *cobra.Command) *metav1.ListOptions {
	labelSelector := ""
	project := common.GetFlagValue(cmd, "project")
	if project != "" {
		labelSelector = fmt.Sprintf("%s=%s", meta.LabelProjectName, project)
	}

	return &metav1.ListOptions{
		LabelSelector: labelSelector,
	}
}

func Get(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		rsc, err := client.K8sClient.CycloneV1alpha1().Resources(common.MetaNamespace(context.GetTenant())).Get(args[0], metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				fmt.Printf("Resource %s %s %s\n", args[0], color.RedString("NOT FOUND"), emoji.Sprint(":beer:"))
			} else {
				console.Error("Get resource error: ", err)
			}
			return
		}

		RenderResourceItems([]v1alpha1.Resource{*rsc})
		return
	}

	rscs, err := client.K8sClient.CycloneV1alpha1().Resources(common.MetaNamespace(context.GetTenant())).List(*listOptions(cmd))
	if err != nil {
		console.Error("Get resource error: ", err)
		return
	}

	if t := common.GetFlagValue(cmd, "type"); t != "" {
		var items []v1alpha1.Resource
		for _, i := range rscs.Items {
			if common.Equal(string(i.Spec.Type), t) {
				items = append(items, i)
			}
		}

		RenderResourceItems(items)
	} else {
		RenderResourceItems(rscs.Items)
	}
}

func getLabel(stg *v1alpha1.Resource, label string) string {
	if stg.Labels != nil && stg.Labels[label] != "" {
		return stg.Labels[label]
	}
	return "--"
}
