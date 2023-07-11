package pipeline

import (
	"fmt"
	"strings"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/lex"
	"github.com/octopipe/cloudx/internal/task"
)

func (p pipelineCtx) validateDependencies(infra commonv1alpha1.Infra) error {
	graph := map[string][]string{}
	for _, task := range infra.Spec.Tasks {
		graph[task.Name] = task.Depends
	}

	for _, p := range infra.Spec.Tasks {
		for _, dep := range p.Depends {
			if _, ok := graph[dep]; !ok {
				return fmt.Errorf("not found the dependency %s specified in task %s", dep, p.Name)
			}
		}
	}

	return nil
}

func (p pipelineCtx) validateInputInterpolations(infra commonv1alpha1.Infra) error {
	graph := map[string][]string{}
	for _, task := range infra.Spec.Tasks {
		graph[task.Name] = task.Depends
	}

	for _, p := range infra.Spec.Tasks {
		for _, i := range p.Inputs {
			tokens := lex.Tokenize(i.Value)

			for _, t := range tokens {
				if t.Type == lex.TokenVariable {
					s := strings.Split(strings.Trim(t.Value, " "), ".")
					if len(s) != 3 {
						return fmt.Errorf("malformed input variable %s with value %s", i.Key, i.Value)
					}

					origin, name := s[0], s[1]
					if origin != task.ThisInterpolationOrigin && origin != task.TaskOutputInterpolationOrigin {
						return fmt.Errorf("invalid origin: %s for input %s interpolation with value %s", origin, i.Key, i.Value)
					}

					if origin == task.ThisInterpolationOrigin {
						if _, ok := graph[name]; !ok {
							return fmt.Errorf("invalid name: %s in origin this for input %s interpolation with value %s", name, i.Key, i.Value)
						}
					}
				}
			}
		}
	}

	return nil
}

func (p pipelineCtx) validateInfra(infra commonv1alpha1.Infra) error {
	err := p.validateDependencies(infra)
	if err != nil {
		return err
	}

	err = p.validateInputInterpolations(infra)
	if err != nil {
		return err
	}

	return err
}
