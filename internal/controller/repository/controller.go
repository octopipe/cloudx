package repository

import (
	"context"
	"time"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/repository"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
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
	scheme            *runtime.Scheme
	k8sClient         *kubernetes.Clientset
	repositoryUseCase repository.UseCase
}

func NewController(logger *zap.Logger, client client.Client, scheme *runtime.Scheme, k8sClient *kubernetes.Clientset, repositoryUseCase repository.UseCase) Controller {

	return &controller{
		Client:            client,
		logger:            logger,
		scheme:            scheme,
		k8sClient:         k8sClient,
		repositoryUseCase: repositoryUseCase,
	}
}

func (c *controller) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	repository := &commonv1alpha1.Repository{}

	err := c.Get(ctx, req.NamespacedName, repository)
	if err != nil {
		return ctrl.Result{}, err
	}

	c.logger.Info("syncing repository", zap.String("name", req.Name), zap.String("namespace", req.Namespace))
	err = c.repositoryUseCase.Sync(ctx, req.Name, req.Namespace)
	if err != nil {
		c.logger.Error("failed to sync repository", zap.Error(err))
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
}

func getControlResult(repository *commonv1alpha1.Repository) ctrl.Result {
	if repository.Spec.Sync.Auto {
		return ctrl.Result{RequeueAfter: 3 * time.Second}
	}

	return ctrl.Result{}
}

func (c *controller) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&commonv1alpha1.Repository{}).
		WithEventFilter(predicate.Or(predicate.GenerationChangedPredicate{}, predicate.LabelChangedPredicate{})).
		Complete(c)
}
