package sharedinfra

import (
	"context"
	"os"
	"time"

	"github.com/google/uuid"
	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/execution"
	"github.com/octopipe/cloudx/internal/runner"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
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

	executionId := uuid.New()
	newSharedInfraExecution := commonv1alpha1.SharedInfraExecutionStatus{
		Id:        executionId.String(),
		StartedAt: time.Now().Format(time.RFC3339),
		Status:    execution.ExecutionRunningStatus,
	}
	currentExecutions := currentSharedInfra.Status.Executions
	currentExecutions = append(currentExecutions, newSharedInfraExecution)
	currentSharedInfra.Status.Executions = currentExecutions

	err = c.Status().Update(ctx, currentSharedInfra)
	if err != nil {
		return ctrl.Result{}, err
	}

	if os.Getenv("ENV") != "local" {
		newRunner := runner.NewRunner(executionId.String(), req.String(), *currentSharedInfra)
		err = c.Create(ctx, newRunner.Pod)
		if err != nil {
			return ctrl.Result{}, err
		}

	}

	return ctrl.Result{}, nil
}

func (c *controller) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&commonv1alpha1.SharedInfra{}).
		WithEventFilter(predicate.Or(predicate.GenerationChangedPredicate{}, predicate.LabelChangedPredicate{})).
		Complete(c)
}
