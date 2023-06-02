package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SecretsManagerSpec struct {
	Author       string `json:"author,omitempty" default:"anonymous"`
	Description  string `json:"description,omitempty"`
	Name         string `json:"name,omitempty"`
	SecretString string `json:"secretString,omitempty"`
	KmsKeyId     string `json:"kmsKeyId"`
}

type SecretsManagerStatus struct {
	// Executions []SharedInfraExecutionStatus `json:"executions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Circle is the Schema for the circles API
type SecretsManager struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SecretsManagerSpec   `json:"spec,omitempty"`
	Status SecretsManagerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CircleList contains a list of Circle
type SecretsManagerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SecretsManager `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SecretsManager{}, &SecretsManagerList{})
}
