package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/rpc"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/controller/sharedinfra"
	providerIO "github.com/octopipe/cloudx/internal/provider/io"
	"github.com/octopipe/cloudx/internal/provider/terraform"
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
	logger            *zap.Logger
	terraformProvider terraform.TerraformProvider
	mu                sync.Mutex
	rpcClient         *rpc.Client

	dependencyGraph map[string][]string
	executionGraph  map[string][]string
	executedNodes   map[string]providerIO.ProviderOutput
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

	rpcClient, err := rpc.DialHTTP("tcp", os.Getenv("RPC_SERVER_ADDRESS"))
	if err != nil {
		logger.Fatal("Error to connect with controllerr", zap.Error(err), zap.String("address", os.Getenv("RPC_SERVER_ADDRESS")))
	}

	terraformProvider, err := terraform.NewTerraformProvider(logger)
	if err != nil {
		panic(err)
	}

	currentSharedInfra := &commonv1alpha1.SharedInfra{}
	commandArgs := os.Args[1:]
	executionId := commandArgs[1]

	namespace, name := "", ""
	s := strings.Split(commandArgs[0], "/")
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

	go func() {
		time.Sleep(5 * time.Minute)

		var reply int
		args := &sharedinfra.RPCSetRunnerTimeoutArgs{
			SharedInfraRef: req,
			ExecutionId:    executionId,
		}

		err := rpcClient.Call("RPCServer.SetRunnerTimeout", args, &reply)
		if err != nil {
			logger.Fatal("call rpc timeout error")
		}

		logger.Fatal("Runner timeout")
	}()

	dependencyGraph, executionGraph := createGraphs(*currentSharedInfra)
	newExecutionContext := executionContext{
		logger:            logger,
		rpcClient:         rpcClient,
		terraformProvider: terraformProvider,
		dependencyGraph:   dependencyGraph,
		executionGraph:    executionGraph,
		executedNodes:     map[string]providerIO.ProviderOutput{},
	}

	errMsg := ""
	status := "Success"
	pluginStatus, err := newExecutionContext.execute(currentSharedInfra.Spec.Plugins)
	if err != nil {
		errMsg = err.Error()
		status = "Error"
	}

	args := &sharedinfra.RPCSetRunnerFinishedArgs{
		ExecutionId: executionId,
		Ref:         req,
		Error:       errMsg,
		Plugins:     pluginStatus,
		Status:      status,
		FinishedAt:  time.Now().Format(time.RFC3339),
	}

	var reply int
	err = rpcClient.Call("RPCServer.SetRunnerFinished", args, &reply)
	if err != nil {
		logger.Fatal("Error to call controller", zap.Error(err))
	}

	logger.Info("Finish runner execution")
}

func (c *executionContext) execute(plugins []commonv1alpha1.SharedInfraPlugin) ([]commonv1alpha1.PluginStatus, error) {
	status := []commonv1alpha1.PluginStatus{}
	eg, _ := errgroup.WithContext(context.Background())
	for _, p := range plugins {
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

	p, err := c.execute(plugins)
	if err != nil {
		return nil, err
	}

	status = append(status, p...)
	return status, nil
}

func (c *executionContext) executeStep(p commonv1alpha1.SharedInfraPlugin) (commonv1alpha1.PluginStatus, providerIO.ProviderOutput) {
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
