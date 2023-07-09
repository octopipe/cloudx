package main

import (
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

	sharedInfraRef, action := newRunnerContext.getDataFromCommandArgs()
	if err != nil {
		panic(err)
	}

	logger.Info("getting last execution")

	currentSharedInfra := &commonv1alpha1.SharedInfra{}
	err = newRunnerContext.rpcClient.Call("RPCServer.GetSharedInfra", sharedinfra.RPCGetSharedInfraArgs{
		Ref: sharedInfraRef,
	}, currentSharedInfra)
	if err != nil {
		panic(err)
	}

	currentExecutionStatusChann := make(chan commonv1alpha1.ExecutionStatus)

	go newRunnerContext.startTimeout(sharedInfraRef, currentExecutionStatusChann)

	logger.Info("installing terraform provider")
	terraformProvider, err := terraform.NewTerraformProvider(logger, "")
	if err != nil {
		panic(err)
	}

	doneChann := make(chan bool)

	go func() {
		currentExecution := engine.NewEngine(logger, rpcClient, terraformProvider)
		if action == "APPLY" {
			currentExecutionStatusChann <- currentExecution.Apply(*currentSharedInfra, currentExecutionStatusChann)
		} else {
			currentExecution.Destroy(*currentSharedInfra, currentExecutionStatusChann)
		}

		doneChann <- true
	}()

	for {
		select {
		case executionStatus := <-currentExecutionStatusChann:
			logger.Info("New status received calling controller...")
			rpcRunnerFinishedArgs := &sharedinfra.RPCSetExecutionStatusArgs{
				Ref:             sharedInfraRef,
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

func (c runnerContext) getDataFromCommandArgs() (types.NamespacedName, string) {
	commandArgs := os.Args[1:]
	action := commandArgs[0]
	rawSharedInfraRef := commandArgs[1]

	fmt.Println("ARGS", commandArgs)
	fmt.Println("RAW", commandArgs[2])

	sharedInfraRef := types.NamespacedName{}

	s := strings.Split(rawSharedInfraRef, "/")
	if len(s) > 1 {
		sharedInfraRef.Namespace = s[0]
		sharedInfraRef.Name = s[1]
	} else {
		sharedInfraRef.Name = s[0]
		sharedInfraRef.Namespace = "default"
	}

	return sharedInfraRef, action
}

func (c runnerContext) startTimeout(sharedInfraRef types.NamespacedName, currentExecutionStatusChann <-chan commonv1alpha1.ExecutionStatus) {
	time.Sleep(5 * time.Minute)

	rpcRunnerFinishedArgs := &sharedinfra.RPCSetExecutionStatusArgs{
		Ref:             sharedInfraRef,
		ExecutionStatus: <-currentExecutionStatusChann,
	}

	var reply int
	err := c.rpcClient.Call("RPCServer.SetExecutionStatus", rpcRunnerFinishedArgs, &reply)
	if err != nil {
		c.logger.Fatal("Error to call controller", zap.Error(err))
	}

	c.logger.Fatal("Runner timeout")
}
