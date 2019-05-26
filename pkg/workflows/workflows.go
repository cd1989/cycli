package workflows

import (
	"fmt"

	"github.com/caicloud/cyclone/pkg/apis/cyclone/v1alpha1"
	"github.com/caicloud/cyclone/pkg/meta"
	"github.com/fatih/color"
	"github.com/kyokomi/emoji"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"strings"

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
		wf, err := client.K8sClient.CycloneV1alpha1().Workflows(common.MetaNamespace(context.GetTenant())).Get(args[0], metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				fmt.Printf("Workflow %s %s %s\n", args[0], color.RedString("NOT FOUND"), emoji.Sprint(":beer:"))
			} else {
				console.Error("Get Workflow error: ", err)
			}
			return
		}

		wfrs, stats := getWfStats(cmd, wf.Name)
		RenderWorkflow(wf, wfrs, stats)
		return
	}

	wfs, err := client.K8sClient.CycloneV1alpha1().Workflows(common.MetaNamespace(context.GetTenant())).List(*listOptions(cmd))
	if err != nil {
		console.Error("List Workflow error: ", err)
		return
	}

	statsMap := getWfsStats(cmd)
	var items []v1alpha1.Workflow
	for _, item := range wfs.Items {
		active, err := cmd.Flags().GetBool("active")
		if err == nil {
			if f := cmd.Flag("active"); f != nil && f.Changed {
				if _, ok := statsMap[item.Name]; ok != active {
					continue
				}
			}
		}

		items = append(items, item)
	}

	RenderWorkflowItems(items, statsMap)
}

func getLabel(wf *v1alpha1.Workflow, label string) string {
	if wf.Labels != nil && wf.Labels[label] != "" {
		return wf.Labels[label]
	}
	return "--"
}

type wfStats struct {
	succeed int
	failed  int
	running int
	pending int
}

func (s *wfStats) String() string {
	percent := "--%"
	if s.succeed+s.failed > 0 {
		p := s.succeed * 100 / (s.succeed + s.failed)
		if p >= 70 {
			percent = color.GreenString("%03d%%", p)
		} else if p >= 50 {
			percent = color.CyanString("%03d%%", p)
		} else {
			percent = color.RedString("%03d%%", p)
		}
	}

	return fmt.Sprintf("%s [%s-%s-%s-%s]",
		percent,
		fmt.Sprintf("%02d", s.succeed),
		fmt.Sprintf("%02d", s.failed),
		fmt.Sprintf("%02d", s.running),
		fmt.Sprintf("%02d", s.pending))
}

func (s *wfStats) Percent() string {
	percent := ""
	if s.succeed+s.failed > 0 {
		p := s.succeed * 100 / (s.succeed + s.failed)
		if p >= 70 {
			percent = color.New(color.FgGreen, color.Bold).Sprint(strings.Repeat("|", p/5+1))
		} else if p >= 50 {
			percent = color.New(color.FgCyan, color.Bold).Sprint(strings.Repeat("|", p/5+1))
		} else {
			percent = color.New(color.FgRed, color.Bold).Sprint(strings.Repeat("|", p/5+1))
		}
	}

	return percent
}

func getWfStats(cmd *cobra.Command, wf string) ([]v1alpha1.WorkflowRun, *wfStats) {
	wfrs, err := client.K8sClient.CycloneV1alpha1().WorkflowRuns(common.MetaNamespace(context.GetTenant())).List(metav1.ListOptions{})
	if err != nil {
		console.Error("List WorkflowRuns for stats error: ", err)
		return nil, &wfStats{}
	}

	var items []v1alpha1.WorkflowRun
	for _, item := range wfrs.Items {
		if item.Spec.WorkflowRef.Name == wf {
			items = append(items, item)
		}
	}

	stats := &wfStats{}
	for _, wfr := range items {
		switch wfr.Status.Overall.Phase {
		case v1alpha1.StatusSucceeded:
			stats.succeed++
		case v1alpha1.StatusFailed, v1alpha1.StatusCancelled:
			stats.failed++
		case v1alpha1.StatusRunning, v1alpha1.StatusWaiting:
			stats.running++
		case v1alpha1.StatusPending:
			stats.pending++
		default:
		}
	}

	return items, stats
}

func getWfsStats(cmd *cobra.Command) map[string]*wfStats {
	wfrs, err := client.K8sClient.CycloneV1alpha1().WorkflowRuns(common.MetaNamespace(context.GetTenant())).List(*listOptions(cmd))
	if err != nil {
		console.Error("List WorkflowRuns for stats error: ", err)
		return nil
	}

	statsMap := make(map[string]*wfStats)
	for _, wfr := range wfrs.Items {
		if wfr.Spec.WorkflowRef.Name == "" {
			continue
		}

		if _, ok := statsMap[wfr.Spec.WorkflowRef.Name]; !ok {
			statsMap[wfr.Spec.WorkflowRef.Name] = &wfStats{}
		}

		switch wfr.Status.Overall.Phase {
		case v1alpha1.StatusSucceeded:
			statsMap[wfr.Spec.WorkflowRef.Name].succeed++
		case v1alpha1.StatusFailed, v1alpha1.StatusCancelled:
			statsMap[wfr.Spec.WorkflowRef.Name].failed++
		case v1alpha1.StatusRunning, v1alpha1.StatusWaiting:
			statsMap[wfr.Spec.WorkflowRef.Name].running++
		case v1alpha1.StatusPending:
			statsMap[wfr.Spec.WorkflowRef.Name].pending++
		default:
		}
	}

	return statsMap
}
