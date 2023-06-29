package providerconfig

import (
	"context"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"k8s.io/apimachinery/pkg/types"
)

type ProviderConfigRPCHandler struct {
	providerConfigRepository Repository
}

func NewProviderConfigRPCHandler(providerConfigRepository Repository) *ProviderConfigRPCHandler {
	return &ProviderConfigRPCHandler{providerConfigRepository: providerConfigRepository}
}

type RPCGetProviderConfigArgs struct {
	Ref types.NamespacedName
}

func (h *ProviderConfigRPCHandler) GetProviderConfig(args *RPCGetProviderConfigArgs, reply *commonv1alpha1.ProviderConfig) error {
	currentProviderConfig, err := h.providerConfigRepository.Get(context.Background(), args.Ref.Name, args.Ref.Namespace)
	*reply = currentProviderConfig
	return err
}
