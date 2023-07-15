package infra

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
func (u useCase) Create(ctx context.Context, infra Infra) (Infra, error) {
	newInfra := commonv1alpha1.Infra{
		Spec: infra.InfraSpec,
	}

	newInfra.SetName(infra.Name)
	newInfra.SetNamespace(infra.Namespace)

	s, err := u.repository.Apply(ctx, newInfra)
	if err != nil {
		return Infra{}, err
	}

	return Infra{
		Name:      s.GetName(),
		Namespace: s.GetNamespace(),
		InfraSpec: s.Spec,
		Status: InfraStatus{
			StartedAt:  s.Status.LastExecution.StartedAt,
			FinishedAt: s.Status.LastExecution.FinishedAt,
			Status:     s.Status.LastExecution.Status,
			Error:      s.Status.LastExecution.Error,
			Tasks:      maskTasksSensitiveData(s.Status.LastExecution.Tasks),
		},
	}, nil
}

// Delete implements UseCase.
func (u useCase) Delete(ctx context.Context, name string, namespace string) error {
	return u.repository.Delete(ctx, name, namespace)
}

// Get implements UseCase.
func (u useCase) Get(ctx context.Context, name string, namespace string) (Infra, error) {
	s, err := u.repository.Get(ctx, name, namespace)
	if err != nil {
		return Infra{}, err
	}

	return Infra{
		Name:      s.GetName(),
		Namespace: s.GetNamespace(),
		InfraSpec: s.Spec,
		Status: InfraStatus{
			StartedAt:  s.Status.LastExecution.StartedAt,
			FinishedAt: s.Status.LastExecution.FinishedAt,
			Status:     s.Status.LastExecution.Status,
			Error:      s.Status.LastExecution.Error,
			Tasks:      maskTasksSensitiveData(s.Status.LastExecution.Tasks),
		},
	}, nil
}

func (u useCase) Reconcile(ctx context.Context, name string, namespace string) error {
	return u.repository.Reconcile(ctx, name, namespace)
}

// List implements UseCase.
func (u useCase) List(ctx context.Context, namespace string, chunkPagination pagination.ChunkingPaginationRequest) (pagination.ChunkingPaginationResponse[Infra], error) {
	l, err := u.repository.List(ctx, namespace, chunkPagination)
	if err != nil {
		return pagination.ChunkingPaginationResponse[Infra]{}, err
	}

	infras := []Infra{}
	for _, i := range l.Items {
		infras = append(infras, Infra{
			Name:      i.GetName(),
			Namespace: i.GetNamespace(),
			InfraSpec: i.Spec,
			Status: InfraStatus{
				StartedAt:  i.Status.LastExecution.StartedAt,
				FinishedAt: i.Status.LastExecution.FinishedAt,
				Status:     i.Status.LastExecution.Status,
				Error:      i.Status.LastExecution.Error,
				Tasks:      maskTasksSensitiveData(i.Status.LastExecution.Tasks),
			},
		})
	}

	return pagination.ChunkingPaginationResponse[Infra]{
		Items: infras,
		Chunk: l.Continue,
	}, nil
}

// Update implements UseCase.
func (u useCase) Update(ctx context.Context, infra Infra) (Infra, error) {
	newInfra := commonv1alpha1.Infra{
		Spec: infra.InfraSpec,
	}

	newInfra.SetName(infra.Name)
	newInfra.SetNamespace(infra.Namespace)

	s, err := u.repository.Apply(ctx, newInfra)
	if err != nil {
		return Infra{}, err
	}

	return Infra{
		Name:      s.GetName(),
		Namespace: s.GetNamespace(),
		InfraSpec: s.Spec,
		Status: InfraStatus{
			StartedAt:  s.Status.LastExecution.StartedAt,
			FinishedAt: s.Status.LastExecution.FinishedAt,
			Status:     s.Status.LastExecution.Status,
			Error:      s.Status.LastExecution.Error,
			Tasks:      maskTasksSensitiveData(s.Status.LastExecution.Tasks),
		},
	}, nil
}

func maskTasksSensitiveData(taskStatus []commonv1alpha1.TaskExecutionStatus) []InfraTaskStatus {
	maskedTasks := []InfraTaskStatus{}

	for _, p := range taskStatus {

		inputs := []commonv1alpha1.InfraTaskInput{}
		for _, i := range p.Inputs {
			if i.Sensitive {
				i.Value = "***"
			}

			inputs = append(inputs, i)
		}
		maskedTasks = append(maskedTasks, InfraTaskStatus{
			Name:       p.Name,
			Depends:    p.Depends,
			Backend:    p.Backend,
			Inputs:     inputs,
			StartedAt:  p.StartedAt,
			FinishedAt: p.FinishedAt,
			Status:     p.Status,
			Error:      p.Error,
		})
	}

	return maskedTasks
}
