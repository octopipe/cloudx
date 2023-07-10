package engine

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/connectioninterface"
	"github.com/octopipe/cloudx/internal/rpcclient"
	"github.com/octopipe/cloudx/internal/task"
	"github.com/octopipe/cloudx/internal/terraform"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"k8s.io/apimachinery/pkg/types"
)

const (
	ThisInterpolationOrigin                = "this"
	ConnectionInterfaceInterpolationOrigin = "connection-interface"
)

type DependencyGraph map[string][]string

type ExecutionContextItem struct {
	Value     string
	Sensitive bool
	Type      string
}

type pipeline struct {
	logger            *zap.Logger
	terraformProvider terraform.TerraformProvider
	mu                sync.Mutex
	executionContext  map[string]map[string]ExecutionContextItem
	rpcClient         rpcclient.Client
}

func NewPipeline(logger *zap.Logger, rpcClient rpcclient.Client, terraformProvider terraform.TerraformProvider) pipeline {
	return pipeline{
		logger:            logger,
		terraformProvider: terraformProvider,
		executionContext:  make(map[string]map[string]ExecutionContextItem),
		rpcClient:         rpcClient,
	}
}

func (p *pipeline) Execute(action ExecutionActionType, graph DependencyGraph, infra commonv1alpha1.Infra, currentExecutionStatusChann chan commonv1alpha1.ExecutionStatus) commonv1alpha1.ExecutionStatus {
	lastExecution := infra.Status.LastExecution
	status := commonv1alpha1.ExecutionStatus{
		Status: ExecutionRunningStatus,
		Tasks:  []commonv1alpha1.TaskExecutionStatus{},
	}

	eg := new(errgroup.Group)
	inDegrees := make(map[string]int)

	for node, deps := range graph {
		inDegrees[node] = len(deps)
	}

	for {
		for node, deps := range inDegrees {
			if _, ok := p.executionContext[node]; !ok && deps == 0 {
				eg.Go(func(node string) func() error {
					return func() error {
						p.logger.Info("starting task execution...", zap.String("name", node), zap.Any("action", action))

						taskExecutionStatus, taskOutput := commonv1alpha1.TaskExecutionStatus{}, map[string]ExecutionContextItem{}
						if action == DestroyAction {
							lastTaskExecution := commonv1alpha1.TaskExecutionStatus{}
							for _, statusTask := range lastExecution.Tasks {
								if statusTask.Name == node {
									lastTaskExecution = statusTask
									break
								}
							}
							taskExecutionStatus = p.destroyTask(lastExecution, lastTaskExecution)
						} else {
							currentTask := commonv1alpha1.InfraTask{}
							for _, specTask := range infra.Spec.Tasks {
								if specTask.Name == node {
									currentTask = specTask
									break
								}
							}
							taskExecutionStatus, taskOutput = p.applyTask(lastExecution, currentTask)
						}

						p.mu.Lock()
						defer p.mu.Unlock()

						status.Tasks = append(status.Tasks, taskExecutionStatus)
						if taskExecutionStatus.Status == ExecutionApplyErrorStatus || taskExecutionStatus.Status == ExecutionDestroyErrorStatus {
							status.Status = ExecutionErrorStatus
							status.Error = taskExecutionStatus.Error
							p.logger.Info("task execution failed", zap.String("task-name", taskExecutionStatus.Name), zap.Error(errors.New(taskExecutionStatus.Error)))
							return errors.New(taskExecutionStatus.Error)
						}

						p.executionContext[node] = taskOutput

						for n, deps := range graph {
							for _, dep := range deps {
								if dep == node {
									inDegrees[n]--
								}
							}
						}

						p.logger.Info("finish task execution...", zap.String("name", node), zap.Any("action", action))
						return nil
					}
				}(node))

			}
		}

		err := eg.Wait()
		currentExecutionStatusChann <- status
		if err != nil {
			p.logger.Info("find errors in parallel execution...")
			break
		}

		if len(p.executionContext) == len(graph) {
			break
		}
	}
	status.Status = ExecutionSuccessStatus
	p.logger.Info("finished pipeline execution")
	return status
}

func (p *pipeline) destroyTask(lastExecution commonv1alpha1.ExecutionStatus, lastExecutionTask commonv1alpha1.TaskExecutionStatus) commonv1alpha1.TaskExecutionStatus {
	status := commonv1alpha1.TaskExecutionStatus{
		Name:      lastExecutionTask.Name,
		Ref:       lastExecutionTask.Ref,
		Depends:   lastExecutionTask.Depends,
		Inputs:    lastExecutionTask.Inputs,
		TaskType:  lastExecutionTask.TaskType,
		Status:    ExecutionDestroyed,
		StartedAt: time.Now().Format(time.RFC3339),
	}

	inputs := lastExecutionTask.Inputs
	if lastExecutionTask.TaskType == task.TerraformTaskType {
		err := p.terraformProvider.Destroy(lastExecutionTask.Ref, inputs, lastExecutionTask.State, lastExecutionTask.DependencyLock)
		if err != nil {
			status.Error = err.Error()
			status.Status = ExecutionDestroyErrorStatus
			return status
		}

		return status
	}

	status.Error = "invalid task type"
	status.Status = ExecutionDestroyErrorStatus

	return status
}

func (p *pipeline) applyTask(lastExecution commonv1alpha1.ExecutionStatus, currentTask commonv1alpha1.InfraTask) (commonv1alpha1.TaskExecutionStatus, map[string]ExecutionContextItem) {
	lastTaskExecutionStatus := commonv1alpha1.TaskExecutionStatus{}

	for _, e := range lastExecution.Tasks {
		if e.Name == currentTask.Name {
			lastTaskExecutionStatus = e
		}
	}

	status := commonv1alpha1.TaskExecutionStatus{
		Name:      currentTask.Name,
		Ref:       currentTask.Ref,
		Depends:   currentTask.Depends,
		Inputs:    currentTask.Inputs,
		TaskType:  currentTask.TaskType,
		Status:    ExecutionAppliedStatus,
		StartedAt: time.Now().Format(time.RFC3339),
	}

	inputs, err := p.interpolateTaskInputsByExecutionContext(currentTask)
	if err != nil {
		status.Error = err.Error()
		status.Status = ExecutionApplyErrorStatus
		return status, nil
	}

	status.Inputs = inputs
	if currentTask.TaskType == task.TerraformTaskType {
		out, state, lockfile, err := p.terraformProvider.Apply(currentTask.Ref, inputs, lastTaskExecutionStatus.State, lastTaskExecutionStatus.DependencyLock)
		status.FinishedAt = time.Now().Format(time.RFC3339)
		if err != nil {
			status.Error = err.Error()
			status.Status = ExecutionApplyErrorStatus
			return status, nil
		}

		status.DependencyLock = lockfile
		status.State = state

		outputs := map[string]ExecutionContextItem{}

		for key, tfMeta := range out {
			outputs[key] = ExecutionContextItem{
				Value:     string(tfMeta.Value),
				Type:      string(tfMeta.Type),
				Sensitive: tfMeta.Sensitive,
			}
		}

		return status, outputs
	}

	status.Error = "invalid task type"
	status.Status = ExecutionApplyErrorStatus

	return status, nil
}

func (p *pipeline) interpolateTaskInputsByExecutionContext(task commonv1alpha1.InfraTask) ([]commonv1alpha1.InfraTaskInput, error) {
	inputs := []commonv1alpha1.InfraTaskInput{}
	for _, i := range task.Inputs {
		tokens := Lex(i.Value)
		data := map[string]string{}
		sensitive := false
		for _, t := range tokens {
			if t.Type == TokenVariable {
				s := strings.Split(strings.Trim(t.Value, " "), ".")
				if len(s) != 3 {
					return nil, fmt.Errorf("malformed input variable %s with value %s", i.Key, i.Value)
				}

				value, isSensitive, err := p.getDataByOrigin(s[0], s[1], s[2])
				if err != nil {
					return nil, err
				}

				if isSensitive {
					sensitive = isSensitive
				}

				data[t.Value] = strings.Trim(value, "\"")
			}
		}

		inputs = append(inputs, commonv1alpha1.InfraTaskInput{
			Key:       i.Key,
			Value:     Interpolate(tokens, data),
			Sensitive: sensitive,
		})
	}

	return inputs, nil
}

func (p *pipeline) getDataByOrigin(origin string, name string, attr string) (string, bool, error) {
	switch origin {
	case ThisInterpolationOrigin:
		p.logger.Info("interpolate this origin")
		execution, ok := p.executionContext[name]
		if !ok {
			return "", false, fmt.Errorf("not found task %s in execution context", name)
		}

		executionAttr, ok := execution[attr]
		if !ok {
			return "", false, fmt.Errorf("not found attr %s in finished task execution %s", attr, name)
		}

		return executionAttr.Value, executionAttr.Sensitive, nil

	case ConnectionInterfaceInterpolationOrigin:
		p.logger.Info("interpolate this connection-interface")
		connectionInterface := commonv1alpha1.ConnectionInterface{}
		err := p.rpcClient.Call("ConnectionInterfaceRPCHandler.GetConnectionInterface", connectioninterface.RPCGetConnectionInterfaceArgs{
			Ref: types.NamespacedName{Name: name, Namespace: "default"},
		}, &connectionInterface)
		if err != nil {
			return "", false, err
		}

		for _, out := range connectionInterface.Spec.Outputs {
			if out.Key == attr {
				return out.Value, out.Sensitive, nil
			}
		}

		return "", false, fmt.Errorf("not found attr in connection-interface %s", name)
	default:
		return "", false, fmt.Errorf("invalid origin type %s", origin)
	}
}
