package secret

import (
	"context"

	v1 "k8s.io/api/core/v1"
)

type Secret struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Data      map[string][]byte `json:"data"`
}

type UseCase interface {
	// List(ctx context.Context, namespace string, chunkPagination pagination.ChunkingPaginationRequest) (pagination.ChunkingPaginationResponse[Repository], error)
	Apply(ctx context.Context, s Secret) (Secret, error)
	Get(ctx context.Context, name string, namespace string) (Secret, error)
	Delete(ctx context.Context, name string, namespace string) error
}

type Repository interface {
	// List(ctx context.Context, namespace string, chunkPagination pagination.ChunkingPaginationRequest) (commonv1alpha1.RepositoryList, error)
	Apply(ctx context.Context, s v1.Secret) (v1.Secret, error)
	Get(ctx context.Context, name string, namespace string) (v1.Secret, error)
	Delete(ctx context.Context, name string, namespace string) error
}
