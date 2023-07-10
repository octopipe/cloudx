package infra

import (
	"context"
	"os"
	"time"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/controller/utils"
	"github.com/octopipe/cloudx/internal/engine"
	"github.com/octopipe/cloudx/internal/provider"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type Controller interface {
	Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error)
	SetupWithManager(mgr ctrl.Manager) error
}

type controller struct {
	client.Client
	logger   *zap.Logger
	scheme   *runtime.Scheme
	provider provider.Provider
}

func NewController(logger *zap.Logger, client client.Client, scheme *runtime.Scheme, provider provider.Provider) Controller {
	return &controller{
		Client:   client,
		logger:   logger,
		scheme:   scheme,
		provider: provider,
	}
}

func (c *controller) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	currentInfra := &commonv1alpha1.Infra{}
	err := c.Get(ctx, req.NamespacedName, currentInfra)
	if err != nil {
		return ctrl.Result{}, nil
	}

	if currentInfra.Status.LastExecution.Status == engine.ExecutionRunningStatus {
		c.logger.Info("This infra has runner in execution, enqueue this request")
		return ctrl.Result{
			RequeueAfter: time.Second * 2,
		}, nil
	}

	action := "APPLY"
	if len(currentInfra.Finalizers) > 0 {
		action = "DESTROY"
	}

	providerConfig := commonv1alpha1.ProviderConfig{}
	err = c.Get(ctx, types.NamespacedName{
		Name:      currentInfra.Spec.ProviderConfigRef.Name,
		Namespace: currentInfra.Spec.ProviderConfigRef.Namespace,
	}, &providerConfig)
	if err != nil {
		c.logger.Error("Failed to get provider config by shared infra", zap.Error(err))
		return ctrl.Result{
			RequeueAfter: time.Second * 2,
		}, err
	}

	if os.Getenv("ENV") != "local" {
		c.logger.Info("creating runner...")
		newRunner, err := c.NewRunner(action, *currentInfra, providerConfig)
		if err != nil {
			c.logger.Error("Failed to create runner", zap.Error(err))
			return ctrl.Result{Requeue: false}, err
		}

		err = c.Create(ctx, newRunner.Pod)
		if err != nil {
			c.logger.Error("Failed to apply runner", zap.Error(err))
			return ctrl.Result{Requeue: false}, err
		}

		currentInfra.Status.LastExecution.Status = engine.ExecutionRunningStatus
		currentInfra.Status.LastExecution.StartedAt = time.Now().Format(time.RFC3339)
		err = utils.UpdateInfraStatus(c.Client, *currentInfra)
		if err != nil {
			c.logger.Error("Failed to update infra status", zap.Error(err))
			return ctrl.Result{Requeue: false}, err
		}
	}

	return ctrl.Result{Requeue: false}, nil
}

func ignoreDeletionPredicate() predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			return e.ObjectOld.GetGeneration() != e.ObjectNew.GetGeneration()
		},
	}
}

func (c *controller) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&commonv1alpha1.Infra{}).
		WithEventFilter(predicate.Or(predicate.GenerationChangedPredicate{}, predicate.LabelChangedPredicate{})).
		Complete(c)
}