package infra

import (
	"context"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/pagination"
)

const (
	ApplyAction   = "APPLY"
	DestroyAction = "DESTROY"
)

type InfraStatus struct {
	Tasks      []commonv1alpha1.TaskExecutionStatus `json:"tasks"`
	StartedAt  string                               `json:"startedAt"`
	FinishedAt string                               `json:"finishedAt"`
	Status     string                               `json:"status"`
	Error      string                               `json:"error"`
}

type Infra struct {
	Name      string      `json:"name"`
	Namespace string      `json:"namespace"`
	Status    InfraStatus `json:"status"`
	commonv1alpha1.InfraSpec
}

type UseCase interface {
	List(ctx context.Context, namespace string, chunkPagination pagination.ChunkingPaginationRequest) (pagination.ChunkingPaginationResponse[Infra], error)
	Create(ctx context.Context, infra Infra) (Infra, error)
	Update(ctx context.Context, infra Infra) (Infra, error)
	Get(ctx context.Context, name string, namespace string) (Infra, error)
	Reconcile(ctx context.Context, name string, namespace string) error
	Delete(ctx context.Context, name string, namespace string) error
}

type Repository interface {
	List(ctx context.Context, namespace string, chunkPagination pagination.ChunkingPaginationRequest) (commonv1alpha1.InfraList, error)
	Apply(ctx context.Context, s commonv1alpha1.Infra) (commonv1alpha1.Infra, error)
	Get(ctx context.Context, name string, namespace string) (commonv1alpha1.Infra, error)
	Reconcile(ctx context.Context, name string, namespace string) error
	Delete(ctx context.Context, name string, namespace string) error
}
