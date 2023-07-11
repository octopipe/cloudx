package engine

import (
	"sync"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/backend"
	"github.com/octopipe/cloudx/internal/rpcclient"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type ExecutionOutputItem struct {
	Value     string
	Sensitive bool
	Type      string
}

type ExecutionContext map[string]map[string]ExecutionOutputItem

type Engine struct {
	logger    *zap.Logger
	backend   backend.Backend
	rpcClient rpcclient.Client

	mu               sync.Mutex
	executionContext ExecutionContext
}

type ActionFuncType func(taskName string, executionOutputs ExecutionContext) (commonv1alpha1.TaskExecutionStatus, map[string]ExecutionOutputItem)

func NewEngine(logger *zap.Logger, rpcClient rpcclient.Client, backend backend.Backend) Engine {
	return Engine{
		logger:           logger,
		backend:          backend,
		rpcClient:        rpcClient,
		executionContext: make(ExecutionContext),
	}
}

func (e *Engine) Run(graph map[string][]string, action ActionFuncType, taskStatusChan chan commonv1alpha1.TaskExecutionStatus) {
	eg := new(errgroup.Group)
	inDegrees := make(map[string]int)

	for node, deps := range graph {
		inDegrees[node] = len(deps)
	}

	for {
		for node, deps := range inDegrees {
			if _, ok := e.executionContext[node]; !ok && deps == 0 {
				eg.Go(func(node string) func() error {
					return func() error {

						taskStatus, taskOutput := action(node, e.executionContext)

						e.mu.Lock()
						defer e.mu.Unlock()

						taskStatusChan <- taskStatus
						e.executionContext[node] = taskOutput

						for n, deps := range graph {
							for _, dep := range deps {
								if dep == node {
									inDegrees[n]--
								}
							}
						}

						return nil
					}
				}(node))
			}
		}

		err := eg.Wait()
		if err != nil {
			e.logger.Info("find errors in parallel execution...")
			break
		}

		if len(e.executionContext) == len(graph) {
			break
		}
	}
}
