package sharedinfra

import (
	"context"
	"fmt"
	"time"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/engine"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	logger *zap.Logger
	scheme *runtime.Scheme
}

func NewController(logger *zap.Logger, client client.Client, scheme *runtime.Scheme) Controller {
	return &controller{
		Client: client,
		logger: logger,
		scheme: scheme,
	}
}

func (c *controller) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	currentSharedInfra := &commonv1alpha1.SharedInfra{}
	err := c.Get(ctx, req.NamespacedName, currentSharedInfra)
	if err != nil {
		return ctrl.Result{}, nil
	}

	action := "APPLY"
	if len(currentSharedInfra.Finalizers) > 0 {
		action = "DESTROY"
	}

	newExecution := commonv1alpha1.Execution{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("execution-%s-%d", currentSharedInfra.GetName(), time.Now().UnixMilli()),
			Namespace: currentSharedInfra.GetNamespace(),
			Labels: map[string]string{
				"commons.cloudx.io/sharedinfra-name":      currentSharedInfra.GetName(),
				"commons.cloudx.io/sharedinfra-namespace": currentSharedInfra.GetNamespace(),
			},
		},
		Spec: commonv1alpha1.ExecutionSpec{
			Action: action,
			SharedInfra: commonv1alpha1.Ref{
				Name:      currentSharedInfra.GetName(),
				Namespace: currentSharedInfra.GetNamespace(),
			},
			StartedAt: time.Now().Format(time.RFC3339),
		},
	}

	hasExecutionRunning, err := c.hasExecutionRunning(ctx, req)
	if err != nil {
		return ctrl.Result{}, err
	}

	if hasExecutionRunning {
		c.logger.Info("the shared infra has a execution in status running", zap.String("shared-infra", currentSharedInfra.Name))
		return ctrl.Result{}, nil
	}

	c.logger.Info("creating new execution", zap.String("shared-infra", currentSharedInfra.Name))
	err = c.Create(ctx, &newExecution)
	if err != nil {
		return ctrl.Result{}, err
	}

	// newExecution.Status = commonv1alpha1.ExecutionStatus{

	// 	Status: engine.ExecutionRunningStatus,
	// }

	// err = utils.UpdateExecutionStatus(c.Client, newExecution)
	// if err != nil {
	// 	return ctrl.Result{}, err
	// }

	return ctrl.Result{}, nil
}

func (c *controller) hasExecutionRunning(ctx context.Context, sharedInfraRef ctrl.Request) (bool, error) {
	executionList := commonv1alpha1.ExecutionList{}

	err := c.List(ctx, &executionList)
	if err != nil {
		return false, err
	}

	for _, i := range executionList.Items {
		if i.Spec.SharedInfra.Name == sharedInfraRef.Name && i.Spec.SharedInfra.Namespace == sharedInfraRef.Namespace && i.Status.Status == engine.ExecutionRunningStatus {
			return true, nil
		}
	}

	return false, nil
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
		For(&commonv1alpha1.SharedInfra{}).
		WithEventFilter(predicate.Or(predicate.GenerationChangedPredicate{}, predicate.LabelChangedPredicate{})).
		Complete(c)
}
