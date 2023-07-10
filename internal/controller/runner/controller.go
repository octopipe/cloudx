package runner

import (
	"bytes"
	"context"
	"io"

	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
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
	logger    *zap.Logger
	scheme    *runtime.Scheme
	k8sClient *kubernetes.Clientset
}

func NewController(logger *zap.Logger, client client.Client, scheme *runtime.Scheme, k8sClient *kubernetes.Clientset) Controller {

	return &controller{
		Client:    client,
		logger:    logger,
		scheme:    scheme,
		k8sClient: k8sClient,
	}
}

func (c *controller) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	currentRunner := &v1.Pod{}
	err := c.Get(ctx, req.NamespacedName, currentRunner)
	if err != nil {
		return ctrl.Result{}, err
	}

	// labels := currentRunner.Labels
	// infraName := labels["commons.cloudx.io/infra-name"]

	logsReq := c.k8sClient.CoreV1().Pods("cloudx-system").GetLogs(currentRunner.GetName(), &v1.PodLogOptions{})
	podLogs, err := logsReq.Stream(ctx)
	if err != nil {
		return ctrl.Result{}, err
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return ctrl.Result{}, err
	}

	logsStr := buf.String()

	c.logger.Info("runner logss", zap.String("name", currentRunner.GetName()), zap.String("logs", logsStr))

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
