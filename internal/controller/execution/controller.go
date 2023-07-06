package execution

import (
	"context"
	"encoding/json"
	"os"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/engine"
	"github.com/octopipe/cloudx/internal/provider"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
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
	currentExecution := &commonv1alpha1.Execution{}
	err := c.Get(ctx, req.NamespacedName, currentExecution)
	if err != nil {
		return ctrl.Result{}, nil
	}

	if currentExecution.Status.Status == engine.ExecutionSuccessStatus || currentExecution.Status.Status == engine.ExecutionErrorStatus {
		return ctrl.Result{}, nil
	}

	c.logger.Info("starting new execution...")

	currentSharedInfra := &commonv1alpha1.SharedInfra{}
	sharedInfraRef := types.NamespacedName{
		Name:      currentExecution.Spec.SharedInfra.Name,
		Namespace: currentExecution.Spec.SharedInfra.Namespace,
	}
	err = c.Get(ctx, sharedInfraRef, currentSharedInfra)
	if err != nil {
		c.logger.Error("Failed to get execution shared-infra", zap.Error(err))
		return ctrl.Result{}, err
	}

	rawSharedInfra, err := json.Marshal(currentSharedInfra)
	if err != nil {
		c.logger.Error("Failed to parse shared-infra", zap.Error(err))
		return ctrl.Result{}, err
	}

	escapedSharedInfra, err := json.Marshal(rawSharedInfra)
	if err != nil {
		c.logger.Error("Failed to parse shared-infra", zap.Error(err))
		return ctrl.Result{}, err
	}

	providerConfig := commonv1alpha1.ProviderConfig{}
	err = c.Get(ctx, types.NamespacedName{
		Name:      currentSharedInfra.Spec.ProviderConfigRef.Name,
		Namespace: currentSharedInfra.Spec.ProviderConfigRef.Namespace,
	}, &providerConfig)
	if err != nil {
		c.logger.Error("Failed to get provider config by shared infra", zap.Error(err))
		return ctrl.Result{}, err
	}

	if os.Getenv("ENV") != "local" {
		c.logger.Info("creating runner...")
		newRunner, err := c.NewRunner(*currentExecution, *currentSharedInfra, string(escapedSharedInfra), providerConfig)
		if err != nil {

			c.logger.Error("Failed to create runner", zap.Error(err))
			return ctrl.Result{}, err
		}

		err = c.Create(ctx, newRunner.Pod)
		if err != nil {
			c.logger.Error("Failed to apply runner", zap.Error(err))
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (c *controller) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&commonv1alpha1.Execution{}).
		WithEventFilter(predicate.Or(predicate.GenerationChangedPredicate{}, predicate.LabelChangedPredicate{})).
		Complete(c)
}
