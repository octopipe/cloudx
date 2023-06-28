package connectioninterface

import (
	"context"
	"fmt"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
)

type ConnectionInterfaceRPCHandler struct {
	connectionInterfaceRepository Repository
}

func NewConnectionInterfaceRPCHandler(connectionInterfaceRepository Repository) *ConnectionInterfaceRPCHandler {
	return &ConnectionInterfaceRPCHandler{connectionInterfaceRepository: connectionInterfaceRepository}
}

type RPCGetConnectionInterfaceArgs struct {
	Ref types.NamespacedName
}

func (h *ConnectionInterfaceRPCHandler) GetConnectionInterface(args *RPCGetConnectionInterfaceArgs, reply *commonv1alpha1.ConnectionInterface) error {
	fmt.Println(args.Ref)
	currentConnectionInterface, err := h.connectionInterfaceRepository.Get(context.Background(), args.Ref.Name, args.Ref.Namespace)
	*reply = currentConnectionInterface
	return err
}
