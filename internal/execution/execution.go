package execution

import (
	"context"
	"sync"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/pluginmanager"
	"golang.org/x/sync/errgroup"
)

type executionContext struct {
	pluginManager pluginmanager.Manager

	mu              sync.Mutex
	dependencyGraph map[string][]string
	executionGraph  map[string][]string
	executedNodes   map[string][]commonv1alpha1.StackSetPluginOutput
}

func NewExecutionManager(pluginManager pluginmanager.Manager, stackset commonv1alpha1.StackSet) error {
	dependencyGraph, executionGraph := createGraphs(stackset)

	ctx := executionContext{
		pluginManager: pluginManager,

		dependencyGraph: dependencyGraph,
		executionGraph:  executionGraph,
		executedNodes:   map[string][]commonv1alpha1.StackSetPluginOutput{},
	}

	return ctx.execute(stackset.Spec.Plugins)
}

func (c *executionContext) execute(plugins []commonv1alpha1.StackSetPlugin) error {
	eg, _ := errgroup.WithContext(context.Background())
	for _, p := range plugins {
		if _, ok := c.executedNodes[p.Name]; !ok && isComplete(c.dependencyGraph[p.Name], c.executedNodes) {
			cb := func(currentPlugin commonv1alpha1.StackSetPlugin) func() error {
				return func() error {
					inputs := map[string]interface{}{}

					for _, i := range currentPlugin.Inputs {
						inputs[i.Key] = i.Value
					}

					out, err := c.pluginManager.Execute(currentPlugin.Ref, inputs)
					if err != nil {
						return err
					}

					c.mu.Lock()
					defer c.mu.Unlock()
					c.executedNodes[currentPlugin.Name] = out
					return nil
				}
			}
			eg.Go(cb(p))
		}
	}

	if err := eg.Wait(); err != nil {
		return err
	}

	if len(c.executedNodes) == len(c.executionGraph) {
		return nil
	}

	return c.execute(plugins)
}

func isComplete(dependencies []string, executedNodes map[string][]commonv1alpha1.StackSetPluginOutput) bool {
	isComplete := true

	for _, d := range dependencies {
		if _, ok := executedNodes[d]; !ok {
			isComplete = false
		}
	}

	return isComplete || len(dependencies) <= 0
}

func createGraphs(stackset commonv1alpha1.StackSet) (map[string][]string, map[string][]string) {

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
