package connectioninterface

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
func (u useCase) Create(ctx context.Context, connectionInterface ConnectionInterface) (ConnectionInterface, error) {
	newConnectionInterface := commonv1alpha1.ConnectionInterface{
		Spec: connectionInterface.ConnectionInterfaceSpec,
	}

	newConnectionInterface.SetName(connectionInterface.Name)
	newConnectionInterface.SetNamespace(connectionInterface.Namespace)

	s, err := u.repository.Apply(ctx, newConnectionInterface)
	if err != nil {
		return ConnectionInterface{}, err
	}

	return ConnectionInterface{
		Name:                    s.GetName(),
		Namespace:               s.GetNamespace(),
		ConnectionInterfaceSpec: s.Spec,
	}, nil
}

// Delete implements UseCase.
func (u useCase) Delete(ctx context.Context, name string, namespace string) error {
	return u.repository.Delete(ctx, name, namespace)
}

// Get implements UseCase.
func (u useCase) Get(ctx context.Context, name string, namespace string) (ConnectionInterface, error) {
	s, err := u.repository.Get(ctx, name, namespace)
	if err != nil {
		return ConnectionInterface{}, err
	}

	return ConnectionInterface{
		Name:                    s.GetName(),
		Namespace:               s.GetNamespace(),
		ConnectionInterfaceSpec: s.Spec,
	}, nil
}

// List implements UseCase.
func (u useCase) List(ctx context.Context, namespace string, chunkPagination pagination.ChunkingPaginationRequest) (pagination.ChunkingPaginationResponse[ConnectionInterface], error) {
	l, err := u.repository.List(ctx, namespace, chunkPagination)
	if err != nil {
		return pagination.ChunkingPaginationResponse[ConnectionInterface]{}, err
	}

	connectionInterfaces := []ConnectionInterface{}
	for _, i := range l.Items {
		connectionInterfaces = append(connectionInterfaces, ConnectionInterface{
			Name:                    i.GetName(),
			Namespace:               i.GetNamespace(),
			ConnectionInterfaceSpec: i.Spec,
		})
	}

	return pagination.ChunkingPaginationResponse[ConnectionInterface]{
		Items: connectionInterfaces,
		Chunk: l.Continue,
	}, nil
}

// Update implements UseCase.
func (u useCase) Update(ctx context.Context, connectionInterface ConnectionInterface) (ConnectionInterface, error) {
	newConnectionInterface := commonv1alpha1.ConnectionInterface{
		Spec: connectionInterface.ConnectionInterfaceSpec,
	}

	newConnectionInterface.SetName(connectionInterface.Name)
	newConnectionInterface.SetNamespace(connectionInterface.Namespace)

	s, err := u.repository.Apply(ctx, newConnectionInterface)
	if err != nil {
		return ConnectionInterface{}, err
	}

	return ConnectionInterface{
		Name:                    s.GetName(),
		Namespace:               s.GetNamespace(),
		ConnectionInterfaceSpec: s.Spec,
	}, nil
}
