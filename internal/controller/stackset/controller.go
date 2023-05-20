package stackset

import (
	"context"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/execution"
	"github.com/octopipe/cloudx/internal/pluginmanager"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Controller interface {
	Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error)
	SetupWithManager(mgr ctrl.Manager) error
}

type controller struct {
	client.Client
	logger        *zap.Logger
	scheme        *runtime.Scheme
	pluginManager pluginmanager.Manager
}

func NewController(logger *zap.Logger, client client.Client, scheme *runtime.Scheme, pluginManager pluginmanager.Manager) Controller {
	return &controller{
		Client:        client,
		logger:        logger,
		scheme:        scheme,
		pluginManager: pluginManager,
	}
}

func (c *controller) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	currentStackSet := &commonv1alpha1.StackSet{}
	err := c.Get(ctx, req.NamespacedName, currentStackSet)
	if err != nil {
		return ctrl.Result{}, nil
	}

	execution.NewExecutionManager(c.pluginManager, *currentStackSet)

	return ctrl.Result{}, nil
}

func (c *controller) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&commonv1alpha1.StackSet{}).
		Complete(c)
}
