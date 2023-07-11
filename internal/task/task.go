package task

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const (
	TerraformTaskType = "terraform"
)

const (
	ThisInterpolationOrigin       = "this"
	TaskOutputInterpolationOrigin = "task-output"
)

type TaskInput struct {
	Label    string `json:"label"`
	Name     string `json:"name"`
	Help     string `json:"help"`
	Required bool   `json:"required"`
	Type     string `json:"type"`
	Default  string `json:"default"`
}

type TaskSpec struct {
	Inputs []TaskInput `json:"inputs"`
}

type Task struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec TaskSpec `json:"spec,omitempty"`
}
