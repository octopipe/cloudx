package sharedinfra

import (
	"github.com/octopipe/cloudx/apis/common/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type k8sRepository struct {
	client.Client
}

func NewK8sRepository() Repository {
	return k8sRepository{}
}

// Apply implements Repository.
func (k8sRepository) Apply(v1alpha1.SharedInfra) (v1alpha1.SharedInfra, error) {
	panic("unimplemented")
}

// Get implements Repository.
func (k8sRepository) Get(name string) (v1alpha1.SharedInfra, error) {
	panic("unimplemented")
}

// Delete implements Repository.
func (k8sRepository) Delete(name string) error {
	panic("unimplemented")
}

// List implements Repository.
func (k8sRepository) List() (v1alpha1.SharedInfraList, error) {
	panic("unimplemented")
}
