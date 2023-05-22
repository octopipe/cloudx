package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SharedInfraPluginRef struct {
	Type string `json:"type"`
	Path string `json:"path"`
}

type SharedInfraPluginOutput struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SharedInfraPluginInput struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type SharedInfraPlugin struct {
	Name       string                   `json:"name"`
	Ref        string                   `json:"ref"`
	Depends    []string                 `json:"depends"`
	PluginType string                   `json:"type"`
	Inputs     []SharedInfraPluginInput `json:"inputs"`
}

type SharedInfraSpec struct {
	Author      string              `json:"author,omitempty" default:"anonymous"`
	Description string              `json:"description,omitempty"`
	Plugins     []SharedInfraPlugin `json:"plugins"`
}

type PluginStatus struct {
	Name            string `json:"name,omitempty"`
	State           string `json:"state,omitempty"`
	ExecutionStatus string `json:"executionStatus,omitempty"`
	ExecutionAt     string `json:"executionAt,omitempty"`
	Error           string `json:"error,omitempty"`
}

type SharedInfraStatus struct {
	LastExecutionStatus string         `json:"lastExecutionStatus,omitempty"`
	LastExecutionAt     string         `json:"lastExecutionAt,omitempty"`
	Error               string         `json:"error,omitempty"`
	Plugins             []PluginStatus `json:"plugins,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Circle is the Schema for the circles API
type SharedInfra struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SharedInfraSpec   `json:"spec,omitempty"`
	Status SharedInfraStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CircleList contains a list of Circle
type SharedInfraList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SharedInfra `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SharedInfra{}, &SharedInfraList{})
}
