package execution

import (
	"context"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/pagination"
)

type Execution struct {
	Name      string                         `json:"name"`
	Namespace string                         `json:"namespace"`
	Status    commonv1alpha1.ExecutionStatus `json:"status"`
	commonv1alpha1.ExecutionSpec
}

type UseCase interface {
	List(ctx context.Context, namespace string, chunkPagination pagination.ChunkingPaginationRequest) (pagination.ChunkingPaginationResponse[Execution], error)
	Create(ctx context.Context, execution Execution) (Execution, error)
	Update(ctx context.Context, execution Execution) (Execution, error)
	Get(ctx context.Context, name string, namespace string) (Execution, error)
	Delete(ctx context.Context, name string, namespace string) error
}

type Repository interface {
	List(ctx context.Context, namespace string, chunkPagination pagination.ChunkingPaginationRequest) (commonv1alpha1.ExecutionList, error)
	Apply(ctx context.Context, s commonv1alpha1.Execution) (commonv1alpha1.Execution, error)
	Get(ctx context.Context, name string, namespace string) (commonv1alpha1.Execution, error)
	Delete(ctx context.Context, name string, namespace string) error
}
