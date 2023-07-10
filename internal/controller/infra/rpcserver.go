package infra

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
	Infra commonv1alpha1.Infra
}

func (s *RPCServer) GetRunnerData(args *RPCGetRunnerDataArgs, reply *RPCGetRunnerDataReply) error {
	s.logger.Info("Received rpc call", zap.String("infra", args.Ref.String()))
	currentInfra := commonv1alpha1.Infra{}
	err := s.Get(context.Background(), args.Ref, &currentInfra)
	if err != nil {
		return err
	}

	reply.Infra = currentInfra
	return nil
}

type RPCSetExecutionStatusArgs struct {
	Ref             types.NamespacedName
	ExecutionStatus commonv1alpha1.ExecutionStatus
}

func (s *RPCServer) SetExecutionStatus(args *RPCSetExecutionStatusArgs, reply *int) error {
	s.logger.Info("received call", zap.String("method", "RPCServer.SetExecutionStatus"), zap.String("infra", args.Ref.String()))
	infra := &commonv1alpha1.Infra{}
	err := s.Get(context.Background(), args.Ref, infra)
	if err != nil {
		s.logger.Error("Failed to get current execution", zap.String("method", "RPCServer.SetExecutionStatus"), zap.String("infra", args.Ref.String()), zap.Error(err))
		return err
	}

	infra.Status.LastExecution = args.ExecutionStatus
	s.logger.Info("updating current execution status", zap.String("method", "RPCServer.SetExecutionStatus"), zap.String("status", args.ExecutionStatus.Status))
	err = utils.UpdateInfraStatus(s.Client, *infra)
	if err != nil {
		s.logger.Error("Failed to update current execution status", zap.String("method", "RPCServer.SetExecutionStatus"), zap.String("infra", args.Ref.String()), zap.Error(err))
		return err
	}

	return nil
}

type RPCSetRunnerTimeoutArgs struct {
	Tasks []commonv1alpha1.TaskExecutionStatus
	Ref   types.NamespacedName
}

func (s *RPCServer) SetRunnerTimeout(args *RPCSetRunnerTimeoutArgs, reply *int) error {
	infra := &commonv1alpha1.Infra{}
	err := s.Get(context.Background(), args.Ref, infra)
	if err != nil {
		return err
	}

	infra.Status.LastExecution.Status = engine.ExecutionTimeout
	infra.Status.LastExecution.FinishedAt = time.Now().Format(time.RFC3339)
	infra.Status.LastExecution.Error = "Runner time exceeded"
	infra.Status.LastExecution.Tasks = args.Tasks

	return utils.UpdateInfraStatus(s.Client, *infra)
}

type RPCGetLastExecutionArgs struct {
	Ref types.NamespacedName
}

func (s *RPCServer) GetLastExecution(args *RPCGetLastExecutionArgs, reply *commonv1alpha1.ExecutionStatus) error {
	s.logger.Info("get last execution rpc all", zap.String("name", args.Ref.String()))

	infra := &commonv1alpha1.Infra{}
	err := s.Get(context.Background(), args.Ref, infra)
	if err != nil {
		s.logger.Error("failed to get current execution", zap.Error(err))
		return err
	}

	*reply = infra.Status.LastExecution
	return nil
}

type RPCGetInfraArgs struct {
	Ref types.NamespacedName
}

func (s *RPCServer) GetInfra(args *RPCGetInfraArgs, reply *commonv1alpha1.Infra) error {
	s.logger.Info("get shared infra rpc call", zap.String("name", args.Ref.String()))

	infra := &commonv1alpha1.Infra{}
	err := s.Get(context.Background(), args.Ref, infra)
	if err != nil {
		s.logger.Error("failed to get shared infra", zap.Error(err))
		return err
	}

	*reply = *infra
	return nil
}
