package pipeline

import (
	"time"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/backend"
	"github.com/octopipe/cloudx/internal/backend/terraform"
	"github.com/octopipe/cloudx/internal/engine"
	"github.com/octopipe/cloudx/internal/rpcclient"
	"github.com/octopipe/cloudx/internal/taskoutput"
	"go.uber.org/zap"
)

const (
	InfraSuccessStatus = "SUCCESS"
	InfraErrorStatus   = "ERROR"
	InfraRunningStatus = "RUNNING"
	InfraTimeoutStatus = "TIMEOUT"
)

const (
	TaskAppliedStatus      = "APPLIED"
	TaskApplyErrorStatus   = "APPLY_ERROR"
	TaskDestroyed          = "DESTROYED"
	TaskDestroyErrorStatus = "DESTROY_ERROR"
)

type pipelineCtx struct {
	logger    *zap.Logger
	backend   backend.Backend
	engine    *engine.Engine
	rpcClient rpcclient.Client
}

type Pipeline interface {
	Start(action string, infra commonv1alpha1.Infra, taskStatusChan chan commonv1alpha1.TaskExecutionStatus) commonv1alpha1.ExecutionStatus
}

func NewPipeline(logger *zap.Logger, rpcClient rpcclient.Client, backend backend.Backend, engine *engine.Engine) Pipeline {
	return &pipelineCtx{
		logger:    logger,
		backend:   backend,
		engine:    engine,
		rpcClient: rpcClient,
	}
}

func (p pipelineCtx) Start(action string, infra commonv1alpha1.Infra, taskStatusChan chan commonv1alpha1.TaskExecutionStatus) commonv1alpha1.ExecutionStatus {
	status := commonv1alpha1.ExecutionStatus{}
	if action == "APPLY" {
		tasksForDestroy := p.diffTasksForApply(infra)
		destroyGraph := p.getDestroyGraph(tasksForDestroy)
		applyGraph := p.getApplyGraph(infra)
		p.engine.Run(destroyGraph, p.destroy(infra), taskStatusChan)
		p.engine.Run(applyGraph, p.apply(infra), taskStatusChan)
	} else {
		tasksForDestroy := p.diffTasksForApply(commonv1alpha1.Infra{})
		destroyGraph := p.getDestroyGraph(tasksForDestroy)
		p.engine.Run(destroyGraph, p.destroy(infra), taskStatusChan)
	}

	return status
}

func (p pipelineCtx) apply(infra commonv1alpha1.Infra) engine.ActionFuncType {
	return func(taskName string, executionContext engine.ExecutionContext) (commonv1alpha1.TaskExecutionStatus, map[string]engine.ExecutionOutputItem) {
		lastTaskExecutionStatus := commonv1alpha1.TaskExecutionStatus{}
		currentTask := commonv1alpha1.InfraTask{}
		for _, specTask := range infra.Spec.Tasks {
			if specTask.Name == taskName {
				currentTask = specTask
				break
			}
		}
		for _, e := range infra.Status.LastExecution.Tasks {
			if e.Name == currentTask.Name {
				lastTaskExecutionStatus = e
			}
		}

		status := commonv1alpha1.TaskExecutionStatus{
			Name:        currentTask.Name,
			Depends:     currentTask.Depends,
			Inputs:      currentTask.Inputs,
			Backend:     currentTask.Backend,
			TaskOutputs: currentTask.TaskOutputs,
			Status:      TaskAppliedStatus,
			StartedAt:   time.Now().Format(time.RFC3339),
		}

		interpolatedInputs, err := p.interpolateTaskInputsByExecutionContext(currentTask, executionContext)
		if err != nil {
			status.Error = err.Error()
			status.Status = TaskApplyErrorStatus
			return status, nil
		}

		status.Inputs = interpolatedInputs
		if currentTask.Backend == backend.TerraformBackend {
			applyInput := terraform.TerraformApplyInput{
				Source:           currentTask.Terraform.Source,
				Version:          currentTask.Terraform.Version,
				TaskInputs:       interpolatedInputs,
				PreviousState:    lastTaskExecutionStatus.Terraform.State,
				PreviousLockDeps: lastTaskExecutionStatus.Terraform.DependencyLock,
			}
			result, err := p.backend.Terraform.Apply(applyInput)
			status.FinishedAt = time.Now().Format(time.RFC3339)
			if err != nil {
				status.Error = err.Error()
				status.Status = TaskApplyErrorStatus
				return status, nil
			}

			status.Terraform = commonv1alpha1.TerraformTaskStatus{
				TerraformTask:  currentTask.Terraform,
				State:          result.State,
				DependencyLock: result.DependenciesLock,
			}

			outputs := map[string]engine.ExecutionOutputItem{}

			for key, tfMeta := range result.Outputs {
				outputs[key] = engine.ExecutionOutputItem{
					Value:     string(tfMeta.Value),
					Type:      string(tfMeta.Type),
					Sensitive: tfMeta.Sensitive,
				}
			}

			p.logger.Info("creating tasks outputs...")
			err = p.createTaskOutputs(currentTask, outputs)
			if err != nil {
				status.Error = err.Error()
				status.Status = TaskApplyErrorStatus
				return status, nil
			}

			return status, outputs
		}

		status.Error = "invalid task type"
		status.Status = TaskApplyErrorStatus
		return status, nil
	}
}

func (p pipelineCtx) createTaskOutputs(task commonv1alpha1.InfraTask, outputs map[string]engine.ExecutionOutputItem) error {
	for _, t := range task.TaskOutputs {
		p.logger.Info("creating task output", zap.String("name", t.Name))
		items := []taskoutput.RPCCreateTaskOutputItem{}
		for key, value := range outputs {
			items = append(items, taskoutput.RPCCreateTaskOutputItem{
				InfraTaskOutputItem: commonv1alpha1.InfraTaskOutputItem{
					Key:       key,
					Sensitive: value.Sensitive,
				},
				Value: value.Value,
			})
		}

		var reply int
		err := p.rpcClient.Call("TaskOutputRPCHandler.ApplyTaskOuput", taskoutput.RPCCreateTaskOutputArgs{
			Name:      t.Name,
			Namespace: "default",
			Items:     items,
		}, &reply)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p pipelineCtx) destroy(infra commonv1alpha1.Infra) engine.ActionFuncType {
	return func(taskName string, exectionContext engine.ExecutionContext) (commonv1alpha1.TaskExecutionStatus, map[string]engine.ExecutionOutputItem) {
		lastTaskExecutionStatus := commonv1alpha1.TaskExecutionStatus{}
		for _, e := range infra.Status.LastExecution.Tasks {
			if e.Name == taskName {
				lastTaskExecutionStatus = e
			}
		}
		status := commonv1alpha1.TaskExecutionStatus{
			Name:        lastTaskExecutionStatus.Name,
			Depends:     lastTaskExecutionStatus.Depends,
			Inputs:      lastTaskExecutionStatus.Inputs,
			TaskOutputs: lastTaskExecutionStatus.TaskOutputs,
			Status:      TaskDestroyed,
			StartedAt:   time.Now().Format(time.RFC3339),
		}

		if lastTaskExecutionStatus.Backend == backend.TerraformBackend {
			destroyInput := terraform.TerraformDestroyInput{
				Source:           lastTaskExecutionStatus.Terraform.Source,
				Version:          lastTaskExecutionStatus.Terraform.Version,
				TaskInputs:       lastTaskExecutionStatus.Inputs,
				PreviousState:    lastTaskExecutionStatus.Terraform.State,
				PreviousLockDeps: lastTaskExecutionStatus.Terraform.DependencyLock,
			}
			err := p.backend.Terraform.Destroy(destroyInput)
			if err != nil {
				status.Error = err.Error()
				status.Status = TaskDestroyErrorStatus
				return status, nil
			}

			err = p.deleteTaskOutputs(lastTaskExecutionStatus)
			if err != nil {
				status.Error = err.Error()
				status.Status = TaskDestroyErrorStatus
				return status, nil
			}

			return status, nil
		}

		status.Error = "invalid task type"
		status.Status = TaskDestroyErrorStatus

		return status, nil
	}
}

func (p pipelineCtx) deleteTaskOutputs(task commonv1alpha1.TaskExecutionStatus) error {
	var reply int
	for _, t := range task.TaskOutputs {
		err := p.rpcClient.Call("TaskOutputRPCHandler.DeleteTaskOutput", taskoutput.RPCCreateTaskOutputArgs{
			Name:      t.Name,
			Namespace: "default",
		}, &reply)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e pipelineCtx) diffTasksForApply(infra commonv1alpha1.Infra) map[string]commonv1alpha1.TaskExecutionStatus {
	forDeletion := map[string]commonv1alpha1.TaskExecutionStatus{}
	for _, lastTaskExecution := range infra.Status.LastExecution.Tasks {
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

func (p pipelineCtx) getApplyGraph(infra commonv1alpha1.Infra) map[string][]string {
	dependencyGraphForApply := map[string][]string{}
	for _, task := range infra.Spec.Tasks {
		dependencyGraphForApply[task.Name] = task.Depends
	}

	return dependencyGraphForApply
}

func (p pipelineCtx) getDestroyGraph(tasksForDestroy map[string]commonv1alpha1.TaskExecutionStatus) map[string][]string {
	dependencyGraphForDeletion := map[string][]string{}
	for task := range tasksForDestroy {
		dependencyGraphForDeletion[task] = []string{}
	}

	for task, execution := range tasksForDestroy {
		for _, dep := range execution.Depends {
			if _, ok := tasksForDestroy[dep]; ok {
				dependencyGraphForDeletion[dep] = append(dependencyGraphForDeletion[dep], task)
			}
		}
	}

	return dependencyGraphForDeletion
}
