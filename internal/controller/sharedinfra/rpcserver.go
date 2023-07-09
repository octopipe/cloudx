package sharedinfra

import (
	"context"
	"time"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/controller/utils"
	"github.com/octopipe/cloudx/internal/engine"
	"go.uber.org/zap"
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
	s.logger.Info("received call", zap.String("method", "RPCServer.SetExecutionStatus"), zap.String("sharedinfra", args.Ref.String()))
	sharedInfra := &commonv1alpha1.SharedInfra{}
	err := s.Get(context.Background(), args.Ref, sharedInfra)
	if err != nil {
		s.logger.Error("Failed to get current execution", zap.String("method", "RPCServer.SetExecutionStatus"), zap.String("sharedinfra", args.Ref.String()), zap.Error(err))
		return err
	}

	sharedInfra.Status.LastExecution = args.ExecutionStatus
	s.logger.Info("updating current execution status", zap.String("method", "RPCServer.SetExecutionStatus"), zap.String("status", args.ExecutionStatus.Status))
	err = utils.UpdateSharedInfraStatus(s.Client, *sharedInfra)
	if err != nil {
		s.logger.Error("Failed to update current execution status", zap.String("method", "RPCServer.SetExecutionStatus"), zap.String("sharedinfra", args.Ref.String()), zap.Error(err))
		return err
	}

	return nil
}

type RPCSetRunnerTimeoutArgs struct {
	Plugins []commonv1alpha1.PluginExecutionStatus
	Ref     types.NamespacedName
}

func (s *RPCServer) SetRunnerTimeout(args *RPCSetRunnerTimeoutArgs, reply *int) error {
	sharedInfra := &commonv1alpha1.SharedInfra{}
	err := s.Get(context.Background(), args.Ref, sharedInfra)
	if err != nil {
		return err
	}

	sharedInfra.Status.LastExecution.Status = engine.ExecutionTimeout
	sharedInfra.Status.LastExecution.FinishedAt = time.Now().Format(time.RFC3339)
	sharedInfra.Status.LastExecution.Error = "Runner time exceeded"
	sharedInfra.Status.LastExecution.Plugins = args.Plugins

	return utils.UpdateSharedInfraStatus(s.Client, *sharedInfra)
}

type RPCGetLastExecutionArgs struct {
	Ref types.NamespacedName
}

func (s *RPCServer) GetLastExecution(args *RPCGetLastExecutionArgs, reply *commonv1alpha1.ExecutionStatus) error {
	s.logger.Info("get last execution rpc all", zap.String("name", args.Ref.String()))

	sharedInfra := &commonv1alpha1.SharedInfra{}
	err := s.Get(context.Background(), args.Ref, sharedInfra)
	if err != nil {
		s.logger.Error("failed to get current execution", zap.Error(err))
		return err
	}

	*reply = sharedInfra.Status.LastExecution
	return nil
}

type RPCGetSharedInfraArgs struct {
	Ref types.NamespacedName
}

func (s *RPCServer) GetSharedInfra(args *RPCGetSharedInfraArgs, reply *commonv1alpha1.SharedInfra) error {
	s.logger.Info("get shared infra rpc call", zap.String("name", args.Ref.String()))

	sharedInfra := &commonv1alpha1.SharedInfra{}
	err := s.Get(context.Background(), args.Ref, sharedInfra)
	if err != nil {
		s.logger.Error("failed to get shared infra", zap.Error(err))
		return err
	}

	*reply = *sharedInfra
	return nil
}
