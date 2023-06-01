package main

import (
	"context"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/pluginmanager"
	"github.com/octopipe/cloudx/internal/terraform"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(commonv1alpha1.AddToScheme(scheme))
}

type executionContext struct {
	pluginManager pluginmanager.Manager
	mu            sync.Mutex

	dependencyGraph map[string][]string
	executionGraph  map[string][]string
	executedNodes   map[string][]commonv1alpha1.SharedInfraPluginOutput
}

func main() {
	logger, _ := zap.NewProduction()
	_ = godotenv.Load()

	config := ctrl.GetConfigOrDie()
	k8sClient, err := client.New(config, client.Options{
		Scheme: scheme,
	})
	if err != nil {
		panic(err)
	}

	terraformProvider, err := terraform.NewTerraformProvider(logger)
	if err != nil {
		panic(err)
	}

	pluginManager := pluginmanager.NewPluginManager(logger, terraformProvider)
	currentSharedInfra := &commonv1alpha1.SharedInfra{}
	sharedInfraRef := os.Args[1:]

	namespace, name := "", ""
	s := strings.Split(sharedInfraRef[0], "/")
	if len(s) <= 1 {
		name = s[0]
	} else {
		namespace, name = s[0], s[1]
	}

	req := types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}
	err = k8sClient.Get(context.Background(), req, currentSharedInfra)
	if err != nil {
		logger.Fatal("error on get shared infra", zap.Error(err))
	}

	logger.Info("Starting runner", zap.String("sharedinfra", req.String()))

	dependencyGraph, executionGraph := createGraphs(*currentSharedInfra)
	newExecutionContext := executionContext{
		pluginManager:   pluginManager,
		dependencyGraph: dependencyGraph,
		executionGraph:  executionGraph,
		executedNodes:   map[string][]commonv1alpha1.SharedInfraPluginOutput{},
	}

	_, err = newExecutionContext.execute(currentSharedInfra.Spec.Plugins)
	if err != nil {
		logger.Fatal("error in execution", zap.Error(err))
	}
}

func (c *executionContext) execute(plugins []commonv1alpha1.SharedInfraPlugin) ([]commonv1alpha1.PluginStatus, error) {
	status := []commonv1alpha1.PluginStatus{}
	eg, _ := errgroup.WithContext(context.Background())
	for _, p := range plugins {
		if _, ok := c.executedNodes[p.Name]; !ok && isComplete(c.dependencyGraph[p.Name], c.executedNodes) {
			cb := func(currentPlugin commonv1alpha1.SharedInfraPlugin) func() error {
				return func() error {
					inputs := map[string]interface{}{}

					for _, i := range currentPlugin.Inputs {
						inputs[i.Key] = i.Value
					}

					if p.PluginType == "terraform" {
						out, state, err := c.pluginManager.ExecuteTerraformPlugin(currentPlugin.Ref, inputs)
						if err != nil {
							status = append(status, commonv1alpha1.PluginStatus{
								Name:            p.Name,
								State:           state,
								ExecutionStatus: "Error",
								ExecutionAt:     strconv.Itoa(int(time.Now().Unix())),
								Error:           err.Error(),
							})
							return err
						}

						c.mu.Lock()
						defer c.mu.Unlock()
						status = append(status, commonv1alpha1.PluginStatus{
							Name:            p.Name,
							State:           state,
							ExecutionStatus: "ExecutedWithSuccess",
							ExecutionAt:     strconv.Itoa(int(time.Now().Unix())),
						})
						c.executedNodes[currentPlugin.Name] = out
					}
					return nil
				}
			}
			eg.Go(cb(p))
		}
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	if len(c.executedNodes) == len(c.executionGraph) {
		return status, nil
	}

	p, err := c.execute(plugins)
	if err != nil {
		return nil, err
	}

	status = append(status, p...)
	return status, nil
}

func isComplete(dependencies []string, executedNodes map[string][]commonv1alpha1.SharedInfraPluginOutput) bool {
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
