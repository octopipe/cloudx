package sharedinfra

import (
	"context"

	"github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/pagination"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type k8sRepository struct {
	client client.Client
}

func NewK8sRepository(c client.Client) Repository {
	return k8sRepository{client: c}
}

// Apply implements Repository.
func (r k8sRepository) Apply(ctx context.Context, s v1alpha1.SharedInfra) (v1alpha1.SharedInfra, error) {
	err := r.client.Create(ctx, &s)
	return s, err
}

// Get implements Repository.
func (r k8sRepository) Get(ctx context.Context, name string, namespace string) (v1alpha1.SharedInfra, error) {
	var sharedInfra v1alpha1.SharedInfra
	err := r.client.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, &sharedInfra)
	return sharedInfra, err
}

// Delete implements Repository.
func (r k8sRepository) Delete(ctx context.Context, name string, namespace string) error {
	sharedInfra, err := r.Get(ctx, name, namespace)
	if err != nil {
		return nil
	}

	err = r.client.Delete(ctx, &sharedInfra)
	return err
}

// List implements Repository.
func (r k8sRepository) List(ctx context.Context, namespace string, chunkPagination pagination.ChunkingPaginationRequest) (v1alpha1.SharedInfraList, error) {
	var sharedInfraList v1alpha1.SharedInfraList
	err := r.client.List(ctx, &sharedInfraList, &client.ListOptions{Limit: chunkPagination.Limit, Continue: chunkPagination.Chunk, Namespace: namespace})
	return sharedInfraList, err
}
