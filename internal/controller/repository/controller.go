package repository

import (
	"context"
	"time"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/annotation"
	"github.com/octopipe/cloudx/internal/repository"
	"github.com/octopipe/cloudx/pkg/twice/reconciler"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type Controller interface {
	Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error)
	SetupWithManager(mgr ctrl.Manager) error
}

type controller struct {
	client.Client
	logger            *zap.Logger
	repositoryUseCase repository.UseCase
	gitopsReconciler  reconciler.Reconciler
}

const (
	RepositoryNameAnnotation      = "cloudx.octopipe.io/repository-name"
	RepositoryNamespaceAnnotation = "cloudx.octopipe.io/repository-namespace"
)

func NewController(logger *zap.Logger, client client.Client, repositoryUseCase repository.UseCase, reconciler reconciler.Reconciler) Controller {

	return &controller{
		Client:            client,
		logger:            logger,
		repositoryUseCase: repositoryUseCase,
		gitopsReconciler:  reconciler,
	}
}

func (c *controller) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	currRepository := &commonv1alpha1.Repository{}

	err := c.Get(ctx, req.NamespacedName, currRepository)
	if err != nil {
		return ctrl.Result{}, err
	}

	c.logger.Info("syncing repository", zap.String("name", req.Name), zap.String("namespace", req.Namespace))
	manifests, err := c.repositoryUseCase.Sync(ctx, req.Name, req.Namespace)
	if err != nil {
		c.logger.Error("failed to sync repository", zap.Error(err))
		return ctrl.Result{}, err
	}

	c.logger.Info("plan repository files")
	planResult, err := c.gitopsReconciler.Plan(ctx, manifests, "default", func(un *unstructured.Unstructured) bool {
		currAnnotations := un.GetAnnotations()
		isSameRepositoryName := currAnnotations[RepositoryNameAnnotation] == currRepository.Name
		isSameRepositoryNamespace := currAnnotations[RepositoryNamespaceAnnotation] == currRepository.Namespace

		return isSameRepositoryName && isSameRepositoryNamespace
	})
	if err != nil {
		return ctrl.Result{}, err
	}

	_, err = c.gitopsReconciler.Apply(ctx, planResult, "default", map[string]string{
		RepositoryNameAnnotation:       currRepository.Name,
		RepositoryNamespaceAnnotation:  currRepository.Namespace,
		annotation.ManagedByAnnotation: "cloudx",
	})
	if err != nil {
		c.logger.Error("failed to apply plan result", zap.Error(err))
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
}

func (c *controller) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&commonv1alpha1.Repository{}).
		WithEventFilter(predicate.Or(predicate.GenerationChangedPredicate{}, predicate.LabelChangedPredicate{})).
		Complete(c)
}
