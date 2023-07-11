package terraform

import (
	"github.com/hashicorp/terraform-exec/tfexec"
	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"go.uber.org/zap"
)

type TerraformApplyInput struct {
	Source           string
	Version          string
	TaskInputs       []commonv1alpha1.InfraTaskInput
	PreviousState    string
	PreviousLockDeps string
}

type TerraformApplyResult struct {
	Outputs          map[string]tfexec.OutputMeta
	DependenciesLock string
	State            string
}

type TerraformDestroyInput struct {
	Source           string
	Version          string
	TaskInputs       []commonv1alpha1.InfraTaskInput
	PreviousState    string
	PreviousLockDeps string
}

type TerraformBackend interface {
	Apply(input TerraformApplyInput) (TerraformApplyResult, error)
	Destroy(input TerraformDestroyInput) error
}

type terraformBackend struct {
	logger *zap.Logger
}

func NewTerraformBackend(logger *zap.Logger) (terraformBackend, error) {
	return terraformBackend{
		logger: logger,
	}, nil
}
