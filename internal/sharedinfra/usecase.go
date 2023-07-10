package sharedinfra

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
func (u useCase) Create(ctx context.Context, sharedInfra SharedInfra) (SharedInfra, error) {
	newSharedInfra := commonv1alpha1.SharedInfra{
		Spec: sharedInfra.SharedInfraSpec,
	}

	newSharedInfra.SetName(sharedInfra.Name)
	newSharedInfra.SetNamespace(sharedInfra.Namespace)

	s, err := u.repository.Apply(ctx, newSharedInfra)
	if err != nil {
		return SharedInfra{}, err
	}

	return SharedInfra{
		Name:            s.GetName(),
		Namespace:       s.GetNamespace(),
		SharedInfraSpec: s.Spec,
		Status: SharedInfraStatus{
			StartedAt:  s.Status.LastExecution.StartedAt,
			FinishedAt: s.Status.LastExecution.FinishedAt,
			Status:     s.Status.LastExecution.Status,
			Error:      s.Status.LastExecution.Error,
			Plugins:    maskPluginsSensitiveData(s.Status.LastExecution.Plugins),
		},
	}, nil
}

// Delete implements UseCase.
func (u useCase) Delete(ctx context.Context, name string, namespace string) error {
	return u.repository.Delete(ctx, name, namespace)
}

// Get implements UseCase.
func (u useCase) Get(ctx context.Context, name string, namespace string) (SharedInfra, error) {
	s, err := u.repository.Get(ctx, name, namespace)
	if err != nil {
		return SharedInfra{}, err
	}

	return SharedInfra{
		Name:            s.GetName(),
		Namespace:       s.GetNamespace(),
		SharedInfraSpec: s.Spec,
		Status: SharedInfraStatus{
			StartedAt:  s.Status.LastExecution.StartedAt,
			FinishedAt: s.Status.LastExecution.FinishedAt,
			Status:     s.Status.LastExecution.Status,
			Error:      s.Status.LastExecution.Error,
			Plugins:    maskPluginsSensitiveData(s.Status.LastExecution.Plugins),
		},
	}, nil
}

func (u useCase) Reconcile(ctx context.Context, name string, namespace string) error {
	return u.repository.Reconcile(ctx, name, namespace)
}

// List implements UseCase.
func (u useCase) List(ctx context.Context, namespace string, chunkPagination pagination.ChunkingPaginationRequest) (pagination.ChunkingPaginationResponse[SharedInfra], error) {
	l, err := u.repository.List(ctx, namespace, chunkPagination)
	if err != nil {
		return pagination.ChunkingPaginationResponse[SharedInfra]{}, err
	}

	sharedInfras := []SharedInfra{}
	for _, i := range l.Items {
		sharedInfras = append(sharedInfras, SharedInfra{
			Name:            i.GetName(),
			Namespace:       i.GetNamespace(),
			SharedInfraSpec: i.Spec,
			Status: SharedInfraStatus{
				StartedAt:  i.Status.LastExecution.StartedAt,
				FinishedAt: i.Status.LastExecution.FinishedAt,
				Status:     i.Status.LastExecution.Status,
				Error:      i.Status.LastExecution.Error,
				Plugins:    maskPluginsSensitiveData(i.Status.LastExecution.Plugins),
			},
		})
	}

	return pagination.ChunkingPaginationResponse[SharedInfra]{
		Items: sharedInfras,
		Chunk: l.Continue,
	}, nil
}

// Update implements UseCase.
func (u useCase) Update(ctx context.Context, sharedInfra SharedInfra) (SharedInfra, error) {
	newSharedInfra := commonv1alpha1.SharedInfra{
		Spec: sharedInfra.SharedInfraSpec,
	}

	newSharedInfra.SetName(sharedInfra.Name)
	newSharedInfra.SetNamespace(sharedInfra.Namespace)

	s, err := u.repository.Apply(ctx, newSharedInfra)
	if err != nil {
		return SharedInfra{}, err
	}

	return SharedInfra{
		Name:            s.GetName(),
		Namespace:       s.GetNamespace(),
		SharedInfraSpec: s.Spec,
		Status: SharedInfraStatus{
			StartedAt:  s.Status.LastExecution.StartedAt,
			FinishedAt: s.Status.LastExecution.FinishedAt,
			Status:     s.Status.LastExecution.Status,
			Error:      s.Status.LastExecution.Error,
			Plugins:    maskPluginsSensitiveData(s.Status.LastExecution.Plugins),
		},
	}, nil
}

func maskPluginsSensitiveData(pluginStatus []commonv1alpha1.PluginExecutionStatus) []commonv1alpha1.PluginExecutionStatus {
	maskedPlugins := []commonv1alpha1.PluginExecutionStatus{}

	for _, p := range pluginStatus {

		inputs := []commonv1alpha1.SharedInfraPluginInput{}
		for _, i := range p.Inputs {
			if i.Sensitive {
				i.Value = "***"
			}

			inputs = append(inputs, i)
		}

		p.Inputs = inputs
		maskedPlugins = append(maskedPlugins, p)
	}

	return maskedPlugins
}
