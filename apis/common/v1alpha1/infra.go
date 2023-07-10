package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type InfraTaskRef struct {
	Type string `json:"type"`
	Path string `json:"path"`
}

type InfraTaskInput struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	Sensitive bool   `json:"sensitive"`
}

type InfraTaskOutput struct {
	Key       string `json:"key"`
	Sensitive bool   `json:"sensitive"`
}

type InfraTask struct {
	Name             string            `json:"name"`
	Ref              string            `json:"ref"`
	Depends          []string          `json:"depends,omitempty"`
	TaskType         string            `json:"type"`
	TerraformVersion string            `json:"terraformVersion,omitempty"`
	Inputs           []InfraTaskInput  `json:"inputs"`
	Outputs          []InfraTaskOutput `json:"outputs,omitempty"`
}

type InfraRunnerConfig struct {
	NodeSelector   string `json:"nodeSelector,omitempty"`
	ServiceAccount string `json:"serviceAccount,omitempty"`
}

type Ref struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

type InfraSpec struct {
	Author            string            `json:"author,omitempty" default:"anonymous"`
	Description       string            `json:"description,omitempty"`
	Generation        string            `json:"generation,omitempty"`
	ProviderConfigRef Ref               `json:"providerConfigRef,omitempty"`
	RunnerConfig      InfraRunnerConfig `json:"runnerConfig,omitempty"`
	Tasks             []InfraTask       `json:"tasks"`
}

type TaskExecutionStatus struct {
	Name           string           `json:"name"`
	Ref            string           `json:"ref"`
	Depends        []string         `json:"depends,omitempty"`
	TaskType       string           `json:"type"`
	Inputs         []InfraTaskInput `json:"inputs"`
	State          string           `json:"state,omitempty"`
	DependencyLock string           `json:"dependencyLock,omitempty"`
	StartedAt      string           `json:"startedAt,omitempty"`
	FinishedAt     string           `json:"finishedAt,omitempty"`
	Status         string           `json:"status,omitempty"`
	Error          string           `json:"error,omitempty"`
}

type ExecutionStatus struct {
	Tasks      []TaskExecutionStatus `json:"tasks,omitempty"`
	StartedAt  string                `json:"startedAt,omitempty"`
	FinishedAt string                `json:"finishedAt,omitempty"`
	Status     string                `json:"status,omitempty"`
	Error      string                `json:"error,omitempty"`
}

type InfraStatus struct {
	LastExecution ExecutionStatus `json:"lastExecution,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Circle is the Schema for the circles API
type Infra struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InfraSpec   `json:"spec,omitempty"`
	Status InfraStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CircleList contains a list of Circle
type InfraList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Infra `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Infra{}, &InfraList{})
}