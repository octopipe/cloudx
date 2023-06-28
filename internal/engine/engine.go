package engine

import (
	"fmt"
	"strings"

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

func (p *engine) validateDependencies(sharedInfra commonv1alpha1.SharedInfra) error {
	graph := map[string][]string{}
	for _, plugin := range sharedInfra.Spec.Plugins {
		graph[plugin.Name] = plugin.Depends
	}

	for _, p := range sharedInfra.Spec.Plugins {
		for _, dep := range p.Depends {
			if _, ok := graph[dep]; !ok {
				return fmt.Errorf("not found the dependency %s specified in plugin %s", dep, p.Name)
			}
		}
	}

	return nil
}

func (p *engine) validateInputInterpolations(sharedInfra commonv1alpha1.SharedInfra) error {
	graph := map[string][]string{}
	for _, plugin := range sharedInfra.Spec.Plugins {
		graph[plugin.Name] = plugin.Depends
	}

	for _, p := range sharedInfra.Spec.Plugins {
		for _, i := range p.Inputs {
			tokens := Lex(i.Value)

			for _, t := range tokens {
				if t.Type == TokenVariable {
					s := strings.Split(strings.Trim(t.Value, " "), ".")
					if len(s) != 3 {
						return fmt.Errorf("malformed input variable %s with value %s", i.Key, i.Value)
					}

					origin, name := s[0], s[1]
					if origin != ThisInterpolationOrigin && origin != ConnectionInterfaceInterpolationOrigin {
						return fmt.Errorf("invalid origin: %s for input %s interpolation with value %s", origin, i.Key, i.Value)
					}

					if origin == ThisInterpolationOrigin {
						if _, ok := graph[name]; !ok {
							return fmt.Errorf("invalid name: %s in origin this for input %s interpolation with value %s", name, i.Key, i.Value)
						}
					}
				}
			}
		}
	}

	return nil
}

func (e engine) Apply(lastExecution commonv1alpha1.Execution, sharedInfra commonv1alpha1.SharedInfra, currentExecutionStatusChann chan<- commonv1alpha1.ExecutionStatus) commonv1alpha1.ExecutionStatus {
	status := commonv1alpha1.ExecutionStatus{
		Status:  ExecutionSuccessStatus,
		Plugins: lastExecution.Status.Plugins,
	}

	if err := e.validateDependencies(sharedInfra); err != nil {
		status.Status = ExecutionErrorStatus
		status.Error = err.Error()
		return status
	}

	if err := e.validateInputInterpolations(sharedInfra); err != nil {
		status.Status = ExecutionErrorStatus
		status.Error = err.Error()
		return status
	}

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
				if _, ok := pluginsForDeletion[dep]; ok {
					dependencyGraphForDeletion[dep] = append(dependencyGraphForDeletion[dep], plugin)
				}
			}
		}

		deletionStatus := pipelineForDeletion.Execute(DestroyAction, dependencyGraphForDeletion, lastExecution, sharedInfra, currentExecutionStatusChann)
		if deletionStatus.Status != ExecutionSuccessStatus {
			status.Status = ExecutionErrorStatus
			status.Error = "Failed to delete diff plugins"
			return status
		}

		// Updating lastexecution after success on destroying plugins
		updatedLastExecutionPlugins := []commonv1alpha1.PluginExecutionStatus{}
		for _, p := range lastExecution.Status.Plugins {
			if _, ok := pluginsForDeletion[p.Name]; !ok {
				updatedLastExecutionPlugins = append(updatedLastExecutionPlugins, p)
			}

		}

		lastExecution.Status.Plugins = updatedLastExecutionPlugins
	}

	dependencyGraphForApply := DependencyGraph{}
	for _, plugin := range sharedInfra.Spec.Plugins {
		dependencyGraphForApply[plugin.Name] = plugin.Depends
	}

	pipelineForApply := NewPipeline(e.logger, e.rpcClient, e.terraformProvider)
	return pipelineForApply.Execute(ApplyAction, dependencyGraphForApply, lastExecution, sharedInfra, currentExecutionStatusChann)
}

func (e engine) Destroy(lastExecution commonv1alpha1.Execution, sharedInfra commonv1alpha1.SharedInfra, currentExecutionStatusChann chan<- commonv1alpha1.ExecutionStatus) commonv1alpha1.ExecutionStatus {
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

	return pipeline.Execute(DestroyAction, dependencyGraphForDestroy, lastExecution, sharedInfra, currentExecutionStatusChann)
}
