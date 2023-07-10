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

func (e engine) getTasksForDeletion(lastExecution commonv1alpha1.ExecutionStatus, infra commonv1alpha1.Infra) map[string]commonv1alpha1.TaskExecutionStatus {
	forDeletion := map[string]commonv1alpha1.TaskExecutionStatus{}
	for _, lastTaskExecution := range lastExecution.Tasks {
		foundTask := false
		for _, currentTask := range infra.Spec.Tasks {
			if lastTaskExecution.Name == currentTask.Name {
				foundTask = true
				break
			}
		}

		if !foundTask {
			forDeletion[lastTaskExecution.Name] = lastTaskExecution
		}
	}

	return forDeletion
}

func (p *engine) validateDependencies(infra commonv1alpha1.Infra) error {
	graph := map[string][]string{}
	for _, task := range infra.Spec.Tasks {
		graph[task.Name] = task.Depends
	}

	for _, p := range infra.Spec.Tasks {
		for _, dep := range p.Depends {
			if _, ok := graph[dep]; !ok {
				return fmt.Errorf("not found the dependency %s specified in task %s", dep, p.Name)
			}
		}
	}

	return nil
}

func (p *engine) validateInputInterpolations(infra commonv1alpha1.Infra) error {
	graph := map[string][]string{}
	for _, task := range infra.Spec.Tasks {
		graph[task.Name] = task.Depends
	}

	for _, p := range infra.Spec.Tasks {
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

func (e engine) Apply(infra commonv1alpha1.Infra, currentExecutionStatusChann chan commonv1alpha1.ExecutionStatus) commonv1alpha1.ExecutionStatus {
	lastExecution := infra.Status.LastExecution
	status := commonv1alpha1.ExecutionStatus{
		Status: ExecutionSuccessStatus,
		Tasks:  lastExecution.Tasks,
	}

	if err := e.validateDependencies(infra); err != nil {
		status.Status = ExecutionErrorStatus
		status.Error = err.Error()
		currentExecutionStatusChann <- status
		return status
	}

	if err := e.validateInputInterpolations(infra); err != nil {
		status.Status = ExecutionErrorStatus
		status.Error = err.Error()
		currentExecutionStatusChann <- status
		return status
	}

	tasksForDeletion := e.getTasksForDeletion(lastExecution, infra)
	if len(tasksForDeletion) > 0 {

		pipelineForDeletion := NewPipeline(e.logger, e.rpcClient, e.terraformProvider)
		dependencyGraphForDeletion := DependencyGraph{}
		for task := range tasksForDeletion {
			dependencyGraphForDeletion[task] = []string{}
		}

		e.logger.Info(fmt.Sprintf("%d tasks for deletion", len(tasksForDeletion)))
		for task, execution := range tasksForDeletion {
			for _, dep := range execution.Depends {
				if _, ok := tasksForDeletion[dep]; ok {
					dependencyGraphForDeletion[dep] = append(dependencyGraphForDeletion[dep], task)
				}
			}
		}

		deletionStatus := pipelineForDeletion.Execute(DestroyAction, dependencyGraphForDeletion, infra, currentExecutionStatusChann)
		if deletionStatus.Status != ExecutionSuccessStatus {
			status.Status = ExecutionErrorStatus
			status.Error = "Failed to delete on diff tasks"
			return status
		}

		// // Updating lastexecution after success on destroying tasks
		// updatedLastExecutionTasks := []commonv1alpha1.TaskExecutionStatus{}
		// for _, p := range lastExecution.Tasks {
		// 	if _, ok := tasksForDeletion[p.Name]; !ok {
		// 		updatedLastExecutionTasks = append(updatedLastExecutionTasks, p)
		// 	}

		// }

		// lastExecution.Tasks = updatedLastExecutionTasks
	}

	dependencyGraphForApply := DependencyGraph{}
	for _, task := range infra.Spec.Tasks {
		dependencyGraphForApply[task.Name] = task.Depends
	}

	pipelineForApply := NewPipeline(e.logger, e.rpcClient, e.terraformProvider)
	return pipelineForApply.Execute(ApplyAction, dependencyGraphForApply, infra, currentExecutionStatusChann)
}

func (e engine) Destroy(infra commonv1alpha1.Infra, currentExecutionStatusChann chan commonv1alpha1.ExecutionStatus) commonv1alpha1.ExecutionStatus {
	pipeline := NewPipeline(e.logger, e.rpcClient, e.terraformProvider)
	dependencyGraphForDestroy := DependencyGraph{}

	for _, task := range infra.Spec.Tasks {
		dependencyGraphForDestroy[task.Name] = []string{}
	}

	for _, task := range infra.Spec.Tasks {
		for _, dep := range task.Depends {
			dependencyGraphForDestroy[dep] = append(dependencyGraphForDestroy[dep], task.Name)
		}
	}

	return pipeline.Execute(DestroyAction, dependencyGraphForDestroy, infra, currentExecutionStatusChann)
}
