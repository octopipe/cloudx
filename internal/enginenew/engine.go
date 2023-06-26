package engine

import (
	"sync"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/terraform"
	"go.uber.org/zap"
)

const (
	ExecutionSuccessStatus = "SUCCESS"
	ExecutionFailedStatus  = "FAILED"
	ExecutionErrorStatus   = "ERROR"
	ExecutionRunningStatus = "RUNNING"
	ExecutionTimeout       = "TIMEOUT"
)

type Engine interface {
}

type engine struct {
	logger            *zap.Logger
	terraformProvider terraform.TerraformProvider
	mu                sync.Mutex

	// currentSharedInfra commonv1alpha1.SharedInfra
	// dependencyGraph    map[string][]string
	// executionGraph     map[string][]string
	// executedNodes      map[string]providerIO.ProviderOutput
}

func NewEngine(logger *zap.Logger, terraformProvider terraform.TerraformProvider) Engine {
	return engine{
		logger:            logger,
		terraformProvider: terraformProvider,
	}
}

func Apply(sharedInfra commonv1alpha1.SharedInfra) {

}
