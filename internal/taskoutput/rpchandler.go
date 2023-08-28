package taskoutput

import (
	"context"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type TaskOutputRPCHandler struct {
	logger               *zap.Logger
	taskOutputRepository Repository
	k8sClient            client.Client
}

func NewTaskOutputRPCHandler(logger *zap.Logger, k8sClient client.Client, taskOutputRepository Repository) *TaskOutputRPCHandler {
	return &TaskOutputRPCHandler{
		logger:               logger,
		taskOutputRepository: taskOutputRepository,
		k8sClient:            k8sClient,
	}
}

type RPCGetTaskOutputArgs struct {
	Ref types.NamespacedName
}

func (h *TaskOutputRPCHandler) GetTaskOutput(args *RPCGetTaskOutputArgs, reply *commonv1alpha1.TaskOutput) error {
	currentTaskOutput, err := h.taskOutputRepository.Get(context.Background(), args.Ref.Name, args.Ref.Namespace)
	*reply = currentTaskOutput
	return err
}

type RPCCreateTaskOutputItem struct {
	commonv1alpha1.InfraTaskOutputItem
	Value string
}

type RPCCreateTaskOutputArgs struct {
	Name      string
	Namespace string
	Items     []RPCCreateTaskOutputItem
	TaskName  string
	InfraRef  commonv1alpha1.Ref
}

func (h TaskOutputRPCHandler) isSecretTaskOutput(items []RPCCreateTaskOutputItem) bool {
	hasSensitiveData := false

	for _, i := range items {
		if i.Sensitive {
			hasSensitiveData = true
			break
		}
	}

	return hasSensitiveData

}

func (h *TaskOutputRPCHandler) applySecret(args *RPCCreateTaskOutputArgs) (types.NamespacedName, error) {
	newSecret := v1.Secret{}
	newSecret.SetName(args.Name)
	newSecret.SetNamespace(args.Namespace)

	for _, i := range args.Items {
		if i.Sensitive {
			newSecret.Data[i.Key] = []byte(i.Value)
		}
	}

	err := h.k8sClient.Create(context.Background(), &newSecret)
	if err != nil && errors.IsAlreadyExists(err) {
		err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
			currentSecret := v1.Secret{}

			err = h.k8sClient.Get(context.Background(), types.NamespacedName{
				Name:      args.Name,
				Namespace: args.Namespace,
			}, &currentSecret)

			currentSecret.Data = newSecret.Data

			return h.k8sClient.Update(context.Background(), &currentSecret)
		})

		return types.NamespacedName{
			Name:      newSecret.Name,
			Namespace: newSecret.Namespace,
		}, err
	}

	return types.NamespacedName{
		Name:      newSecret.Name,
		Namespace: newSecret.Namespace,
	}, err
}

func (h *TaskOutputRPCHandler) ApplyTaskOuput(args *RPCCreateTaskOutputArgs, reply *int) error {
	h.logger.Info("received call", zap.String("method", "TaskOutputRPCHandler.ApplyTaskOuput"), zap.String("taskoutput", args.Name))
	newTaskOutput := commonv1alpha1.TaskOutput{}
	if h.isSecretTaskOutput(args.Items) {
		secretRef, err := h.applySecret(args)
		if err != nil {
			return err
		}

		newTaskOutput.Spec.Secret = commonv1alpha1.Ref{
			Name:      secretRef.Name,
			Namespace: secretRef.Namespace,
		}
	}

	newTaskOutput.SetName(args.Name)
	newTaskOutput.SetNamespace(args.Namespace)
	for _, item := range args.Items {
		newOutput := commonv1alpha1.TaskOutputSpecItem{
			Key:       item.Key,
			Sensitive: item.Sensitive,
		}

		if !item.Sensitive {
			newOutput.Value = item.Value
		}

		newTaskOutput.Spec.Outputs = append(newTaskOutput.Spec.Outputs, newOutput)
	}

	newTaskOutput.Spec.Infra = args.InfraRef
	newTaskOutput.Spec.TaskName = args.TaskName

	_, err := h.taskOutputRepository.Apply(context.Background(), newTaskOutput)
	if err != nil {
		return err
	}

	_, err = h.taskOutputRepository.Get(context.Background(), args.Name, args.Namespace)
	// *reply = currentTaskOutput
	return err
}

type RPCDeleteTaskOutputArgs struct {
	Name      string
	Namespace string
}

func (h *TaskOutputRPCHandler) DeleteTaskOutput(args *RPCCreateTaskOutputArgs, reply *int) error {
	newTaskOutput := commonv1alpha1.TaskOutput{}
	if h.isSecretTaskOutput(args.Items) {
		secretRef, err := h.applySecret(args)
		if err != nil {
			return err
		}

		newTaskOutput.Spec.Secret = commonv1alpha1.Ref{
			Name:      secretRef.Name,
			Namespace: secretRef.Namespace,
		}
	}
	currentSecret := v1.Secret{}
	currentSecret.SetName(args.Name)
	currentSecret.SetNamespace(args.Namespace)
	err := h.k8sClient.Delete(context.Background(), &currentSecret)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	err = h.taskOutputRepository.Delete(context.Background(), args.Name, args.Namespace)
	return err
}
