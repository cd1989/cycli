package workflowruns

import (
	"fmt"
	"sort"

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
		wfr, err := client.K8sClient.CycloneV1alpha1().WorkflowRuns(common.MetaNamespace(context.GetTenant())).Get(args[0], metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				fmt.Printf("WorkflowRun %s %s %s\n", args[0], color.RedString("NOT FOUND"), emoji.Sprint(":beer:"))
			} else {
				console.Error("Get WorkflowRun error: ", err)
			}
			return
		}

		wf, err := client.K8sClient.CycloneV1alpha1().Workflows(common.MetaNamespace(context.GetTenant())).Get(wfr.Spec.WorkflowRef.Name, metav1.GetOptions{})
		if err != nil {
			console.Error(fmt.Sprintf("Get Workflow %s error: %v", wfr.Spec.WorkflowRef.Name, err))
			return
		}

		RenderWorkflowRun(wfr, wf)
		return
	}

	wfrs, err := client.K8sClient.CycloneV1alpha1().WorkflowRuns(common.MetaNamespace(context.GetTenant())).List(*listOptions(cmd))
	if err != nil {
		console.Error("List WorkflowRun error: ", err)
		return
	}

	var items []v1alpha1.WorkflowRun
	for _, item := range wfrs.Items {
		if s := common.GetFlagValue(cmd, "status"); s != "" {
			if !common.PrefixMatch(string(item.Status.Overall.Phase), s) {
				continue
			}
		}

		c, err := cmd.Flags().GetBool("cleaned")
		if err == nil {
			if f := cmd.Flag("cleaned"); f != nil && f.Changed {
				if c != item.Status.Cleaned {
					continue
				}
			}
		}

		if w := common.GetFlagValue(cmd, "wf"); w != "" {
			if item.Spec.WorkflowRef.Name != w {
				continue
			}
		}

		items = append(items, item)
	}

	sort.Sort(SortByCreationTime(items))
	RenderWorkflowRunItems(items)
}

func getLabel(wfr *v1alpha1.WorkflowRun, label string) string {
	if wfr.Labels != nil && wfr.Labels[label] != "" {
		return wfr.Labels[label]
	}
	return "--"
}
