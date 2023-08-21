package repository

import (
	"context"

	"github.com/octopipe/cloudx/apis/common/v1alpha1"
	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/pagination"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type k8sRepository struct {
	client client.Client
}

func NewK8sRepository(c client.Client) RepositoryType {
	return k8sRepository{client: c}
}

// Apply implements Repository.
func (r k8sRepository) Apply(ctx context.Context, s v1alpha1.Repository) (v1alpha1.Repository, error) {

	err := r.client.Create(ctx, &s)
	if err != nil && errors.IsAlreadyExists(err) {
		err = retry.RetryOnConflict(retry.DefaultBackoff, func() error {
			current := commonv1alpha1.Repository{}
			err = r.client.Get(ctx, types.NamespacedName{
				Name:      s.Name,
				Namespace: s.Namespace,
			}, &current)
			if err != nil {
				return err
			}

			current.Spec = s.Spec

			return r.client.Update(ctx, &current)
		})

	}

	return s, err
}

// Get implements Repository.
func (r k8sRepository) Get(ctx context.Context, name string, namespace string) (v1alpha1.Repository, error) {
	var repository v1alpha1.Repository
	err := r.client.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, &repository)
	return repository, err
}

func (r k8sRepository) Sync(ctx context.Context, name string, namespace string) error {
	return retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		current := commonv1alpha1.Repository{}
		err := r.client.Get(ctx, types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		}, &current)
		if err != nil {
			return err
		}

		return r.client.Update(ctx, &current)
	})
}

// Delete implements Repository.
func (r k8sRepository) Delete(ctx context.Context, name string, namespace string) error {
	repository, err := r.Get(ctx, name, namespace)
	if err != nil {
		return nil
	}

	err = r.client.Delete(ctx, &repository)
	return err
}

// List implements Repository.
func (r k8sRepository) List(ctx context.Context, namespace string, chunkPagination pagination.ChunkingPaginationRequest) (v1alpha1.RepositoryList, error) {
	var repositoryList v1alpha1.RepositoryList
	err := r.client.List(ctx, &repositoryList, &client.ListOptions{Limit: chunkPagination.Limit, Continue: chunkPagination.Chunk, Namespace: namespace})
	return repositoryList, err
}
