package providerconfig

import (
	"context"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/pagination"
)

type ProviderConfig struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	commonv1alpha1.ProviderConfigSpec
}

type UseCase interface {
	List(ctx context.Context, namespace string, chunkPagination pagination.ChunkingPaginationRequest) (pagination.ChunkingPaginationResponse[ProviderConfig], error)
	Create(ctx context.Context, providerConfig ProviderConfig) (ProviderConfig, error)
	Update(ctx context.Context, providerConfig ProviderConfig) (ProviderConfig, error)
	Get(ctx context.Context, name string, namespace string) (ProviderConfig, error)
	Delete(ctx context.Context, name string, namespace string) error
}

type Repository interface {
	List(ctx context.Context, namespace string, chunkPagination pagination.ChunkingPaginationRequest) (commonv1alpha1.ProviderConfigList, error)
	Apply(ctx context.Context, s commonv1alpha1.ProviderConfig) (commonv1alpha1.ProviderConfig, error)
	Get(ctx context.Context, name string, namespace string) (commonv1alpha1.ProviderConfig, error)
	Delete(ctx context.Context, name string, namespace string) error
}
