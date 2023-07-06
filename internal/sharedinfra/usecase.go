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
		Status:          s.Status,
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
		Status:          s.Status,
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
			Status:          i.Status,
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
		Status:          s.Status,
	}, nil
}
