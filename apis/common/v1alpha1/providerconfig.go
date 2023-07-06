package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AWSProviderConfig struct {
	Role   string `json:"role,omitempty"`
	Region string `json:"region"`
}

type ProviderConfigSpec struct {
	Type      string            `json:"type,omitempty"`
	Source    string            `json:"source,omitempty"`
	AWSConfig AWSProviderConfig `json:"awsConfig,omitempty"`
	SecretRef Ref               `json:"secretRef,omitempty"`
}

type ProviderConfigStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ProviderConfig is the Schema for the circles API
type ProviderConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProviderConfigSpec   `json:"spec,omitempty"`
	Status ProviderConfigStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ProviderConfigList contains a list of ProviderConfig
type ProviderConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProviderConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ProviderConfig{}, &ProviderConfigList{})
}
