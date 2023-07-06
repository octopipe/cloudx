package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
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

	logger.Info("getting last execution")

	lastExecution := newRunnerContext.getLastExecution(executionRef)
	currentExecutionStatusChann := make(chan commonv1alpha1.ExecutionStatus)

	go newRunnerContext.startTimeout(executionRef, lastExecution, currentExecutionStatusChann)

	logger.Info("installing terraform provider")
	terraformProvider, err := terraform.NewTerraformProvider(logger, "")
	if err != nil {
		panic(err)
	}

	doneChann := make(chan bool)
	// go newRunnerContext.setExecutionStatusLive(executionRef, currentExecutionStatusChann, doneChann)

	go func() {
		currentExecution := engine.NewEngine(logger, rpcClient, terraformProvider)
		if action == "APPLY" {
			currentExecutionStatusChann <- currentExecution.Apply(lastExecution, currentSharedInfra, currentExecutionStatusChann)
		} else {
			currentExecution.Destroy(lastExecution, currentSharedInfra, currentExecutionStatusChann)
		}

		doneChann <- true
		// close(currentExecutionStatusChann)
		// close(doneChann)
	}()

	for {
		select {
		case executionStatus := <-currentExecutionStatusChann:
			logger.Info("New status received calling controller...")
			rpcRunnerFinishedArgs := &sharedinfra.RPCSetExecutionStatusArgs{
				Ref:             executionRef,
				ExecutionStatus: executionStatus,
			}

			var reply int
			err := rpcClient.Call("RPCServer.SetExecutionStatus", rpcRunnerFinishedArgs, &reply)
			if err != nil {
				logger.Fatal("Failed to call rpc execution status", zap.Error(err))
			}
		case done := <-doneChann:
			if done {
				logger.Info("Finish engine execution")
				return
			}

		}
	}
}

func (c runnerContext) getDataFromCommandArgs() (commonv1alpha1.SharedInfra, types.NamespacedName, string, error) {
	commandArgs := os.Args[1:]
	rawExecutionRef := commandArgs[1]
	action := commandArgs[0]
	encodedRawSharedInfra := commandArgs[2]

	fmt.Println("ARGS", commandArgs)
	fmt.Println("RAW", commandArgs[2])

	rawDecodedSharedInfra, err := base64.StdEncoding.DecodeString(strings.Trim(encodedRawSharedInfra, "\""))
	if err != nil {
		panic(err)
	}

	var sharedInfra commonv1alpha1.SharedInfra
	err = json.Unmarshal(rawDecodedSharedInfra, &sharedInfra)
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

func (c runnerContext) startTimeout(executionRef types.NamespacedName, lastExecution commonv1alpha1.Execution, currentExecutionStatusChann <-chan commonv1alpha1.ExecutionStatus) {
	time.Sleep(5 * time.Minute)

	lastExecution.Status.Status = engine.ExecutionTimeout
	lastExecution.Status.Error = "Time exceeded"
	rpcRunnerFinishedArgs := &sharedinfra.RPCSetExecutionStatusArgs{
		Ref:             executionRef,
		ExecutionStatus: <-currentExecutionStatusChann,
	}

	var reply int
	err := c.rpcClient.Call("RPCServer.SetExecutionStatus", rpcRunnerFinishedArgs, &reply)
	if err != nil {
		c.logger.Fatal("Error to call controller", zap.Error(err))
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
