package sharedinfra

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/google/uuid"
	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/execution"
	"github.com/octopipe/cloudx/internal/runner"
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

	rawSharedInfra, err := json.Marshal(currentSharedInfra)
	if err != nil {
		return ctrl.Result{}, err
	}

	action := "APPLY"
	if len(currentSharedInfra.Finalizers) > 0 {
		action = "DESTROY"
	}

	c.logger.Info("reconcile shared-infra", zap.String("shared-infra", currentSharedInfra.GetName()), zap.String("action", action))

	providerConfig := commonv1alpha1.ProviderConfig{}
	err = c.Get(ctx, types.NamespacedName{
		Name:      currentSharedInfra.Spec.ProviderConfigRef.Name,
		Namespace: currentSharedInfra.Spec.ProviderConfigRef.Namespace,
	}, &providerConfig)
	if err != nil {
		return ctrl.Result{}, err
	}

	if os.Getenv("ENV") != "local" {
		c.logger.Info("creating runner")
		newRunner, err := runner.NewRunner(executionId.String(), *currentSharedInfra, string(rawSharedInfra), action, providerConfig)
		if err != nil {
			return ctrl.Result{}, err
		}

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
