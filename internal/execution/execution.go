package execution

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	providerIO "github.com/octopipe/cloudx/internal/io"
	"github.com/octopipe/cloudx/internal/plugin"
	"github.com/octopipe/cloudx/internal/terraform"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	ExecutionSuccessStatus = "SUCCESS"
	ExecutionFailedStatus  = "FAILED"
	ExecutionErrorStatus   = "ERROR"
	ExecutionRunningStatus = "RUNNING"
)

type execution struct {
	logger            *zap.Logger
	terraformProvider terraform.TerraformProvider
	mu                sync.Mutex

	currentSharedInfra commonv1alpha1.SharedInfra
	dependencyGraph    map[string][]string
	executionGraph     map[string][]string
	executedNodes      map[string]providerIO.ProviderOutput
	cb                 func(plugins []commonv1alpha1.PluginStatus) error
}

func NewExecution(logger *zap.Logger, terraformProvider terraform.TerraformProvider, sharedInfra commonv1alpha1.SharedInfra, cb func(plugins []commonv1alpha1.PluginStatus) error) execution {
	dependencyGraph, executionGraph := createGraphs(sharedInfra)
	return execution{
		logger:             logger,
		terraformProvider:  terraformProvider,
		dependencyGraph:    dependencyGraph,
		executionGraph:     executionGraph,
		currentSharedInfra: sharedInfra,
		executedNodes:      map[string]providerIO.ProviderOutput{},
		cb:                 cb,
	}
}

func (c *execution) Destroy() error {
	lastExecution := c.getLastFinishedExecution()

	for _, executionPlugin := range lastExecution.Plugins {
		providerInputs := providerIO.ToProviderInput(executionPlugin.Inputs)
		if executionPlugin.PluginType == plugin.TerraformPluginType {
			err := c.terraformProvider.Destroy(executionPlugin.Ref, providerInputs, executionPlugin.State, executionPlugin.DependencyLock)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *execution) Apply() ([]commonv1alpha1.PluginStatus, error) {
	status := []commonv1alpha1.PluginStatus{}

	err := c.deleteDiffExecutionPlugins()
	if err != nil {
		return nil, err
	}

	for {
		if len(c.executedNodes) == len(c.executionGraph) {
			return status, nil
		}

		s, err := c.execute()
		status = append(status, s...)
		if err != nil {
			return nil, err
		}

	}
}

func (c *execution) execute() ([]commonv1alpha1.PluginStatus, error) {
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

					startedAt := time.Now().Format(time.RFC3339)
					finalInputs, err := c.interpolateInputs(currentPlugin.Inputs)
					if err != nil {
						return err
					}
					c.mu.Lock()
					status = append(status, commonv1alpha1.PluginStatus{
						Name:       currentPlugin.Name,
						Ref:        currentPlugin.Ref,
						PluginType: currentPlugin.PluginType,
						Depends:    currentPlugin.Depends,
						Status:     ExecutionRunningStatus,
						StartedAt:  startedAt,
						Inputs:     finalInputs,
					})
					c.cb(status)
					c.mu.Unlock()
					c.logger.Info("Start plugin execution...", zap.String("name", currentPlugin.Name))
					pluginStatus, pluginOutput, err := c.executeStep(finalInputs, startedAt, currentPlugin)
					c.mu.Lock()
					for i := range status {
						if status[i].Name == pluginStatus.Name {
							status[i] = pluginStatus
						}
					}
					c.cb(status)
					c.mu.Unlock()
					if err != nil {
						return err
					}

					c.mu.Lock()
					defer c.mu.Unlock()
					c.executedNodes[currentPlugin.Name] = pluginOutput
					return nil
				}
			}(p))
		}
	}

	err := eg.Wait()
	return status, err
}

func (c *execution) executeStep(inputs []commonv1alpha1.SharedInfraPluginInput, startedAt string, p commonv1alpha1.SharedInfraPlugin) (commonv1alpha1.PluginStatus, providerIO.ProviderOutput, error) {
	providerInputs := providerIO.ToProviderInput(inputs)
	lastPluginStatus := c.getLastFinishedPluginStatus(p)
	if p.PluginType == plugin.TerraformPluginType {
		out, state, lockFile, err := c.terraformProvider.Apply(p.Ref, providerInputs, lastPluginStatus.State, lastPluginStatus.DependencyLock)
		if err != nil {
			return getPluginStatusError(p, inputs, startedAt, err), providerIO.ProviderOutput{}, nil
		}

		return getTerraformPluginStatusSuccess(p, inputs, startedAt, state, lockFile), out, nil
	}

	return commonv1alpha1.PluginStatus{}, providerIO.ProviderOutput{}, errors.New("invalid plugin type")
}

func (c *execution) manipulateExpression(text string) (string, error) {
	start, end := 0, 0
	for index, c := range text {
		if c == '{' && index > 0 && text[index-1] == '{' {
			start = index + 1
		}

		if c == '}' && index < len(text)-1 && text[index+1] == '}' {
			end = index - 1
		}
	}
	expression := strings.Split(strings.Trim(text[start:end], " "), ".")

	if len(expression) == 3 {

		pluginName, _, outKey := expression[0], expression[1], expression[2]

		out, ok := c.executedNodes[pluginName]
		if !ok {
			return "", fmt.Errorf("plugin %s not found", pluginName)
		}

		outVal, ok := out[outKey]
		if !ok {
			return "", fmt.Errorf("output %s in plugin %s not found", outKey, pluginName)
		}

		reg := regexp.MustCompile(`"([^"]*)"`)
		newExpression := fmt.Sprintf("%s%s%s", text[0:start-2], reg.ReplaceAllString(outVal.Value, "${1}"), text[end+3:])

		return newExpression, nil
	}

	return text, nil

}

func (c *execution) interpolateInputs(inputs []commonv1alpha1.SharedInfraPluginInput) ([]commonv1alpha1.SharedInfraPluginInput, error) {
	newInputs := []commonv1alpha1.SharedInfraPluginInput{}
	for _, i := range inputs {

		inputParsed, err := c.manipulateExpression(i.Value)
		if err != nil {
			return nil, err
		}

		newInputs = append(newInputs, commonv1alpha1.SharedInfraPluginInput{
			Key:   i.Key,
			Value: inputParsed,
		})
	}

	return newInputs, nil
}

// TODO: A more inteligent way to do the diff-deleting plugins
func (c *execution) deleteDiffExecutionPlugins() error {
	lastExecution := c.getLastFinishedExecution()

	for _, executionPlugin := range lastExecution.Plugins {
		foundPlugin := false
		for _, currentPlugin := range c.currentSharedInfra.Spec.Plugins {
			if executionPlugin.Name == currentPlugin.Name {
				foundPlugin = true
				break
			}
		}

		if !foundPlugin {
			c.logger.Info("Not found plugin in current shared infra, deleting plugin", zap.String("plugin", executionPlugin.Name))
			providerInputs := providerIO.ToProviderInput(executionPlugin.Inputs)
			if executionPlugin.PluginType == plugin.TerraformPluginType {
				err := c.terraformProvider.Destroy(executionPlugin.Ref, providerInputs, executionPlugin.State, executionPlugin.DependencyLock)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (c *execution) getLastFinishedExecution() commonv1alpha1.SharedInfraExecutionStatus {
	for _, e := range c.currentSharedInfra.Status.Executions {
		if e.Status != "RUNNING" {
			return e
		}
	}

	return commonv1alpha1.SharedInfraExecutionStatus{}
}

func (c *execution) getLastFinishedPluginStatus(currentPlugin commonv1alpha1.SharedInfraPlugin) commonv1alpha1.PluginStatus {
	lastExecution := c.getLastFinishedExecution()

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

func getTerraformPluginStatusSuccess(plugin commonv1alpha1.SharedInfraPlugin, inputs []commonv1alpha1.SharedInfraPluginInput, startedAt string, state string, lockFile string) commonv1alpha1.PluginStatus {
	escapedState, err := json.Marshal(state)
	if err != nil {
		return commonv1alpha1.PluginStatus{
			Name:       plugin.Name,
			Ref:        plugin.Ref,
			Inputs:     inputs,
			Depends:    plugin.Depends,
			PluginType: plugin.PluginType,
			Status:     ExecutionErrorStatus,
			StartedAt:  startedAt,
			FinishedAt: time.Now().Format(time.RFC3339),
			Error:      err.Error(),
		}
	}

	escapedLockFile, err := json.Marshal(lockFile)
	if err != nil {
		return commonv1alpha1.PluginStatus{
			Name:       plugin.Name,
			Ref:        plugin.Ref,
			Depends:    plugin.Depends,
			Inputs:     inputs,
			PluginType: plugin.PluginType,
			Status:     ExecutionErrorStatus,
			StartedAt:  startedAt,
			FinishedAt: time.Now().Format(time.RFC3339),
			Error:      err.Error(),
		}
	}

	return commonv1alpha1.PluginStatus{
		Name:           plugin.Name,
		Ref:            plugin.Ref,
		Inputs:         inputs,
		Depends:        plugin.Depends,
		PluginType:     plugin.PluginType,
		State:          string(escapedState),
		DependencyLock: string(escapedLockFile),
		Status:         ExecutionSuccessStatus,
		StartedAt:      startedAt,
		FinishedAt:     time.Now().Format(time.RFC3339),
	}
}

func getPluginStatusError(plugin commonv1alpha1.SharedInfraPlugin, inputs []commonv1alpha1.SharedInfraPluginInput, startedAt string, err error) commonv1alpha1.PluginStatus {
	escapedError, err := json.Marshal(err.Error())
	if err != nil {
		return commonv1alpha1.PluginStatus{
			Name:       plugin.Name,
			Ref:        plugin.Ref,
			Inputs:     inputs,
			Depends:    plugin.Depends,
			PluginType: plugin.PluginType,
			Status:     ExecutionErrorStatus,
			StartedAt:  startedAt,
			FinishedAt: time.Now().Format(time.RFC3339),
			Error:      err.Error(),
		}
	}

	return commonv1alpha1.PluginStatus{
		Name:       plugin.Name,
		Ref:        plugin.Ref,
		Inputs:     inputs,
		Depends:    plugin.Depends,
		PluginType: plugin.PluginType,
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
