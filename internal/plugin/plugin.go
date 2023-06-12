package plugin

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const (
	TerraformPluginType = "terraform"
)

type PluginInput struct {
	Label    string `json:"label"`
	Name     string `json:"name"`
	Help     string `json:"help"`
	Required bool   `json:"required"`
	Type     string `json:"type"`
	Default  string `json:"default"`
}

type PluginSpec struct {
	Inputs []PluginInput `json:"inputs"`
}

type Plugin struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec PluginSpec `json:"spec,omitempty"`
}
