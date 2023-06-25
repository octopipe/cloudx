package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ExecutionSpec struct {
	Author      string `json:"author,omitempty" default:"anonymous"`
	Action      string `json:"action"`
	SharedInfra Ref    `json:"sharedInfra"`
}

type PluginExecutionStatus struct {
	Name           string                   `json:"name"`
	Ref            string                   `json:"ref"`
	Depends        []string                 `json:"depends,omitempty"`
	PluginType     string                   `json:"type"`
	Inputs         []SharedInfraPluginInput `json:"inputs"`
	State          string                   `json:"state,omitempty"`
	DependencyLock string                   `json:"dependencyLock,omitempty"`
	StartedAt      string                   `json:"startedAt,omitempty"`
	FinishedAt     string                   `json:"finishedAt,omitempty"`
	Status         string                   `json:"status,omitempty"`
	Error          string                   `json:"error,omitempty"`
}

type ExecutionStatus struct {
	Plugins    []PluginExecutionStatus `json:"plugins,omitempty"`
	StartedAt  string                  `json:"startedAt,omitempty"`
	FinishedAt string                  `json:"finishedAt,omitempty"`
	Status     string                  `json:"status,omitempty"`
	Error      string                  `json:"error,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Circle is the Schema for the circles API
type Execution struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ExecutionSpec   `json:"spec,omitempty"`
	Status ExecutionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CircleList contains a list of Circle
type ExecutionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Execution `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Execution{}, &ExecutionList{})
}
