package execution

import (
	"context"
	"time"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/pagination"
	"github.com/octopipe/cloudx/internal/sharedinfra"
)

type useCase struct {
	repository         Repository
	sharedInfraUseCase sharedinfra.UseCase
}

func NewUseCase(repository Repository, sharedInfraUseCase sharedinfra.UseCase) UseCase {
	return useCase{repository: repository, sharedInfraUseCase: sharedInfraUseCase}
}

// Create implements UseCase.
func (u useCase) Create(ctx context.Context, sharedInfraName string, namespace string, execution Execution) (Execution, error) {
	sharedInfra, err := u.sharedInfraUseCase.Get(ctx, sharedInfraName, namespace)
	if err != nil {
		return Execution{}, err
	}

	newExecution := commonv1alpha1.Execution{
		Spec: commonv1alpha1.ExecutionSpec{
			Author:    execution.Author,
			Action:    execution.Action,
			StartedAt: time.Now().Format(time.RFC3339),
			SharedInfra: commonv1alpha1.Ref{
				Name:      sharedInfra.Name,
				Namespace: sharedInfra.Namespace,
			},
		},
	}

	newExecution.SetName(execution.Name)
	newExecution.SetNamespace(execution.Namespace)

	s, err := u.repository.Apply(ctx, newExecution)
	if err != nil {
		return Execution{}, err
	}

	return Execution{
		Name:      s.GetName(),
		Namespace: s.GetNamespace(),
		Author:    s.Spec.Author,
		Action:    s.Spec.Action,
		Status:    s.Status,
	}, nil
}

// Delete implements UseCase.
func (u useCase) Delete(ctx context.Context, sharedInfraName string, name string, namespace string) error {
	return u.repository.Delete(ctx, name, namespace)
}

// Get implements UseCase.
func (u useCase) Get(ctx context.Context, sharedInfraName string, name string, namespace string) (Execution, error) {
	s, err := u.repository.Get(ctx, name, namespace)
	if err != nil {
		return Execution{}, err
	}

	return Execution{
		Name:      s.GetName(),
		Namespace: s.GetNamespace(),
		Author:    s.Spec.Author,
		Action:    s.Spec.Action,
		Status:    s.Status,
	}, nil
}

// List implements UseCase.
func (u useCase) List(ctx context.Context, sharedInfraName string, namespace string, chunkPagination pagination.ChunkingPaginationRequest) (pagination.ChunkingPaginationResponse[Execution], error) {
	l, err := u.repository.List(ctx, sharedInfraName, namespace, chunkPagination)
	if err != nil {
		return pagination.ChunkingPaginationResponse[Execution]{}, err
	}

	executions := []Execution{}
	for _, i := range l.Items {
		executions = append(executions, Execution{
			Name:      i.GetName(),
			Namespace: i.GetNamespace(),
			Author:    i.Spec.Author,
			Action:    i.Spec.Action,
			Status:    i.Status,
			StartedAt: i.Spec.StartedAt,
		})
	}

	return pagination.ChunkingPaginationResponse[Execution]{
		Items: executions,
		Chunk: l.Continue,
	}, nil
}

// Update implements UseCase.
func (u useCase) Update(ctx context.Context, sharedInfraName string, namespace string, execution Execution) (Execution, error) {
	sharedInfra, err := u.sharedInfraUseCase.Get(ctx, sharedInfraName, namespace)
	if err != nil {
		return Execution{}, err
	}

	newExecution := commonv1alpha1.Execution{
		Spec: commonv1alpha1.ExecutionSpec{
			Author:    execution.Author,
			Action:    execution.Action,
			StartedAt: time.Now().Format(time.RFC3339),
			SharedInfra: commonv1alpha1.Ref{
				Name:      sharedInfra.Name,
				Namespace: sharedInfra.Namespace,
			},
		},
	}

	newExecution.SetName(execution.Name)
	newExecution.SetNamespace(execution.Namespace)

	s, err := u.repository.Apply(ctx, newExecution)
	if err != nil {
		return Execution{}, err
	}

	return Execution{
		Name:      s.GetName(),
		Namespace: s.GetNamespace(),
		Author:    s.Spec.Author,
		Action:    s.Spec.Action,
		Status:    s.Status,
	}, nil
}
