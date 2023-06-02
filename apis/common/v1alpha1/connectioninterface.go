package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConnectionInterfaceSpec struct {
	GeneratedFrom string `json:"GeneratedFrom,omitempty"`
}

type ConnectionInterfaceStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Circle is the Schema for the circles API
type ConnectionInterface struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConnectionInterfaceSpec   `json:"spec,omitempty"`
	Status ConnectionInterfaceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CircleList contains a list of Circle
type ConnectionInterfaceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ConnectionInterface `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ConnectionInterface{}, &ConnectionInterfaceList{})
}
