package pipeline

import (
	"fmt"
	"sync"
	"time"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/backend"
	"github.com/octopipe/cloudx/internal/backend/terraform"
	"github.com/octopipe/cloudx/internal/customerror"
	"github.com/octopipe/cloudx/internal/rpcclient"
	"github.com/octopipe/cloudx/internal/taskoutput"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
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

type ExecutionOutputItem struct {
	Value     string
	Sensitive bool
	Type      string
}

type ExecutionContext map[string]map[string]ExecutionOutputItem

type ActionFuncType func(taskName string, exectionContext ExecutionContext) (commonv1alpha1.TaskExecutionStatus, map[string]ExecutionOutputItem)

type pipelineCtx struct {
	logger    *zap.Logger
	backend   backend.Backend
	rpcClient rpcclient.Client

	mu               sync.Mutex
	executionContext ExecutionContext
}

type Pipeline interface {
	Start(action string, infra commonv1alpha1.Infra, statusChan chan commonv1alpha1.ExecutionStatus)
}

func NewPipeline(logger *zap.Logger, rpcClient rpcclient.Client, backend backend.Backend) Pipeline {
	return &pipelineCtx{
		logger:           logger,
		backend:          backend,
		rpcClient:        rpcClient,
		executionContext: make(ExecutionContext),
	}
}

func (p *pipelineCtx) Start(action string, infra commonv1alpha1.Infra, statusChan chan commonv1alpha1.ExecutionStatus) {
	if action == "APPLY" {
		tasksForDestroy := p.diffTasksForApply(infra)
		destroyGraph := p.getDestroyGraph(tasksForDestroy)
		applyGraph := p.getApplyGraph(infra)
		p.logger.Info("destroying diff tasks...")
		p.Run(destroyGraph, p.destroy(infra), nil)
		p.logger.Info("apply diff tasks...")
		p.Run(applyGraph, p.apply(infra), statusChan)
	} else {
		tasksForDestroy := p.diffTasksForApply(commonv1alpha1.Infra{})
		destroyGraph := p.getDestroyGraph(tasksForDestroy)
		p.logger.Info("destroying all tasks...")
		p.Run(destroyGraph, p.destroy(infra), statusChan)
	}

}

func (e *pipelineCtx) Run(graph map[string][]string, action ActionFuncType, statusChan chan commonv1alpha1.ExecutionStatus) {
	eg := new(errgroup.Group)
	inDegrees := make(map[string]int)
	status := commonv1alpha1.ExecutionStatus{}

	if len(graph) == 0 {
		e.logger.Info("nothing to execute")
		return
	}

	for node, deps := range graph {
		inDegrees[node] = len(deps)
	}

	go func() {
		time.Sleep(10 * time.Minute)
		e.logger.Info("time limit exceeded")
		statusChan <- commonv1alpha1.ExecutionStatus{
			Status: InfraTimeoutStatus,
			Error: commonv1alpha1.Error{
				Message: "time limit exceeded",
				Code:    "TIME_LIMIT_EXCEEDED",
				Tip:     "Verify if your infrastructure is not stuck in some task.",
			},
		}
	}()

	for {
		for node, deps := range inDegrees {
			if _, ok := e.executionContext[node]; !ok && deps == 0 {
				eg.Go(func(node string) func() error {
					return func() error {

						taskStatus, taskOutput := action(node, e.executionContext)

						e.mu.Lock()
						defer e.mu.Unlock()

						status.Tasks = append(status.Tasks, taskStatus)

						statusChan <- status
						e.executionContext[node] = taskOutput

						if taskStatus.Status == TaskApplyErrorStatus || taskStatus.Status == TaskDestroyErrorStatus {
							return customerror.New(taskStatus.Error.Message, taskStatus.Error.Code, taskStatus.Error.Tip)
						}

						for n, deps := range graph {
							for _, dep := range deps {
								if dep == node {
									inDegrees[n]--
								}
							}
						}

						return nil
					}
				}(node))
			}
		}

		err := eg.Wait()
		if err != nil {
			e.logger.Error("find errors in execution...", zap.Error(err))
			cErr := customerror.Unwrap(err)
			status.Error = commonv1alpha1.Error{
				Message: cErr.Message,
				Code:    cErr.Code,
				Tip:     cErr.Tip,
			}
			status.Status = InfraErrorStatus
			statusChan <- status
			return
		}

		if len(e.executionContext) == len(graph) {
			e.logger.Info("executed all tasks successfully")
			status.Status = InfraSuccessStatus
			statusChan <- status
			return
		}
	}
}

func (p *pipelineCtx) apply(infra commonv1alpha1.Infra) ActionFuncType {
	return func(taskName string, executionContext ExecutionContext) (commonv1alpha1.TaskExecutionStatus, map[string]ExecutionOutputItem) {
		p.logger.Info("applying task", zap.String("name", taskName))
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
			status.Error = commonv1alpha1.Error{
				Message: err.Error(),
				Code:    "TASK_INPUT_INTERPOLATION_ERROR",
				Tip:     "Verify that the task inputs are valid",
			}
			status.Status = TaskApplyErrorStatus
			return status, nil
		}

		status.Inputs = interpolatedInputs
		if currentTask.Backend == backend.TerraformBackend {
			applyInput := terraform.TerraformApplyInput{
				Source:           currentTask.Terraform.Source,
				Version:          currentTask.Terraform.Version,
				TaskInputs:       interpolatedInputs,
				PreviousState:    lastTaskExecutionStatus.Task.State,
				PreviousLockDeps: lastTaskExecutionStatus.Task.DependencyLock,
			}
			result, err := p.backend.Terraform.Apply(applyInput)
			status.FinishedAt = time.Now().Format(time.RFC3339)
			if err != nil {
				status.Error = commonv1alpha1.Error{
					Message: err.Error(),
					Code:    "TASK_APPLY_TERRAFORM_ERROR",
					Tip:     fmt.Sprintf("Verify that the terraform code of task %s is valid", taskName),
				}
				status.Status = TaskApplyErrorStatus
				return status, nil
			}

			status.Task = commonv1alpha1.TaskStatus{
				Terraform:      currentTask.Terraform,
				State:          result.State,
				DependencyLock: result.DependenciesLock,
			}

			outputs := map[string]ExecutionOutputItem{}

			for key, tfMeta := range result.Outputs {
				outputs[key] = ExecutionOutputItem{
					Value:     string(tfMeta.Value),
					Type:      string(tfMeta.Type),
					Sensitive: tfMeta.Sensitive,
				}
			}

			p.logger.Info("creating tasks outputs...")
			err = p.createTaskOutputs(infra, currentTask, outputs)
			if err != nil {
				status.Error = commonv1alpha1.Error{
					Message: err.Error(),
					Code:    "TASK_OUTPUT_CREATION_ERROR",
					Tip:     fmt.Sprintf("An error occurred while creating task outputs for task %s, please retry the execution", taskName),
				}
				status.Status = TaskApplyErrorStatus
				return status, nil
			}

			status.FinishedAt = time.Now().Format(time.RFC3339)
			return status, outputs
		}

		status.Error = commonv1alpha1.Error{
			Message: "invalid task backend",
			Code:    "INVALID_TASK_BACKEND",
			Tip:     "Verify that the task backend is valid",
		}
		status.Status = TaskApplyErrorStatus
		return status, nil
	}
}

func (p pipelineCtx) createTaskOutputs(infra commonv1alpha1.Infra, task commonv1alpha1.InfraTask, outputs map[string]ExecutionOutputItem) error {
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
			Namespace: infra.Namespace,
			Items:     items,
			TaskName:  task.Name,
			InfraRef: commonv1alpha1.Ref{
				Name:      infra.Name,
				Namespace: infra.Namespace,
			},
		}, &reply)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *pipelineCtx) destroy(infra commonv1alpha1.Infra) ActionFuncType {
	return func(taskName string, exectionContext ExecutionContext) (commonv1alpha1.TaskExecutionStatus, map[string]ExecutionOutputItem) {
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
				Source:           lastTaskExecutionStatus.Task.Source,
				Version:          lastTaskExecutionStatus.Task.Version,
				TaskInputs:       lastTaskExecutionStatus.Inputs,
				PreviousState:    lastTaskExecutionStatus.Task.State,
				PreviousLockDeps: lastTaskExecutionStatus.Task.DependencyLock,
			}
			err := p.backend.Terraform.Destroy(destroyInput)
			if err != nil {
				status.Error = commonv1alpha1.Error{
					Message: err.Error(),
					Code:    "TASK_DESTROY_TERRAFORM_ERROR",
				}
				status.Status = TaskDestroyErrorStatus
				return status, nil
			}

			err = p.deleteTaskOutputs(lastTaskExecutionStatus)
			if err != nil {
				status.Error = commonv1alpha1.Error{
					Message: err.Error(),
					Code:    "DESTROY_TASK_OUTPUTS_ERROR",
				}
				status.Status = TaskDestroyErrorStatus
				return status, nil
			}

			return status, nil
		}

		status.Error = commonv1alpha1.Error{
			Message: "invalid task backend",
			Code:    "INVALID_TASK_BACKEND",
			Tip:     "Verify that the task backend is valid",
		}
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

func (p *pipelineCtx) getApplyGraph(infra commonv1alpha1.Infra) map[string][]string {
	dependencyGraphForApply := map[string][]string{}
	for _, task := range infra.Spec.Tasks {
		dependencyGraphForApply[task.Name] = task.Depends
	}

	return dependencyGraphForApply
}

func (p *pipelineCtx) getDestroyGraph(tasksForDestroy map[string]commonv1alpha1.TaskExecutionStatus) map[string][]string {
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
