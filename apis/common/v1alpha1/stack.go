package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type StackPluginRef struct {
	Type string `json:"type"`
	Path string `json:"path"`
}

type StackPluginOutput struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type StackPluginInput struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type StackPlugin struct {
	Name    string             `json:"name"`
	Ref     StackPluginRef     `json:"ref"`
	Depends []string           `json:"depends"`
	Inputs  []StackPluginInput `json:"inputs"`
}

type CircleSpec struct {
	Author      string        `json:"author,omitempty" default:"anonymous"`
	Description string        `json:"description,omitempty"`
	Plugins     []StackPlugin `json:"plugins"`
}

type StackStatus struct {
	LastExecutionStatus string `json:"lastExecutionStatus,omitempty"`
	LastExecutionAt     string `json:"lastExecutionAt,omitempty"`
	Error               string `json:"error,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Circle is the Schema for the circles API
type Stack struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CircleSpec  `json:"spec,omitempty"`
	Status StackStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CircleList contains a list of Circle
type StackList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Stack `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Stack{}, &StackList{})
}
