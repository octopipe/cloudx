package execution

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	providerIO "github.com/octopipe/cloudx/internal/provider/io"
	"github.com/octopipe/cloudx/internal/provider/terraform"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
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

					pluginStatus, pluginOutput := c.executeStep(currentPlugin)
					status = append(status, pluginStatus)
					if pluginStatus.Error != "" {
						return errors.New(pluginStatus.Error)
					}

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

func (c *execution) executeStep(p commonv1alpha1.SharedInfraPlugin) (commonv1alpha1.PluginStatus, providerIO.ProviderOutput) {
	startedAt := time.Now().Format(time.RFC3339)
	providerInputs := providerIO.ToProviderInput(p.Inputs)
	if p.PluginType == "terraform" {
		out, state, err := c.terraformProvider.Apply(p.Ref, providerInputs)
		if err != nil {
			escapedErrorMsg, err := json.Marshal(err.Error())
			if err != nil {
				c.logger.Error("error to escape plugin error", zap.Error(err))
				return commonv1alpha1.PluginStatus{
					Status:     "ERROR",
					StartedAt:  startedAt,
					FinishedAt: time.Now().Format(time.RFC3339),
					Error:      err.Error(),
				}, providerIO.ProviderOutput{}
			}

			return commonv1alpha1.PluginStatus{
				Name:       p.Name,
				Status:     "ERROR",
				StartedAt:  startedAt,
				FinishedAt: time.Now().Format(time.RFC3339),
				Error:      string(escapedErrorMsg),
			}, providerIO.ProviderOutput{}
		}

		escapedState, err := json.Marshal(state)
		if err != nil {
			return commonv1alpha1.PluginStatus{
				Status:     "ERROR",
				StartedAt:  startedAt,
				FinishedAt: time.Now().Format(time.RFC3339),
				Error:      err.Error(),
			}, providerIO.ProviderOutput{}
		}

		return commonv1alpha1.PluginStatus{
			Name:       p.Name,
			State:      string(escapedState),
			Status:     "SUCCESS",
			StartedAt:  startedAt,
			FinishedAt: time.Now().Format(time.RFC3339),
		}, out
	}

	return commonv1alpha1.PluginStatus{}, providerIO.ProviderOutput{}
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
