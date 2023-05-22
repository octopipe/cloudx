package sharedinfra

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/pluginmanager"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
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
	logger           *zap.Logger
	scheme           *runtime.Scheme
	pluginManager    pluginmanager.Manager
	mu               sync.Mutex
	executionContext executionContext
}

type executionContext struct {
	dependencyGraph map[string][]string
	executionGraph  map[string][]string
	executedNodes   map[string][]commonv1alpha1.SharedInfraPluginOutput
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
	currentSharedInfra := &commonv1alpha1.SharedInfra{}
	err := c.Get(ctx, req.NamespacedName, currentSharedInfra)
	if err != nil {
		return ctrl.Result{}, nil
	}

	dependencyGraph, executionGraph := createGraphs(*currentSharedInfra)
	newExecutionContext := executionContext{
		dependencyGraph: dependencyGraph,
		executionGraph:  executionGraph,
		executedNodes:   map[string][]commonv1alpha1.SharedInfraPluginOutput{},
	}
	pluginStatus, err := c.execute(&newExecutionContext, currentSharedInfra.Spec.Plugins)
	if err != nil {
		return ctrl.Result{}, nil
	}

	currentSharedInfra.Status.Plugins = pluginStatus
	err = c.Status().Update(ctx, currentSharedInfra)
	if err != nil {
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

func (c *controller) execute(executionContext *executionContext, plugins []commonv1alpha1.SharedInfraPlugin) ([]commonv1alpha1.PluginStatus, error) {
	status := []commonv1alpha1.PluginStatus{}
	eg, _ := errgroup.WithContext(context.Background())
	for _, p := range plugins {
		if _, ok := executionContext.executedNodes[p.Name]; !ok && isComplete(executionContext.dependencyGraph[p.Name], executionContext.executedNodes) {
			cb := func(currentPlugin commonv1alpha1.SharedInfraPlugin) func() error {
				return func() error {
					inputs := map[string]interface{}{}

					for _, i := range currentPlugin.Inputs {
						inputs[i.Key] = i.Value
					}

					if p.PluginType == "terraform" {
						out, state, err := c.pluginManager.ExecuteTerraformPlugin(currentPlugin.Ref, inputs)
						if err != nil {
							status = append(status, commonv1alpha1.PluginStatus{
								Name:            p.Name,
								State:           state,
								ExecutionStatus: "Error",
								ExecutionAt:     strconv.Itoa(int(time.Now().Unix())),
								Error:           err.Error(),
							})
							return err
						}

						c.mu.Lock()
						defer c.mu.Unlock()
						status = append(status, commonv1alpha1.PluginStatus{
							Name:            p.Name,
							State:           state,
							ExecutionStatus: "ExecutedWithSuccess",
							ExecutionAt:     strconv.Itoa(int(time.Now().Unix())),
						})
						executionContext.executedNodes[currentPlugin.Name] = out
					}
					return nil
				}
			}
			eg.Go(cb(p))
		}
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	fmt.Println(len(executionContext.executedNodes) == len(executionContext.executionGraph))

	if len(executionContext.executedNodes) == len(executionContext.executionGraph) {
		return status, nil
	}

	p, err := c.execute(executionContext, plugins)
	if err != nil {
		return nil, err
	}

	status = append(status, p...)
	return status, nil
}

func isComplete(dependencies []string, executedNodes map[string][]commonv1alpha1.SharedInfraPluginOutput) bool {
	isComplete := true

	for _, d := range dependencies {
		if _, ok := executedNodes[d]; !ok {
			isComplete = false
		}
	}

	return isComplete || len(dependencies) <= 0
}

func createGraphs(stackset commonv1alpha1.SharedInfra) (map[string][]string, map[string][]string) {

	dependencyGraph := map[string][]string{}
	executionGraph := map[string][]string{}

	for _, p := range stackset.Spec.Plugins {
		dependencyGraph[p.Name] = p.Depends
		executionGraph[p.Name] = []string{}
	}

	for _, p := range stackset.Spec.Plugins {
		for _, d := range p.Depends {
			executionGraph[d] = append(executionGraph[d], p.Name)
		}
	}

	return dependencyGraph, executionGraph
}

func (c *controller) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&commonv1alpha1.SharedInfra{}).
		Complete(c)
}
