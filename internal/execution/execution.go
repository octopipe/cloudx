package execution

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/plugin"
	providerIO "github.com/octopipe/cloudx/internal/provider/io"
	"github.com/octopipe/cloudx/internal/provider/terraform"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	ExecutionSuccessStatus = "SUCCESS"
	ExecutionFailedStatus  = "FAILED"
	ExecutionErrorStatus   = "ERRORS"
)

type execution struct {
	logger            *zap.Logger
	terraformProvider terraform.TerraformProvider
	mu                sync.Mutex

	currentSharedInfra commonv1alpha1.SharedInfra
	dependencyGraph    map[string][]string
	executionGraph     map[string][]string
	executedNodes      map[string]providerIO.ProviderOutput
}

func NewExecution(logger *zap.Logger, terraformProvider terraform.TerraformProvider, sharedInfra commonv1alpha1.SharedInfra) execution {
	dependencyGraph, executionGraph := createGraphs(sharedInfra)
	return execution{
		logger:             logger,
		terraformProvider:  terraformProvider,
		dependencyGraph:    dependencyGraph,
		executionGraph:     executionGraph,
		currentSharedInfra: sharedInfra,
		executedNodes:      map[string]providerIO.ProviderOutput{},
	}
}

func (c *execution) Start() ([]commonv1alpha1.PluginStatus, error) {
	status := []commonv1alpha1.PluginStatus{}
	eg, _ := errgroup.WithContext(context.Background())
	for _, p := range c.currentSharedInfra.Spec.Plugins {
		if _, ok := c.executedNodes[p.Name]; !ok && isComplete(c.dependencyGraph[p.Name], c.executedNodes) {
			eg.Go(func(currentPlugin commonv1alpha1.SharedInfraPlugin) func() error {
				return func() error {
					inputs := map[string]interface{}{}

					for _, i := range currentPlugin.Inputs {
						inputs[i.Key] = i.Value
					}

					pluginStatus, pluginOutput, err := c.executeStep(currentPlugin)
					status = append(status, pluginStatus)
					if err != nil {
						return err
					}

					if pluginStatus.Error != "" {
						return errors.New(pluginStatus.Error)
					}

					c.mu.Lock()
					defer c.mu.Unlock()
					c.executedNodes[currentPlugin.Name] = pluginOutput
					return nil
				}
			}(p))
		}
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	if len(c.executedNodes) == len(c.executionGraph) {
		return status, nil
	}

	p, err := c.Start()
	if err != nil {
		return nil, err
	}

	status = append(status, p...)
	return status, nil
}

func (c *execution) executeStep(p commonv1alpha1.SharedInfraPlugin) (commonv1alpha1.PluginStatus, providerIO.ProviderOutput, error) {
	startedAt := time.Now().Format(time.RFC3339)
	providerInputs := providerIO.ToProviderInput(p.Inputs)
	lastPluginStatus := c.getLastPluginStatus(p)
	if p.PluginType == plugin.TerraformPluginType {
		out, state, err := c.terraformProvider.Apply(p.Ref, providerInputs, lastPluginStatus.State, lastPluginStatus.DependencyLock)
		if err != nil {
			return getPluginStatusError(p.Name, startedAt, err), providerIO.ProviderOutput{}, nil
		}

		return getTerraformPluginStatusSuccess(p.Name, startedAt, state), out, nil
	}

	return commonv1alpha1.PluginStatus{}, providerIO.ProviderOutput{}, errors.New("invalid plugin type")
}

func (c *execution) getLastExecution() commonv1alpha1.SharedInfraExecutionStatus {
	lastExecution := commonv1alpha1.SharedInfraExecutionStatus{}

	if len(c.currentSharedInfra.Status.Executions) > 0 {
		lastExecution = c.currentSharedInfra.Status.Executions[0]
	}

	return lastExecution
}

func (c *execution) getLastPluginStatus(currentPlugin commonv1alpha1.SharedInfraPlugin) commonv1alpha1.PluginStatus {
	lastExecution := c.getLastExecution()

	if len(lastExecution.Plugins) <= 0 {
		return commonv1alpha1.PluginStatus{}
	}

	for _, p := range lastExecution.Plugins {
		if p.Name == currentPlugin.Name {
			return p
		}
	}

	return commonv1alpha1.PluginStatus{}
}

func getTerraformPluginStatusSuccess(name string, startedAt string, state string) commonv1alpha1.PluginStatus {
	escapedState, err := json.Marshal(state)
	if err != nil {
		return commonv1alpha1.PluginStatus{
			Name:       name,
			Status:     ExecutionErrorStatus,
			StartedAt:  startedAt,
			FinishedAt: time.Now().Format(time.RFC3339),
			Error:      err.Error(),
		}
	}

	return commonv1alpha1.PluginStatus{
		Name:       name,
		State:      string(escapedState),
		Status:     ExecutionSuccessStatus,
		StartedAt:  startedAt,
		FinishedAt: time.Now().Format(time.RFC3339),
	}
}

func getPluginStatusError(name string, startedAt string, err error) commonv1alpha1.PluginStatus {
	escapedError, err := json.Marshal(err.Error())
	if err != nil {
		return commonv1alpha1.PluginStatus{
			Name:       name,
			Status:     ExecutionErrorStatus,
			StartedAt:  startedAt,
			FinishedAt: time.Now().Format(time.RFC3339),
			Error:      err.Error(),
		}
	}

	return commonv1alpha1.PluginStatus{
		Name:       name,
		Status:     ExecutionErrorStatus,
		StartedAt:  startedAt,
		FinishedAt: time.Now().Format(time.RFC3339),
		Error:      string(escapedError),
	}
}

func isComplete(dependencies []string, executedNodes map[string]providerIO.ProviderOutput) bool {
	isComplete := true

	for _, d := range dependencies {
		if _, ok := executedNodes[d]; !ok {
			isComplete = false
		}
	}

	return isComplete || len(dependencies) <= 0
}

func createGraphs(stackset commonv1alpha1.SharedInfra) (map[string][]string, map[string][]string) {
	dependencyGraph := map[string][]string{}
	executionGraph := map[string][]string{}

	for _, p := range stackset.Spec.Plugins {
		dependencyGraph[p.Name] = p.Depends
		executionGraph[p.Name] = []string{}
	}

	for _, p := range stackset.Spec.Plugins {
		for _, d := range p.Depends {
			executionGraph[d] = append(executionGraph[d], p.Name)
		}
	}

	return dependencyGraph, executionGraph
}
