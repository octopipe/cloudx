package engine

import (
	"fmt"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/terraform"
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	ExecutionSuccessStatus = "SUCCESS"
	ExecutionFailedStatus  = "FAILED"
	ExecutionErrorStatus   = "ERROR"
	ExecutionRunningStatus = "RUNNING"
	ExecutionTimeout       = "TIMEOUT"
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
	k8sClient         client.Client
}

func NewEngine(logger *zap.Logger, k8sClient client.Client, terraformProvider terraform.TerraformProvider) engine {
	return engine{
		logger:            logger,
		terraformProvider: terraformProvider,
		k8sClient:         k8sClient,
	}
}

func (e engine) getPluginsForDeletion(lastExecution commonv1alpha1.Execution, sharedInfra commonv1alpha1.SharedInfra) map[string][]string {
	forDeletion := map[string][]string{}
	for _, lastPluginExecution := range lastExecution.Status.Plugins {
		foundPlugin := false
		for _, currentPlugin := range sharedInfra.Spec.Plugins {
			if lastPluginExecution.Name == currentPlugin.Name {
				foundPlugin = true
				break
			}
		}

		if !foundPlugin {
			forDeletion[lastPluginExecution.Name] = lastPluginExecution.Depends
		}
	}

	return forDeletion
}

func (e engine) Apply(lastExecution commonv1alpha1.Execution, sharedInfra commonv1alpha1.SharedInfra) commonv1alpha1.ExecutionStatus {
	pluginsForDeletion := e.getPluginsForDeletion(lastExecution, sharedInfra)

	pipelineForDeletion := NewPipeline(e.logger, e.k8sClient, e.terraformProvider)
	dependencyGraphForDeletion := DependencyGraph{}
	for plugin := range pluginsForDeletion {
		dependencyGraphForDeletion[plugin] = []string{}
	}

	e.logger.Info(fmt.Sprintf("%d plugins for deletion", len(pluginsForDeletion)))
	for plugin, deps := range pluginsForDeletion {
		for _, dep := range deps {
			dependencyGraphForDeletion[dep] = append(dependencyGraphForDeletion[dep], plugin)
		}
	}

	pipelineForDeletion.Execute(DestroyAction, dependencyGraphForDeletion, lastExecution, sharedInfra)

	dependencyGraphForApply := DependencyGraph{}
	for _, plugin := range sharedInfra.Spec.Plugins {
		dependencyGraphForApply[plugin.Name] = plugin.Depends
	}

	pipelineForApply := NewPipeline(e.logger, e.k8sClient, e.terraformProvider)
	return pipelineForApply.Execute(ApplyAction, dependencyGraphForApply, lastExecution, sharedInfra)
}

func (e engine) Destroy(lastExecution commonv1alpha1.Execution, sharedInfra commonv1alpha1.SharedInfra) commonv1alpha1.ExecutionStatus {
	pipeline := NewPipeline(e.logger, e.k8sClient, e.terraformProvider)
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
