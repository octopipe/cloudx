package engine

import (
	"fmt"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/rpcclient"
	"github.com/octopipe/cloudx/internal/terraform"
	"go.uber.org/zap"
)

const (
	ExecutionSuccessStatus      = "SUCCESS"
	ExecutionErrorStatus        = "ERROR"
	ExecutionAppliedStatus      = "APPLIED"
	ExecutionFailedStatus       = "FAILED"
	ExecutionApplyErrorStatus   = "APPLY_ERROR"
	ExecutionRunningStatus      = "RUNNING"
	ExecutionTimeout            = "TIMEOUT"
	ExecutionDestroyed          = "DESTROYED"
	ExecutionDestroyErrorStatus = "DESTROY_ERROR"
)

type ExecutionActionType int

const (
	ApplyAction ExecutionActionType = iota
	DestroyAction
)

type Engine interface {
}

type engine struct {
	logger            *zap.Logger
	terraformProvider terraform.TerraformProvider
	rpcClient         rpcclient.Client
}

func NewEngine(logger *zap.Logger, rpcClient rpcclient.Client, terraformProvider terraform.TerraformProvider) engine {
	return engine{
		logger:            logger,
		terraformProvider: terraformProvider,
		rpcClient:         rpcClient,
	}
}

func (e engine) getPluginsForDeletion(lastExecution commonv1alpha1.Execution, sharedInfra commonv1alpha1.SharedInfra) map[string]commonv1alpha1.PluginExecutionStatus {
	forDeletion := map[string]commonv1alpha1.PluginExecutionStatus{}
	for _, lastPluginExecution := range lastExecution.Status.Plugins {
		foundPlugin := false
		for _, currentPlugin := range sharedInfra.Spec.Plugins {
			if lastPluginExecution.Name == currentPlugin.Name {
				foundPlugin = true
				break
			}
		}

		if !foundPlugin {
			forDeletion[lastPluginExecution.Name] = lastPluginExecution
		}
	}

	return forDeletion
}

func (e engine) Apply(lastExecution commonv1alpha1.Execution, sharedInfra commonv1alpha1.SharedInfra) commonv1alpha1.ExecutionStatus {
	pluginsForDeletion := e.getPluginsForDeletion(lastExecution, sharedInfra)

	if len(pluginsForDeletion) > 0 {

		pipelineForDeletion := NewPipeline(e.logger, e.rpcClient, e.terraformProvider)
		dependencyGraphForDeletion := DependencyGraph{}
		for plugin := range pluginsForDeletion {
			dependencyGraphForDeletion[plugin] = []string{}
		}

		e.logger.Info(fmt.Sprintf("%d plugins for deletion", len(pluginsForDeletion)))
		for plugin, execution := range pluginsForDeletion {
			for _, dep := range execution.Depends {
				dependencyGraphForDeletion[dep] = append(dependencyGraphForDeletion[dep], plugin)
			}
		}

		status := pipelineForDeletion.Execute(DestroyAction, dependencyGraphForDeletion, lastExecution, sharedInfra)
		if status.Status != ExecutionSuccessStatus {
			return status
		}
	}

	dependencyGraphForApply := DependencyGraph{}
	for _, plugin := range sharedInfra.Spec.Plugins {
		dependencyGraphForApply[plugin.Name] = plugin.Depends
	}

	pipelineForApply := NewPipeline(e.logger, e.rpcClient, e.terraformProvider)
	return pipelineForApply.Execute(ApplyAction, dependencyGraphForApply, lastExecution, sharedInfra)
}

func (e engine) Destroy(lastExecution commonv1alpha1.Execution, sharedInfra commonv1alpha1.SharedInfra) commonv1alpha1.ExecutionStatus {
	pipeline := NewPipeline(e.logger, e.rpcClient, e.terraformProvider)
	dependencyGraphForDestroy := DependencyGraph{}

	for _, plugin := range sharedInfra.Spec.Plugins {
		dependencyGraphForDestroy[plugin.Name] = []string{}
	}

	for _, plugin := range sharedInfra.Spec.Plugins {
		for _, dep := range plugin.Depends {
			dependencyGraphForDestroy[dep] = append(dependencyGraphForDestroy[dep], plugin.Name)
		}
	}

	return pipeline.Execute(DestroyAction, dependencyGraphForDestroy, lastExecution, sharedInfra)
}
