package sharedinfra

import (
	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/pagination"
)

type SharedInfra struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	commonv1alpha1.SharedInfraSpec
}

type UseCase interface {
	List(chunkPagination pagination.ChunkingPaginationRequest) (pagination.ChunkingPaginationResponse[SharedInfra], error)
	Create(sharedInfra SharedInfra) (SharedInfra, error)
	Get(name string) (SharedInfra, error)
	Update(name string, sharedInfra SharedInfra) (SharedInfra, error)
	Delete(name string) error
}

type Repository interface {
	List() (commonv1alpha1.SharedInfraList, error)
	Apply(commonv1alpha1.SharedInfra) (commonv1alpha1.SharedInfra, error)
	Get(name string) (commonv1alpha1.SharedInfra, error)
	Delete(name string) error
}
