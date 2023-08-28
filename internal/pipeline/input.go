package pipeline

import (
	"fmt"
	"strings"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/lex"
	"github.com/octopipe/cloudx/internal/task"
	"github.com/octopipe/cloudx/internal/taskoutput"
	"k8s.io/apimachinery/pkg/types"
)

func (p *pipelineCtx) interpolateTaskInputsByExecutionContext(task commonv1alpha1.InfraTask, executionContext ExecutionContext) ([]commonv1alpha1.InfraTaskInput, error) {
	inputs := []commonv1alpha1.InfraTaskInput{}
	for _, i := range task.Inputs {
		tokens := lex.Tokenize(i.Value)
		data := map[string]string{}
		sensitive := false
		for _, t := range tokens {
			if t.Type == lex.TokenVariable {
				s := strings.Split(strings.Trim(t.Value, " "), ".")
				if len(s) != 3 {
					return nil, fmt.Errorf("malformed input variable %s with value %s", i.Key, i.Value)
				}

				value, isSensitive, err := p.getDataByOrigin(s[0], s[1], s[2], executionContext)
				if err != nil {
					return nil, err
				}

				if isSensitive {
					sensitive = isSensitive
				}

				data[t.Value] = strings.Trim(value, "\"")
			}
		}

		inputs = append(inputs, commonv1alpha1.InfraTaskInput{
			Key:       i.Key,
			Value:     lex.Interpolate(tokens, data),
			Sensitive: sensitive,
		})
	}

	return inputs, nil
}

func (p *pipelineCtx) getDataByOrigin(origin string, name string, attr string, executionContext ExecutionContext) (string, bool, error) {
	switch origin {
	case task.ThisInterpolationOrigin:
		p.logger.Info("interpolate this origin")
		execution, ok := executionContext[name]
		if !ok {
			return "", false, fmt.Errorf("not found task %s in execution context", name)
		}

		executionAttr, ok := execution[attr]
		if !ok {
			return "", false, fmt.Errorf("not found attr %s in finished task execution %s", attr, name)
		}

		return executionAttr.Value, executionAttr.Sensitive, nil

	case task.TaskOutputInterpolationOrigin:
		p.logger.Info("interpolate this task-output")
		taskOutput := commonv1alpha1.TaskOutput{}
		err := p.rpcClient.Call("TaskOutputRPCHandler.GetTaskOutput", taskoutput.RPCGetTaskOutputArgs{
			Ref: types.NamespacedName{Name: name, Namespace: "default"},
		}, &taskOutput)
		if err != nil {
			return "", false, err
		}

		for _, out := range taskOutput.Spec.Outputs {
			if out.Key == attr {
				return out.Value, out.Sensitive, nil
			}
		}

		return "", false, fmt.Errorf("not found attr in connection-interface %s", name)
	default:
		return "", false, fmt.Errorf("invalid origin type %s", origin)
	}
}
