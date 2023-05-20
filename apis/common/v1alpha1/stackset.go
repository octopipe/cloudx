package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type StackSetPluginRef struct {
	Type string `json:"type"`
	Path string `json:"path"`
}

type StackSetPluginOutput struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type StackSetPluginInput struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type StackSetPlugin struct {
	Name       string                `json:"name"`
	Ref        string                `json:"ref"`
	Depends    []string              `json:"depends"`
	PluginType string                `json:"type"`
	Inputs     []StackSetPluginInput `json:"inputs"`
}

type StackSetSpec struct {
	Author      string           `json:"author,omitempty" default:"anonymous"`
	Description string           `json:"description,omitempty"`
	Plugins     []StackSetPlugin `json:"plugins"`
}

type StackSetStatus struct {
	LastExecutionStatus string `json:"lastExecutionStatus,omitempty"`
	LastExecutionAt     string `json:"lastExecutionAt,omitempty"`
	Error               string `json:"error,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Circle is the Schema for the circles API
type StackSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StackSetSpec   `json:"spec,omitempty"`
	Status StackSetStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CircleList contains a list of Circle
type StackSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StackSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&StackSet{}, &StackSetList{})
}
