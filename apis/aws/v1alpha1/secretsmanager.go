package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SecretManagerSpec struct {
	Author      string `json:"author,omitempty" default:"anonymous"`
	Description string `json:"description,omitempty"`
}

type SecretManagerStatus struct {
	// Executions []SharedInfraExecutionStatus `json:"executions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Circle is the Schema for the circles API
type SecretManager struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SecretManagerSpec   `json:"spec,omitempty"`
	Status SecretManagerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CircleList contains a list of Circle
type SecretManagerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SecretManager `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SecretManager{}, &SecretManagerList{})
}
