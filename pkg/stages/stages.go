package stages

import (
	"fmt"
	"time"

	"github.com/caicloud/cyclone/pkg/apis/cyclone/v1alpha1"
	"github.com/caicloud/cyclone/pkg/meta"
	"github.com/fatih/color"
	"github.com/kyokomi/emoji"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/cd1989/cycli/pkg/client"
	"github.com/cd1989/cycli/pkg/common"
	"github.com/cd1989/cycli/pkg/console"
)

func Get(args []string, template bool) {
	if template {
		getTemplate(args)
	} else {
		getStage(args)
	}
}

func getTemplate(args []string) {
	// Whether get a given stage/template, args[0] gives the name
	if len(args) > 0 {
		tpl, err := client.K8sClient.CycloneV1alpha1().Stages("default").Get(args[0], metav1.GetOptions{})
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

	templates, err := client.K8sClient.CycloneV1alpha1().Stages("default").List(metav1.ListOptions{
		LabelSelector: meta.LabelStageTemplate,
	})
	if err != nil {
		console.Error("List templates error: ", err)
		return
	}

	RenderTplItems(templates.Items)
}

func getStage(args []string) {
	// Whether get a given stage/template, args[0] gives the name
	if len(args) > 0 {
		stg, err := client.K8sClient.CycloneV1alpha1().Stages("default").Get(args[0], metav1.GetOptions{})
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

		RenderItem(stg)
		return
	}

	stages, err := client.K8sClient.CycloneV1alpha1().Stages("default").List(metav1.ListOptions{
		LabelSelector: fmt.Sprintf("!%s", meta.LabelStageTemplate),
	})
	if err != nil {
		console.Error("List stages error: ", err)
		return
	}

	for _, stg := range stages.Items {
		RenderItem(&stg)
	}
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

func getCreated(stg *v1alpha1.Stage) string {
	return stg.CreationTimestamp.Format(time.RFC3339)
}

func getDesc(stg *v1alpha1.Stage) string {
	if stg.Annotations != nil && stg.Annotations[meta.AnnotationDescription] != "" {
		return stg.Annotations[meta.AnnotationDescription]
	}

	return "--"
}
