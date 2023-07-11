package taskoutput

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
func (u useCase) Create(ctx context.Context, connectionInterface TaskOutput) (TaskOutput, error) {
	newTaskOutput := commonv1alpha1.TaskOutput{
		Spec: connectionInterface.TaskOutputSpec,
	}

	newTaskOutput.SetName(connectionInterface.Name)
	newTaskOutput.SetNamespace(connectionInterface.Namespace)

	s, err := u.repository.Apply(ctx, newTaskOutput)
	if err != nil {
		return TaskOutput{}, err
	}

	return TaskOutput{
		Name:           s.GetName(),
		Namespace:      s.GetNamespace(),
		TaskOutputSpec: s.Spec,
	}, nil
}

// Delete implements UseCase.
func (u useCase) Delete(ctx context.Context, name string, namespace string) error {
	return u.repository.Delete(ctx, name, namespace)
}

// Get implements UseCase.
func (u useCase) Get(ctx context.Context, name string, namespace string) (TaskOutput, error) {
	s, err := u.repository.Get(ctx, name, namespace)
	if err != nil {
		return TaskOutput{}, err
	}

	return TaskOutput{
		Name:           s.GetName(),
		Namespace:      s.GetNamespace(),
		TaskOutputSpec: s.Spec,
	}, nil
}

// List implements UseCase.
func (u useCase) List(ctx context.Context, namespace string, chunkPagination pagination.ChunkingPaginationRequest) (pagination.ChunkingPaginationResponse[TaskOutput], error) {
	l, err := u.repository.List(ctx, namespace, chunkPagination)
	if err != nil {
		return pagination.ChunkingPaginationResponse[TaskOutput]{}, err
	}

	connectionInterfaces := []TaskOutput{}
	for _, i := range l.Items {
		connectionInterfaces = append(connectionInterfaces, TaskOutput{
			Name:           i.GetName(),
			Namespace:      i.GetNamespace(),
			TaskOutputSpec: i.Spec,
		})
	}

	return pagination.ChunkingPaginationResponse[TaskOutput]{
		Items: connectionInterfaces,
		Chunk: l.Continue,
	}, nil
}

// Update implements UseCase.
func (u useCase) Update(ctx context.Context, connectionInterface TaskOutput) (TaskOutput, error) {
	newTaskOutput := commonv1alpha1.TaskOutput{
		Spec: connectionInterface.TaskOutputSpec,
	}

	newTaskOutput.SetName(connectionInterface.Name)
	newTaskOutput.SetNamespace(connectionInterface.Namespace)

	s, err := u.repository.Apply(ctx, newTaskOutput)
	if err != nil {
		return TaskOutput{}, err
	}

	return TaskOutput{
		Name:           s.GetName(),
		Namespace:      s.GetNamespace(),
		TaskOutputSpec: s.Spec,
	}, nil
}
