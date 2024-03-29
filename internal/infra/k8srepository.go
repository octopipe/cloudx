package infra

import (
	"context"

	"github.com/google/uuid"
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

func NewK8sRepository(c client.Client) Repository {
	return k8sRepository{client: c}
}

// Apply implements Repository.
func (r k8sRepository) Apply(ctx context.Context, s v1alpha1.Infra) (v1alpha1.Infra, error) {
	err := r.client.Create(ctx, &s)
	if err != nil && errors.IsAlreadyExists(err) {
		err = retry.RetryOnConflict(retry.DefaultBackoff, func() error {
			current := commonv1alpha1.Infra{}
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
func (r k8sRepository) Get(ctx context.Context, name string, namespace string) (v1alpha1.Infra, error) {
	var infra v1alpha1.Infra
	err := r.client.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, &infra)
	return infra, err
}

func (r k8sRepository) Reconcile(ctx context.Context, name string, namespace string) error {
	return retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		current := commonv1alpha1.Infra{}
		err := r.client.Get(ctx, types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		}, &current)
		if err != nil {
			return err
		}

		current.Spec.Generation = uuid.NewString()

		return r.client.Update(ctx, &current)
	})
}

// Delete implements Repository.
func (r k8sRepository) Delete(ctx context.Context, name string, namespace string) error {
	infra, err := r.Get(ctx, name, namespace)
	if err != nil {
		return nil
	}

	err = r.client.Delete(ctx, &infra)
	return err
}

// List implements Repository.
func (r k8sRepository) List(ctx context.Context, namespace string, chunkPagination pagination.ChunkingPaginationRequest) (v1alpha1.InfraList, error) {
	var infraList v1alpha1.InfraList
	err := r.client.List(ctx, &infraList, &client.ListOptions{Limit: chunkPagination.Limit, Continue: chunkPagination.Chunk, Namespace: namespace})
	return infraList, err
}
