package workflowruns

import "github.com/caicloud/cyclone/pkg/apis/cyclone/v1alpha1"

type SortByCreationTime []v1alpha1.WorkflowRun

func (s SortByCreationTime) Len() int {
	return len(s)
}

func (s SortByCreationTime) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortByCreationTime) Less(i, j int) bool {
	return s[i].CreationTimestamp.After(s[j].CreationTimestamp.Time)
}
