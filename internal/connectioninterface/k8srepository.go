package connectioninterface

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
func (r k8sRepository) Apply(ctx context.Context, s v1alpha1.ConnectionInterface) (v1alpha1.ConnectionInterface, error) {
	err := r.client.Create(ctx, &s)
	return s, err
}

// Get implements Repository.
func (r k8sRepository) Get(ctx context.Context, name string, namespace string) (v1alpha1.ConnectionInterface, error) {
	var connectionInterface v1alpha1.ConnectionInterface
	err := r.client.Get(ctx, types.NamespacedName{Name: name, Namespace: namespace}, &connectionInterface)
	return connectionInterface, err
}

// Delete implements Repository.
func (r k8sRepository) Delete(ctx context.Context, name string, namespace string) error {
	connectionInterface, err := r.Get(ctx, name, namespace)
	if err != nil {
		return nil
	}

	err = r.client.Delete(ctx, &connectionInterface)
	return err
}

// List implements Repository.
func (r k8sRepository) List(ctx context.Context, namespace string, chunkPagination pagination.ChunkingPaginationRequest) (v1alpha1.ConnectionInterfaceList, error) {
	var connectionInterfaceList v1alpha1.ConnectionInterfaceList
	err := r.client.List(ctx, &connectionInterfaceList, &client.ListOptions{Limit: chunkPagination.Limit, Continue: chunkPagination.Chunk, Namespace: namespace})
	return connectionInterfaceList, err
}
