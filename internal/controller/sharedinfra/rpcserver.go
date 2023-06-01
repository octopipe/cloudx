package sharedinfra

import (
	"context"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
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

type RPCSetRunnerFinishedArgs struct {
	Ref       types.NamespacedName
	Execution commonv1alpha1.SharedInfraExecutionStatus
}

func (s *RPCServer) SetRunnerFinished(args *RPCSetRunnerFinishedArgs, reply *int) error {
	s.logger.Info("Received rpc call", zap.String("sharedinfra", args.Ref.String()))
	currentSharedInfra := &commonv1alpha1.SharedInfra{}
	err := s.Get(context.Background(), args.Ref, currentSharedInfra)
	if err != nil {
		return err
	}

	s.logger.Info("rpc execution", zap.String("status", args.Execution.Status))

	currentExecutions := currentSharedInfra.Status.Executions
	currentExecutions = append(currentExecutions, args.Execution)
	currentSharedInfra.Status.Executions = currentExecutions

	return s.Status().Update(context.TODO(), currentSharedInfra)
}
