package stages

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

func Get(cmd *cobra.Command, args []string, template bool) {
	if template {
		getTemplate(cmd, args)
	} else {
		getStage(cmd, args)
	}
}

func getTemplate(cmd *cobra.Command, args []string) {
	// Whether get a given stage/template, args[0] gives the name
	if len(args) > 0 {
		tpl, err := client.K8sClient.CycloneV1alpha1().Stages(common.MetaNamespace(context.GetTenant())).Get(args[0], metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				fmt.Printf("Template %s %s %s\n", args[0], color.RedString("NOT FOUND"), emoji.Sprint(":beer:"))
			} else {
				console.Error("Get template error: ", err)
			}
			return
		}

		if tpl.Labels == nil || tpl.Labels[meta.LabelStageTemplate] == "" {
			fmt.Printf("Template %s %s %s\n", args[0], color.RedString("NOT FOUND"), emoji.Sprint(":beer:"))
		}

		RenderTplItems([]v1alpha1.Stage{*tpl})
		return
	}

	labelSelector := fmt.Sprintf("%s", meta.LabelStageTemplate)
	if p := common.GetFlagValue(cmd, "project"); p != "" {
		labelSelector = fmt.Sprintf("%s,%s=%s", labelSelector, meta.LabelProjectName, p)
	}
	templates, err := client.K8sClient.CycloneV1alpha1().Stages(common.MetaNamespace(context.GetTenant())).List(metav1.ListOptions{
		LabelSelector: meta.LabelStageTemplate,
	})
	if err != nil {
		console.Error("List templates error: ", err)
		return
	}

	RenderTplItems(templates.Items)
}

func getStage(cmd *cobra.Command, args []string) {
	// Whether get a given stage/template, args[0] gives the name
	if len(args) > 0 {
		stg, err := client.K8sClient.CycloneV1alpha1().Stages(common.MetaNamespace(context.GetTenant())).Get(args[0], metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				fmt.Printf("Stage %s %s %s\n", args[0], color.RedString("NOT FOUND"), emoji.Sprint(":beer:"))
			} else {
				console.Error("Get stage error: ", err)
			}
			return
		}

		if stg.Labels != nil && stg.Labels[meta.LabelStageTemplate] != "" {
			fmt.Printf("Stage %s %s %s\n", args[0], color.RedString("NOT FOUND"), emoji.Sprint(":beer:"))
		}

		RenderStgItems([]v1alpha1.Stage{*stg})
		return
	}

	labelSelector := fmt.Sprintf("!%s", meta.LabelStageTemplate)
	if p := common.GetFlagValue(cmd, "project"); p != "" {
		labelSelector = fmt.Sprintf("%s,%s=%s", labelSelector, meta.LabelProjectName, p)
	}
	stages, err := client.K8sClient.CycloneV1alpha1().Stages(common.MetaNamespace(context.GetTenant())).List(metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		console.Error("List stages error: ", err)
		return
	}

	RenderStgItems(stages.Items)
}

func getScope(tpl *v1alpha1.Stage) string {
	if tpl.Labels != nil && tpl.Labels[meta.LabelBuiltin] == meta.LabelValueTrue {
		return "system"
	} else {
		return "tenant"
	}
}

func getType(tpl *v1alpha1.Stage) string {
	if tpl.Labels != nil && tpl.Labels[common.LabelTemplateKind] != "" {
		return tpl.Labels[common.LabelTemplateKind]
	}
	return "--"
}

func getLabel(stg *v1alpha1.Stage, label string) string {
	if stg.Labels != nil && stg.Labels[label] != "" {
		return stg.Labels[label]
	}
	return "--"
}

func getAnnotation(stg *v1alpha1.Stage, annotation string) string {
	if stg.Annotations != nil && stg.Annotations[annotation] != "" {
		return stg.Annotations[annotation]
	}
	return "--"
}
