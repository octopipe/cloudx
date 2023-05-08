package terraform

import (
	"context"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
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
	logger   *zap.Logger
	scheme   *runtime.Scheme
	execPath string
}

func NewController(logger *zap.Logger, client client.Client, scheme *runtime.Scheme, execPath string) Controller {
	return &controller{
		Client:   client,
		logger:   logger,
		scheme:   scheme,
		execPath: execPath,
	}
}

func (c *controller) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	currentStack := &commonv1alpha1.Stack{}
	err := c.Get(ctx, req.NamespacedName, currentStack)
	if err != nil {
		return ctrl.Result{}, nil
	}

	execution := NewExecution(c.execPath)
	execution.executeNextSteps(currentStack.Spec.Plugins)

	return ctrl.Result{}, nil
}

func (c *controller) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&commonv1alpha1.Stack{}).
		Complete(c)
}
