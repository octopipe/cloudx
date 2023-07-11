package taskoutput

import (
	"context"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
)

type TaskOutputRPCHandler struct {
	connectionInterfaceRepository Repository
}

func NewTaskOutputRPCHandler(connectionInterfaceRepository Repository) *TaskOutputRPCHandler {
	return &TaskOutputRPCHandler{connectionInterfaceRepository: connectionInterfaceRepository}
}

type RPCGetTaskOutputArgs struct {
	Ref types.NamespacedName
}

func (h *TaskOutputRPCHandler) GetTaskOutput(args *RPCGetTaskOutputArgs, reply *commonv1alpha1.TaskOutput) error {
	currentTaskOutput, err := h.connectionInterfaceRepository.Get(context.Background(), args.Ref.Name, args.Ref.Namespace)
	*reply = currentTaskOutput
	return err
}
