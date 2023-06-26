package engine

import (
	"context"
	"fmt"
	"sync"
	"time"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/terraform"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type DependencyGraph map[string][]string

type pipeline struct {
	logger            *zap.Logger
	terraformProvider terraform.TerraformProvider
	mu                sync.Mutex
}

func NewPipeline(terraformProvider terraform.TerraformProvider) pipeline {
	return pipeline{
		terraformProvider: terraformProvider,
	}
}

func (p *pipeline) Execute(sharedInfra commonv1alpha1.SharedInfra) {
	eg, _ := errgroup.WithContext(context.Background())
	graph := createDependencyGraph(sharedInfra)
	inDegrees := make(map[string]int)
	queue := make(chan string)
	result := make(chan string)

	for node := range graph {
		inDegrees[node] = 0
	}

	for _, dependencies := range graph {
		for _, dependency := range dependencies {
			inDegrees[dependency]++
		}
	}

	for {
		select {
		case item := <-queue:

		}
	}

	for node := range graph {
		if inDegrees[node] == 0 {
			eg.Go(func() error {
				pluginExecutionStatus, err := p.processNode(node, graph, sharedInfra, inDegrees, queue)
				if err != nil {

				}
			})
		}
	}

	if err := eg.Wait(); err != nil {
		return
	}

	var executionOrder []string
	for node := range result {
		executionOrder = append(executionOrder, node)
	}

	return executionOrder
}

func (p *pipeline) processNode(node string, graph DependencyGraph, sharedInfra commonv1alpha1.SharedInfra, inDegrees map[string]int, queue chan<- string) (commonv1alpha1.PluginExecutionStatus, error) {

	// for _, plugin := range sharedInfra.Spec.Plugins {
	// 	if plugin.Name == node {
	// 		p.terraformProvider.Apply()
	// 	}
	// }

	fmt.Println(node)

	time.Sleep(3 * time.Second)

	if dependencies, ok := graph[node]; ok {
		for _, dependency := range dependencies {
			if inDegrees[dependency] > 0 {
				inDegrees[dependency]--
				if inDegrees[dependency] == 0 {
					queue <- dependency
				}
			}
		}
	}

	return commonv1alpha1.PluginExecutionStatus{}, nil
}

func createDependencyGraph(sharedInfra commonv1alpha1.SharedInfra) DependencyGraph {
	dependencyGraph := DependencyGraph{}

	for _, plugin := range sharedInfra.Spec.Plugins {
		dependencyGraph[plugin.Name] = plugin.Depends
	}

	return dependencyGraph
}
