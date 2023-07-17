package taskoutput

import (
	"context"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/pagination"
)

type TaskOutput struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	commonv1alpha1.TaskOutputSpec
}

type UseCase interface {
	List(ctx context.Context, namespace string, chunkPagination pagination.ChunkingPaginationRequest) (pagination.ChunkingPaginationResponse[TaskOutput], error)
	Create(ctx context.Context, taskOutput TaskOutput) (TaskOutput, error)
	Update(ctx context.Context, taskOutput TaskOutput) (TaskOutput, error)
	Get(ctx context.Context, name string, namespace string) (TaskOutput, error)
	Delete(ctx context.Context, name string, namespace string) error
}

type Repository interface {
	List(ctx context.Context, namespace string, chunkPagination pagination.ChunkingPaginationRequest) (commonv1alpha1.TaskOutputList, error)
	Apply(ctx context.Context, s commonv1alpha1.TaskOutput) (commonv1alpha1.TaskOutput, error)
	Get(ctx context.Context, name string, namespace string) (commonv1alpha1.TaskOutput, error)
	Delete(ctx context.Context, name string, namespace string) error
}
