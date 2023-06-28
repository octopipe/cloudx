package main

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/controller/sharedinfra"
	"github.com/octopipe/cloudx/internal/engine"
	"github.com/octopipe/cloudx/internal/rpcclient"
	"github.com/octopipe/cloudx/internal/terraform"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(commonv1alpha1.AddToScheme(scheme))
}

type runnerContext struct {
	logger    *zap.Logger
	rpcClient rpcclient.Client
}

func main() {
	logger, _ := zap.NewProduction()
	_ = godotenv.Load()

	logger.Info("starting runner")

	rpcClient, err := rpcclient.NewRPCClient(os.Getenv("RPC_SERVER_ADDRESS"))
	if err != nil {
		logger.Fatal("Error to connect with controllerr", zap.Error(err), zap.String("address", os.Getenv("RPC_SERVER_ADDRESS")))
	}

	newRunnerContext := runnerContext{
		rpcClient: rpcClient,
		logger:    logger,
	}

	currentSharedInfra, executionRef, action, err := newRunnerContext.getDataFromCommandArgs()
	if err != nil {
		panic(err)
	}

	go newRunnerContext.startTimeout(executionRef)

	logger.Info("getting last execution")

	lastExecution := newRunnerContext.getLastExecution(executionRef)

	logger.Info("installing terraform provider")
	terraformProvider, err := terraform.NewTerraformProvider(logger, "")
	if err != nil {
		panic(err)
	}

	rpcRunnerFinishedArgs := &sharedinfra.RPCSetExecutionStatusArgs{
		Ref: executionRef,
	}

	currentExecution := engine.NewEngine(logger, rpcClient, terraformProvider)
	if action == "APPLY" {
		rpcRunnerFinishedArgs.ExecutionStatus = currentExecution.Apply(lastExecution, currentSharedInfra)
	} else {
		rpcRunnerFinishedArgs.ExecutionStatus = currentExecution.Destroy(lastExecution, currentSharedInfra)
	}

	var reply int
	err = rpcClient.Call("RPCServer.SetExecutionStatus", rpcRunnerFinishedArgs, &reply)
	if err != nil {
		logger.Fatal("Error to call controller", zap.Error(err))
	}

	logger.Info("Finish runner execution")
}

func (c runnerContext) hasErrorsInPluginExecutions(pluginStatus []commonv1alpha1.PluginExecutionStatus) bool {
	for _, p := range pluginStatus {
		if p.Status == engine.ExecutionErrorStatus || p.Status == engine.ExecutionFailedStatus {
			return true
		}
	}

	return false
}

func (c runnerContext) getDataFromCommandArgs() (commonv1alpha1.SharedInfra, types.NamespacedName, string, error) {
	commandArgs := os.Args[1:]
	rawExecutionRef := commandArgs[1]
	action := commandArgs[0]

	var rawSharedInfra string
	err := json.Unmarshal([]byte(commandArgs[2]), &rawSharedInfra)
	if err != nil {
		return commonv1alpha1.SharedInfra{}, types.NamespacedName{}, "", err
	}

	var sharedInfra commonv1alpha1.SharedInfra
	err = json.Unmarshal([]byte(rawSharedInfra), &sharedInfra)
	if err != nil {
		return commonv1alpha1.SharedInfra{}, types.NamespacedName{}, "", err
	}

	executionRef := types.NamespacedName{}

	s := strings.Split(rawExecutionRef, "/")
	if len(s) > 1 {
		executionRef.Namespace = s[0]
		executionRef.Name = s[1]
	} else {
		executionRef.Name = s[0]
		executionRef.Namespace = "default"
	}

	return sharedInfra, executionRef, action, nil
}

func (c runnerContext) startTimeout(executionRef types.NamespacedName) {
	time.Sleep(5 * time.Minute)

	var reply int
	args := &sharedinfra.RPCSetRunnerTimeoutArgs{
		Ref: executionRef,
	}

	err := c.rpcClient.Call("RPCServer.SetRunnerTimeout", args, &reply)
	if err != nil {
		c.logger.Fatal("call rpc timeout error")
	}

	c.logger.Fatal("Runner timeout")
}

func (c runnerContext) getLastExecution(executionRef types.NamespacedName) commonv1alpha1.Execution {
	var reply commonv1alpha1.Execution
	args := &sharedinfra.RPCGetLastExecutionArgs{
		Ref: executionRef,
	}

	err := c.rpcClient.Call("RPCServer.GetLastExecution", args, &reply)
	if err != nil {
		c.logger.Fatal("call rpc get last execution error")
	}

	return reply
}
