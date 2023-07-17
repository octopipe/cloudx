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
func (u useCase) Create(ctx context.Context, taskOutput TaskOutput) (TaskOutput, error) {
	newTaskOutput := commonv1alpha1.TaskOutput{
		Spec: taskOutput.TaskOutputSpec,
	}

	newTaskOutput.SetName(taskOutput.Name)
	newTaskOutput.SetNamespace(taskOutput.Namespace)

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

	taskOutputs := []TaskOutput{}
	for _, i := range l.Items {
		taskOutputs = append(taskOutputs, TaskOutput{
			Name:           i.GetName(),
			Namespace:      i.GetNamespace(),
			TaskOutputSpec: i.Spec,
		})
	}

	return pagination.ChunkingPaginationResponse[TaskOutput]{
		Items: taskOutputs,
		Chunk: l.Continue,
	}, nil
}

// Update implements UseCase.
func (u useCase) Update(ctx context.Context, taskOutput TaskOutput) (TaskOutput, error) {
	newTaskOutput := commonv1alpha1.TaskOutput{
		Spec: taskOutput.TaskOutputSpec,
	}

	newTaskOutput.SetName(taskOutput.Name)
	newTaskOutput.SetNamespace(taskOutput.Namespace)

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
