package main

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/backend"
	"github.com/octopipe/cloudx/internal/backend/terraform"
	"github.com/octopipe/cloudx/internal/controller/infra"
	"github.com/octopipe/cloudx/internal/pipeline"
	"github.com/octopipe/cloudx/internal/rpcclient"
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

	infraRef, action := newRunnerContext.getDataFromCommandArgs()
	if err != nil {
		panic(err)
	}

	logger.Info("getting last execution")
	currentInfra := &commonv1alpha1.Infra{}
	err = newRunnerContext.rpcClient.Call("RPCServer.GetInfra", infra.RPCGetInfraArgs{
		Ref: infraRef,
	}, currentInfra)
	if err != nil {
		panic(err)
	}

	statusChan := make(chan commonv1alpha1.ExecutionStatus)
	terraformBackend, err := terraform.NewTerraformBackend(logger)
	if err != nil {
		panic(err)
	}

	newBackend := backend.NewBackend(terraformBackend)
	newPipeline := pipeline.NewPipeline(logger, rpcClient, newBackend)

	go func() {
		logger.Info("start pipeline execution")
		newPipeline.Start(action, *currentInfra, statusChan)
	}()

	for executionStatus := range statusChan {
		err = newRunnerContext.setExecutionStatus(infraRef, executionStatus)
		if err != nil {
			logger.Fatal("Failed to call rpc execution status", zap.Error(err))
		}

		if executionStatus.Status == pipeline.InfraSuccessStatus || executionStatus.Status == pipeline.InfraErrorStatus || executionStatus.Status == pipeline.InfraTimeoutStatus {
			logger.Info("Finish engine execution")
			return
		}
	}
}

func (c runnerContext) setExecutionStatus(infraRef types.NamespacedName, executionStatus commonv1alpha1.ExecutionStatus) error {
	c.logger.Info("New status received calling controller...")
	rpcRunnerFinishedArgs := &infra.RPCSetExecutionStatusArgs{
		Ref:             infraRef,
		ExecutionStatus: executionStatus,
	}

	var reply int
	err := c.rpcClient.Call("RPCServer.SetExecutionStatus", rpcRunnerFinishedArgs, &reply)
	if err != nil {
		return err
	}

	return nil
}

func (c runnerContext) getDataFromCommandArgs() (types.NamespacedName, string) {
	commandArgs := os.Args[1:]
	action := commandArgs[0]
	rawInfraRef := commandArgs[1]

	infraRef := types.NamespacedName{}

	s := strings.Split(rawInfraRef, "/")
	if len(s) > 1 {
		infraRef.Namespace = s[0]
		infraRef.Name = s[1]
	} else {
		infraRef.Name = s[0]
		infraRef.Namespace = "default"
	}

	return infraRef, action
}
