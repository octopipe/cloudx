package engine

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
	ExecutionTimeout       = "TIMEOUT"
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

func (c *execution) Destroy(lastExecution commonv1alpha1.Execution) error {
	for _, executionPlugin := range lastExecution.Status.Plugins {
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

func (c *execution) Apply(lastExecution commonv1alpha1.Execution) ([]commonv1alpha1.PluginExecutionStatus, error) {
	status := []commonv1alpha1.PluginExecutionStatus{}

	err := c.deleteDiffExecutionPlugins(lastExecution)
	if err != nil {
		return nil, err
	}

	for {
		if len(c.executedNodes) == len(c.executionGraph) {
			return status, nil
		}

		s, err := c.execute(lastExecution)
		status = append(status, s...)
		if err != nil {
			return status, err
		}
	}
}

func (c *execution) execute(lastExecution commonv1alpha1.Execution) ([]commonv1alpha1.PluginExecutionStatus, error) {
	status := []commonv1alpha1.PluginExecutionStatus{}
	eg, _ := errgroup.WithContext(context.Background())

	for _, p := range c.currentSharedInfra.Spec.Plugins {
		if _, ok := c.executedNodes[p.Name]; !ok && isComplete(c.dependencyGraph[p.Name], c.executedNodes) {
			eg.Go(func(currentPlugin commonv1alpha1.SharedInfraPlugin) func() error {
				return func() error {
					c.logger.Info("Start plugin execution...", zap.String("name", currentPlugin.Name))
					c.logger.Info("Resolve inputs...", zap.String("name", currentPlugin.Name))
					startedAt := time.Now().Format(time.RFC3339)
					finalInputs, err := c.interpolateInputs(currentPlugin.Inputs)
					if err != nil {
						return err
					}

					pluginExecutionStatus, pluginOutput, err := c.executeStep(lastExecution, finalInputs, startedAt, currentPlugin)
					c.mu.Lock()
					status = append(status, pluginExecutionStatus)
					c.mu.Unlock()

					if err != nil {
						return err
					}

					if pluginExecutionStatus.Error != "" {
						c.logger.Error("Plugin execution failed", zap.String("name", currentPlugin.Name), zap.String("error", pluginExecutionStatus.Error))
						return errors.New(pluginExecutionStatus.Error)
					}

					c.logger.Info("Finish plugin execution...", zap.String("name", currentPlugin.Name))

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

func (c *execution) executeStep(lastExecution commonv1alpha1.Execution, inputs []commonv1alpha1.SharedInfraPluginInput, startedAt string, p commonv1alpha1.SharedInfraPlugin) (commonv1alpha1.PluginExecutionStatus, providerIO.ProviderOutput, error) {
	providerInputs := providerIO.ToProviderInput(inputs)
	lastPluginExecutionStatus := commonv1alpha1.PluginExecutionStatus{}

	for _, e := range lastExecution.Status.Plugins {
		if e.Name == p.Name {
			lastPluginExecutionStatus = e
		}
	}

	if p.PluginType == plugin.TerraformPluginType {
		out, state, lockFile, err := c.terraformProvider.Apply(p.Ref, providerInputs, lastPluginExecutionStatus.State, lastPluginExecutionStatus.DependencyLock)
		if err != nil {
			return getPluginExecutionStatusError(p, inputs, startedAt, err), providerIO.ProviderOutput{}, nil
		}

		return getTerraformPluginExecutionStatusSuccess(p, inputs, startedAt, state, lockFile), out, nil
	}

	return commonv1alpha1.PluginExecutionStatus{}, providerIO.ProviderOutput{}, errors.New("invalid plugin type")
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
func (c *execution) deleteDiffExecutionPlugins(lastExecution commonv1alpha1.Execution) error {
	for _, executionPlugin := range lastExecution.Status.Plugins {
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

func getTerraformPluginExecutionStatusSuccess(plugin commonv1alpha1.SharedInfraPlugin, inputs []commonv1alpha1.SharedInfraPluginInput, startedAt string, state string, lockFile string) commonv1alpha1.PluginExecutionStatus {
	escapedState, err := json.Marshal(state)
	if err != nil {
		return commonv1alpha1.PluginExecutionStatus{
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
		return commonv1alpha1.PluginExecutionStatus{
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

	return commonv1alpha1.PluginExecutionStatus{
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

func getPluginExecutionStatusError(plugin commonv1alpha1.SharedInfraPlugin, inputs []commonv1alpha1.SharedInfraPluginInput, startedAt string, err error) commonv1alpha1.PluginExecutionStatus {
	escapedError, err := json.Marshal(err.Error())
	if err != nil {
		return commonv1alpha1.PluginExecutionStatus{
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

	return commonv1alpha1.PluginExecutionStatus{
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
