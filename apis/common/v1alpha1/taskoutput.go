package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TaskOutputSpecItemSecret struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

type TaskOutputSpecItem struct {
	Key       string `json:"key,omitempty"`
	Value     string `json:"value,omitempty"`
	Sensitive bool   `json:"sensitive,omitempty"`
}

type TaskOutputSpec struct {
	Infra    Ref                  `json:"infra,omitempty"`
	TaskName string               `json:"taskName,omitempty"`
	Outputs  []TaskOutputSpecItem `json:"outputs,omitempty"`
	Secret   Ref                  `json:"secret,omitempty"`
}

type TaskOutputStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Circle is the Schema for the circles API
type TaskOutput struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TaskOutputSpec   `json:"spec,omitempty"`
	Status TaskOutputStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CircleList contains a list of Circle
type TaskOutputList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TaskOutput `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TaskOutput{}, &TaskOutputList{})
}
