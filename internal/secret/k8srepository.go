package secret

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type k8sRepository struct {
	client client.Client
}

func NewK8sRepository(c client.Client) Repository {
	return k8sRepository{client: c}
}

// Apply implements Repository.
func (r k8sRepository) Apply(ctx context.Context, s v1.Secret) (v1.Secret, error) {
	err := r.client.Create(ctx, &s)
	if err != nil && errors.IsAlreadyExists(err) {
		err = retry.RetryOnConflict(retry.DefaultBackoff, func() error {
			current := v1.Secret{}
			err = r.client.Get(ctx, types.NamespacedName{
				Name:      s.Name,
				Namespace: s.Namespace,
			}, &current)
			if err != nil {
				return err
			}

			current.Data = s.Data

			return r.client.Update(ctx, &current)
		})

	}

	return s, err
}

// Get implements Repository.
func (r k8sRepository) Get(ctx context.Context, name string, namespace string) (v1.Secret, error) {
	var secret v1.Secret
	err := r.client.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, &secret)
	return secret, err
}

// Delete implements Repository.
func (r k8sRepository) Delete(ctx context.Context, name string, namespace string) error {
	secret, err := r.Get(ctx, name, namespace)
	if err != nil {
		return nil
	}

	err = r.client.Delete(ctx, &secret)
	return err
}

// List implements Repository.
// func (r k8sRepository) List(ctx context.Context, namespace string, chunkPagination pagination.ChunkingPaginationRequest) (v1alpha1.RepositoryList, error) {
// 	var secretList v1alpha1.RepositoryList
// 	err := r.client.List(ctx, &secretList, &client.ListOptions{Limit: chunkPagination.Limit, Continue: chunkPagination.Chunk, Namespace: namespace})
// 	return secretList, err
// }
