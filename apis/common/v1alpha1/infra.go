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
	Sensitive bool   `json:"sensitive,omitempty"`
}

type InfraTaskOutputItem struct {
	Key       string `json:"key"`
	Sensitive bool   `json:"sensitive,omitempty"`
}

type InfraTaskOutput struct {
	Name  string                `json:"name"`
	Items []InfraTaskOutputItem `json:"items"`
}

type TerraformTask struct {
	Source         string `json:"source"`
	Version        string `json:"version,omitempty"`
	CredentialsRef Ref    `json:"credentialsRef,omitempty"`
}

type TemplateTask struct {
	Resource string `json:"resource"`
}

type InfraTask struct {
	Name        string                `json:"name"`
	Depends     []string              `json:"depends,omitempty"`
	Backend     string                `json:"backend"`
	Terraform   TerraformTask         `json:"terraform,omitempty"`
	Template    TemplateTask          `json:"template,omitempty"`
	Inputs      []InfraTaskInput      `json:"inputs"`
	TaskOutputs []InfraTaskOutput     `json:"taskOutputs,omitempty"`
	Outputs     []InfraTaskOutputItem `json:"outputs,omitempty"`
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

type TerraformTaskStatus struct {
	TerraformTask  `json:"task"`
	DependencyLock string `json:"dependencyLock,omitempty"`
	State          string `json:"state,omitempty"`
}

type TemplateTaskStatus struct {
	TerraformTask  `json:"task"`
	DependencyLock string `json:"dependencyLock,omitempty"`
	State          string `json:"state,omitempty"`
}

type TaskExecutionStatus struct {
	Name        string              `json:"name"`
	Depends     []string            `json:"depends,omitempty"`
	Backend     string              `json:"backend"`
	Inputs      []InfraTaskInput    `json:"inputs"`
	Terraform   TerraformTaskStatus `json:"terraform"`
	TaskOutputs []InfraTaskOutput   `json:"taskOutputs,omitempty"`
	Template    TemplateTaskStatus  `json:"template"`
	StartedAt   string              `json:"startedAt,omitempty"`
	FinishedAt  string              `json:"finishedAt,omitempty"`
	Status      string              `json:"status,omitempty"`
	Error       string              `json:"error,omitempty"`
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
