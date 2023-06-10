package sharedinfra

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	newRunner := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-runner-%d", currentSharedInfra.GetName(), time.Now().Unix()),
			Namespace: "cloudx-system",
			Labels: map[string]string{
				"commons.cloudx.io/sharedinfra-name":      currentSharedInfra.GetName(),
				"commons.cloudx.io/sharedinfra-namespace": currentSharedInfra.GetNamespace(),
				"commons.cloudx.io/execution-id":          executionId.String(),
				"app.kubernetes.io/managed-by":            "cloudx",
			},
		},
		Spec: v1.PodSpec{
			ServiceAccountName: "controller-sa",
			RestartPolicy:      v1.RestartPolicyNever,
			Containers: []v1.Container{
				{
					Name:            "runner",
					Image:           "mayconjrpacheco/cloudx-runner:latest",
					Args:            []string{req.String(), executionId.String()},
					ImagePullPolicy: v1.PullAlways,
					Env: []v1.EnvVar{
						{
							Name:  "TF_VERSION",
							Value: "latest",
						},
						{
							Name:  "RPC_SERVER_ADDRESS",
							Value: "controller.cloudx-system:9000",
						},
					},
				},
			},
		},
	}

	newSharedInfraExecution := commonv1alpha1.SharedInfraExecutionStatus{
		Id:        executionId.String(),
		StartedAt: time.Now().Format(time.RFC3339),
		Status:    "RUNNING",
	}
	currentExecutions := currentSharedInfra.Status.Executions
	currentExecutions = append(currentExecutions, newSharedInfraExecution)
	currentSharedInfra.Status.Executions = currentExecutions

	err = c.Status().Update(ctx, currentSharedInfra)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = c.Create(ctx, newRunner)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (c *controller) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&commonv1alpha1.SharedInfra{}).
		WithEventFilter(predicate.Or(predicate.GenerationChangedPredicate{}, predicate.LabelChangedPredicate{})).
		Complete(c)
}
