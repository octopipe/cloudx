package secret

import (
	"context"

	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
)

type useCase struct {
	logger     *zap.Logger
	repository Repository
}

func NewUseCase(logger *zap.Logger, repository Repository) UseCase {
	return useCase{logger: logger, repository: repository}
}

// Create implements UseCase.
func (u useCase) Apply(ctx context.Context, s Secret) (Secret, error) {
	newSecret := v1.Secret{
		Data: s.Data,
	}

	newSecret.SetName(s.Name)
	newSecret.SetNamespace(s.Namespace)

	appliedSecret, err := u.repository.Apply(ctx, newSecret)
	if err != nil {
		return Secret{}, err
	}

	return Secret{
		Name:      appliedSecret.GetName(),
		Namespace: appliedSecret.GetNamespace(),
		Data:      appliedSecret.Data,
	}, nil
}

// Delete implements UseCase.
func (u useCase) Delete(ctx context.Context, name string, namespace string) error {
	return u.repository.Delete(ctx, name, namespace)
}

// Get implements UseCase.
func (u useCase) Get(ctx context.Context, name string, namespace string) (Secret, error) {
	s, err := u.repository.Get(ctx, name, namespace)
	if err != nil {
		return Secret{}, err
	}

	return Secret{
		Name:      s.GetName(),
		Namespace: s.GetNamespace(),
		Data:      s.Data,
	}, nil
}

// List implements UseCase.
// func (u useCase) List(ctx context.Context, namespace string, chunkPagination pagination.ChunkingPaginationRequest) (pagination.ChunkingPaginationResponse[Repository], error) {
// 	l, err := u.repository.List(ctx, namespace, chunkPagination)
// 	if err != nil {
// 		return pagination.ChunkingPaginationResponse[Repository]{}, err
// 	}

// 	repositorys := []Repository{}
// 	for _, i := range l.Items {
// 		repositorys = append(repositorys, Repository{
// 			Name:           i.GetName(),
// 			Namespace:      i.GetNamespace(),
// 			RepositorySpec: i.Spec,
// 		})
// 	}

// 	return pagination.ChunkingPaginationResponse[Repository]{
// 		Items: repositorys,
// 		Chunk: l.Continue,
// 	}, nil
// }
