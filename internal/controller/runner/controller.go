package runner

import (
	"context"

	"github.com/octopipe/cloudx/internal/pluginmanager"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
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

	currentRunner := &v1.Pod{}
	err := c.Get(ctx, req.NamespacedName, currentRunner)
	if err != nil {
		return ctrl.Result{}, nil
	}

	labels := currentRunner.Labels
	sharedInfraName := labels["commons.cloudx.io/sharedinfra-name"]

	// fmt.Println(currentJob.Status.Conditions)
	c.logger.Info("runner request", zap.String("name", currentRunner.GetName()), zap.String("sharedinfra", sharedInfraName))

	return ctrl.Result{}, nil
}

func (c *controller) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.Pod{}).
		WithEventFilter(ignoreNonControlledPods()).
		Complete(c)
}

func ignoreNonControlledPods() predicate.Predicate {
	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			labels := e.Object.GetLabels()
			value, ok := labels["app.kubernetes.io/managed-by"]
			return labels != nil && ok && value == "cloudx"
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			labels := e.ObjectNew.GetLabels()
			value, ok := labels["app.kubernetes.io/managed-by"]
			return labels != nil && ok && value == "cloudx"

		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			labels := e.Object.GetLabels()
			value, ok := labels["app.kubernetes.io/managed-by"]
			return labels != nil && ok && value == "cloudx"
		},
	}
}
