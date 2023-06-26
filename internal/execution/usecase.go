package execution

import (
	"context"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/pagination"
)

type useCase struct {
	repository Repository
}

func NewUseCase(repository Repository) UseCase {
	return useCase{repository: repository}
}

// Create implements UseCase.
func (u useCase) Create(ctx context.Context, execution Execution) (Execution, error) {
	newExecution := commonv1alpha1.Execution{
		Spec: execution.ExecutionSpec,
	}

	newExecution.SetName(execution.Name)
	newExecution.SetNamespace(execution.Namespace)

	s, err := u.repository.Apply(ctx, newExecution)
	if err != nil {
		return Execution{}, err
	}

	return Execution{
		Name:          s.GetName(),
		Namespace:     s.GetNamespace(),
		ExecutionSpec: s.Spec,
		Status:        s.Status,
	}, nil
}

// Delete implements UseCase.
func (u useCase) Delete(ctx context.Context, name string, namespace string) error {
	return u.repository.Delete(ctx, name, namespace)
}

// Get implements UseCase.
func (u useCase) Get(ctx context.Context, name string, namespace string) (Execution, error) {
	s, err := u.repository.Get(ctx, name, namespace)
	if err != nil {
		return Execution{}, err
	}

	return Execution{
		Name:          s.GetName(),
		Namespace:     s.GetNamespace(),
		ExecutionSpec: s.Spec,
		Status:        s.Status,
	}, nil
}

// List implements UseCase.
func (u useCase) List(ctx context.Context, namespace string, chunkPagination pagination.ChunkingPaginationRequest) (pagination.ChunkingPaginationResponse[Execution], error) {
	l, err := u.repository.List(ctx, namespace, chunkPagination)
	if err != nil {
		return pagination.ChunkingPaginationResponse[Execution]{}, err
	}

	executions := []Execution{}
	for _, i := range l.Items {
		executions = append(executions, Execution{
			Name:          i.GetName(),
			Namespace:     i.GetNamespace(),
			ExecutionSpec: i.Spec,
			Status:        i.Status,
		})
	}

	return pagination.ChunkingPaginationResponse[Execution]{
		Items: executions,
		Chunk: l.Continue,
	}, nil
}

// Update implements UseCase.
func (u useCase) Update(ctx context.Context, execution Execution) (Execution, error) {
	newExecution := commonv1alpha1.Execution{
		Spec: execution.ExecutionSpec,
	}

	newExecution.SetName(execution.Name)
	newExecution.SetNamespace(execution.Namespace)

	s, err := u.repository.Apply(ctx, newExecution)
	if err != nil {
		return Execution{}, err
	}

	return Execution{
		Name:          s.GetName(),
		Namespace:     s.GetNamespace(),
		ExecutionSpec: s.Spec,
		Status:        s.Status,
	}, nil
}
