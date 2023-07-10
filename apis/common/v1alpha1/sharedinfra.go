package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SharedInfraPluginRef struct {
	Type string `json:"type"`
	Path string `json:"path"`
}

type SharedInfraPluginInput struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	Sensitive bool   `json:"sensitive"`
}

type SharedInfraPluginOutput struct {
	Key       string `json:"key"`
	Sensitive bool   `json:"sensitive"`
}

type SharedInfraPlugin struct {
	Name             string                    `json:"name"`
	Ref              string                    `json:"ref"`
	Depends          []string                  `json:"depends,omitempty"`
	PluginType       string                    `json:"type"`
	TerraformVersion string                    `json:"terraformVersion,omitempty"`
	Inputs           []SharedInfraPluginInput  `json:"inputs"`
	Outputs          []SharedInfraPluginOutput `json:"outputs,omitempty"`
}

type SharedInfraRunnerConfig struct {
	NodeSelector   string `json:"nodeSelector,omitempty"`
	ServiceAccount string `json:"serviceAccount,omitempty"`
}

type Ref struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

type SharedInfraSpec struct {
	Author            string                  `json:"author,omitempty" default:"anonymous"`
	Description       string                  `json:"description,omitempty"`
	Generation        string                  `json:"generation,omitempty"`
	ProviderConfigRef Ref                     `json:"providerConfigRef,omitempty"`
	RunnerConfig      SharedInfraRunnerConfig `json:"runnerConfig,omitempty"`
	Plugins           []SharedInfraPlugin     `json:"plugins"`
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

type SharedInfraStatus struct {
	LastExecution ExecutionStatus `json:"lastExecution,omitempty"`
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
