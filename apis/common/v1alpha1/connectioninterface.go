package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConnectionInterfaceSpecItemSecret struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

type ConnectionInterfaceSpecItem struct {
	Key       string `json:"key,omitempty"`
	Value     string `json:"value,omitempty"`
	Sensitive bool   `json:"sensitive,omitempty"`
}

type ConnectionInterfaceSpec struct {
	SharedInfra Ref                           `json:"sharedInfra,omitempty"`
	Outputs     []ConnectionInterfaceSpecItem `json:"outputs,omitempty"`
	Secret      Ref                           `json:"secret,omitempty"`
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
