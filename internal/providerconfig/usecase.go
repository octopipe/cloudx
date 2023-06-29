package providerconfig

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
func (u useCase) Create(ctx context.Context, providerConfig ProviderConfig) (ProviderConfig, error) {
	newProviderConfig := commonv1alpha1.ProviderConfig{
		Spec: providerConfig.ProviderConfigSpec,
	}

	newProviderConfig.SetName(providerConfig.Name)
	newProviderConfig.SetNamespace(providerConfig.Namespace)

	s, err := u.repository.Apply(ctx, newProviderConfig)
	if err != nil {
		return ProviderConfig{}, err
	}

	return ProviderConfig{
		Name:               s.GetName(),
		Namespace:          s.GetNamespace(),
		ProviderConfigSpec: s.Spec,
	}, nil
}

// Delete implements UseCase.
func (u useCase) Delete(ctx context.Context, name string, namespace string) error {
	return u.repository.Delete(ctx, name, namespace)
}

// Get implements UseCase.
func (u useCase) Get(ctx context.Context, name string, namespace string) (ProviderConfig, error) {
	s, err := u.repository.Get(ctx, name, namespace)
	if err != nil {
		return ProviderConfig{}, err
	}

	return ProviderConfig{
		Name:               s.GetName(),
		Namespace:          s.GetNamespace(),
		ProviderConfigSpec: s.Spec,
	}, nil
}

// List implements UseCase.
func (u useCase) List(ctx context.Context, namespace string, chunkPagination pagination.ChunkingPaginationRequest) (pagination.ChunkingPaginationResponse[ProviderConfig], error) {
	l, err := u.repository.List(ctx, namespace, chunkPagination)
	if err != nil {
		return pagination.ChunkingPaginationResponse[ProviderConfig]{}, err
	}

	providerConfigs := []ProviderConfig{}
	for _, i := range l.Items {
		providerConfigs = append(providerConfigs, ProviderConfig{
			Name:               i.GetName(),
			Namespace:          i.GetNamespace(),
			ProviderConfigSpec: i.Spec,
		})
	}

	return pagination.ChunkingPaginationResponse[ProviderConfig]{
		Items: providerConfigs,
		Chunk: l.Continue,
	}, nil
}

// Update implements UseCase.
func (u useCase) Update(ctx context.Context, providerConfig ProviderConfig) (ProviderConfig, error) {
	newProviderConfig := commonv1alpha1.ProviderConfig{
		Spec: providerConfig.ProviderConfigSpec,
	}

	newProviderConfig.SetName(providerConfig.Name)
	newProviderConfig.SetNamespace(providerConfig.Namespace)

	s, err := u.repository.Apply(ctx, newProviderConfig)
	if err != nil {
		return ProviderConfig{}, err
	}

	return ProviderConfig{
		Name:               s.GetName(),
		Namespace:          s.GetNamespace(),
		ProviderConfigSpec: s.Spec,
	}, nil
}
