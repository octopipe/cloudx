package connectioninterface

import (
	"context"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
)

type rpcHandler struct {
	connectionInterfaceRepository Repository
}

func NewRPCHandler(connectionInterfaceRepository Repository) rpcHandler {
	return rpcHandler{connectionInterfaceRepository: connectionInterfaceRepository}
}

type RPCGetConnectionInterfaceArgs struct {
	Ref types.NamespacedName
}

func (h *rpcHandler) GetConnectionInterface(args *RPCGetConnectionInterfaceArgs, reply *commonv1alpha1.ConnectionInterface) error {
	currentConnectionInterface, err := h.connectionInterfaceRepository.Get(context.Background(), args.Ref.Name, args.Ref.Namespace)
	reply = &currentConnectionInterface
	return err
}
