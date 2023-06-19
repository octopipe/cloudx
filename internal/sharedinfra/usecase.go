package sharedinfra

import (
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
func (u useCase) Create(sharedInfra SharedInfra) (SharedInfra, error) {
	newSharedInfra := commonv1alpha1.SharedInfra{
		Spec: sharedInfra.SharedInfraSpec,
	}

	newSharedInfra.SetName(sharedInfra.Name)
	newSharedInfra.SetNamespace(sharedInfra.Namespace)

	s, err := u.repository.Apply(newSharedInfra)
	if err != nil {
		return SharedInfra{}, err
	}

	return SharedInfra{
		Name:            s.GetName(),
		Namespace:       s.GetNamespace(),
		SharedInfraSpec: s.Spec,
	}, nil
}

// Delete implements UseCase.
func (u useCase) Delete(name string) error {
	return u.repository.Delete(name)
}

// Get implements UseCase.
func (u useCase) Get(name string) (SharedInfra, error) {
	s, err := u.repository.Get(name)
	if err != nil {
		return SharedInfra{}, err
	}

	return SharedInfra{
		Name:            s.GetName(),
		Namespace:       s.GetNamespace(),
		SharedInfraSpec: s.Spec,
	}, nil
}

// List implements UseCase.
func (u useCase) List() ([]pagination.ChunkingPaginationResponse[SharedInfra], error) {
	var sharedInfra SharedInfra
	l, err := u.repository.List()
	if err != nil {
		return pagination.ChunkingPaginationResponse[T sharedInfra]{}, err
	}


	for _, i := range l.Items {
		
	}
}

// Update implements UseCase.
func (useCase) Update(name string, sharedInfra SharedInfra) (SharedInfra, error) {
	panic("unimplemented")
}
