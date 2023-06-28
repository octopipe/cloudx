package sharedinfra

import (
	"context"
	"fmt"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/engine"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type RPCServer struct {
	client.Client
	logger *zap.Logger
}

func NewRPCServer(client client.Client, logger *zap.Logger) *RPCServer {
	return &RPCServer{Client: client, logger: logger}
}

type RPCGetRunnerDataArgs struct {
	Ref         types.NamespacedName
	ExecutionId string
}

type RPCGetRunnerDataReply struct {
	SharedInfra commonv1alpha1.SharedInfra
}

func (s *RPCServer) GetRunnerData(args *RPCGetRunnerDataArgs, reply *RPCGetRunnerDataReply) error {
	s.logger.Info("Received rpc call", zap.String("sharedinfra", args.Ref.String()))
	currentSharedInfra := commonv1alpha1.SharedInfra{}
	err := s.Get(context.Background(), args.Ref, &currentSharedInfra)
	if err != nil {
		return err
	}

	reply.SharedInfra = currentSharedInfra
	return nil
}

type RPCSetExecutionStatusArgs struct {
	Ref             types.NamespacedName
	ExecutionStatus commonv1alpha1.ExecutionStatus
}

func (s *RPCServer) SetExecutionStatus(args *RPCSetExecutionStatusArgs, reply *int) error {
	s.logger.Info("Received rpc call", zap.String("sharedinfra", args.Ref.String()))
	currentExecution := &commonv1alpha1.Execution{}
	err := s.Get(context.Background(), args.Ref, currentExecution)
	if err != nil {
		return err
	}

	s.logger.Info("rpc execution", zap.String("status", args.ExecutionStatus.Status))

	currentExecution.Status = args.ExecutionStatus

	return updateExecutionStatus(s.Client, currentExecution)
}

type RPCSetRunnerTimeoutArgs struct {
	Plugins    []commonv1alpha1.PluginExecutionStatus
	Ref        types.NamespacedName
	FinishedAt string
}

func (s *RPCServer) SetRunnerTimeout(args *RPCSetRunnerTimeoutArgs, reply *int) error {
	currentExecution := &commonv1alpha1.Execution{}
	err := s.Get(context.Background(), args.Ref, currentExecution)
	if err != nil {
		return err
	}

	runnerList := &v1.PodList{}
	selector, _ := labels.Parse(fmt.Sprintf("commons.cloudx.io/execution=%s", args.Ref.String()))
	err = s.List(context.Background(), runnerList, &client.ListOptions{LabelSelector: selector})
	if err != nil {
		return err
	}

	for _, r := range runnerList.Items {
		err = s.Delete(context.Background(), &r)
		if err != nil {
			return err
		}
	}

	currentExecutionStatus := commonv1alpha1.ExecutionStatus{
		StartedAt:  currentExecution.Status.StartedAt,
		Status:     engine.ExecutionTimeout,
		FinishedAt: args.FinishedAt,
		Error:      "Runner time exceeded",
		Plugins:    args.Plugins,
	}

	currentExecution.Status = currentExecutionStatus

	return updateExecutionStatus(s.Client, currentExecution)
}

type RPCGetLastExecutionArgs struct {
	Ref types.NamespacedName
}

func (s *RPCServer) GetLastExecution(args *RPCGetLastExecutionArgs, reply *commonv1alpha1.Execution) error {

	currentExecution := &commonv1alpha1.Execution{}
	err := s.Get(context.Background(), args.Ref, currentExecution)
	if err != nil {
		return err
	}

	s.logger.Info("current execution", zap.String("name", args.Ref.String()))

	currentSharedInfra := &commonv1alpha1.SharedInfra{}
	sharedInfraRef := types.NamespacedName{
		Name:      currentExecution.Spec.SharedInfra.Name,
		Namespace: currentExecution.Spec.SharedInfra.Namespace,
	}
	err = s.Get(context.Background(), sharedInfraRef, currentSharedInfra)
	if err != nil {
		return err
	}

	for _, e := range currentSharedInfra.Status.Executions {
		executionApi := &commonv1alpha1.Execution{}
		err = s.Get(context.Background(), types.NamespacedName{Name: e.Name, Namespace: e.Namespace}, executionApi)
		if err != nil {
			return err
		}

		if executionApi.Status.Status != engine.ExecutionRunningStatus {
			*reply = *executionApi
			return nil
		}
	}

	return nil
}
