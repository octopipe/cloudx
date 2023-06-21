package connectioninterface

import (
	"context"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/pagination"
)

type ConnectionInterface struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	commonv1alpha1.ConnectionInterfaceSpec
}

type UseCase interface {
	List(ctx context.Context, namespace string, chunkPagination pagination.ChunkingPaginationRequest) (pagination.ChunkingPaginationResponse[ConnectionInterface], error)
	Create(ctx context.Context, connectionInterface ConnectionInterface) (ConnectionInterface, error)
	Update(ctx context.Context, connectionInterface ConnectionInterface) (ConnectionInterface, error)
	Get(ctx context.Context, name string, namespace string) (ConnectionInterface, error)
	Delete(ctx context.Context, name string, namespace string) error
}

type Repository interface {
	List(ctx context.Context, namespace string, chunkPagination pagination.ChunkingPaginationRequest) (commonv1alpha1.ConnectionInterfaceList, error)
	Apply(ctx context.Context, s commonv1alpha1.ConnectionInterface) (commonv1alpha1.ConnectionInterface, error)
	Get(ctx context.Context, name string, namespace string) (commonv1alpha1.ConnectionInterface, error)
	Delete(ctx context.Context, name string, namespace string) error
}
