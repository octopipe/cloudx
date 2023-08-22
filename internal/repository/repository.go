package repository

import (
	"context"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/pagination"
)

const (
	ApplyAction   = "APPLY"
	DestroyAction = "DESTROY"
)

type RepositoryAuth struct {
	Type      string `json:"type,omitempty"`
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
	PublicKey string `json:"publicKey,omitempty"`
}

type Repository struct {
	Name      string                          `json:"name"`
	Namespace string                          `json:"namespace"`
	Auth      RepositoryAuth                  `json:"auth"`
	Status    commonv1alpha1.RepositoryStatus `json:"status"`
	commonv1alpha1.RepositorySpec
}

type UseCase interface {
	List(ctx context.Context, namespace string, chunkPagination pagination.ChunkingPaginationRequest) (pagination.ChunkingPaginationResponse[Repository], error)
	Create(ctx context.Context, repository Repository) (Repository, error)
	Update(ctx context.Context, repository Repository) (Repository, error)
	Get(ctx context.Context, name string, namespace string) (Repository, error)
	Sync(ctx context.Context, name string, namespace string) ([]string, error)
	Delete(ctx context.Context, name string, namespace string) error
}

type RepositoryType interface {
	List(ctx context.Context, namespace string, chunkPagination pagination.ChunkingPaginationRequest) (commonv1alpha1.RepositoryList, error)
	Apply(ctx context.Context, s commonv1alpha1.Repository) (commonv1alpha1.Repository, error)
	Get(ctx context.Context, name string, namespace string) (commonv1alpha1.Repository, error)
	Sync(ctx context.Context, name string, namespace string) error
	Delete(ctx context.Context, name string, namespace string) error
}
