package sharedinfra

import (
	"context"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func updateSharedInfraStatus(client client.Client, sharedInfra *commonv1alpha1.SharedInfra) error {
	return retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		return client.Status().Update(context.TODO(), sharedInfra)
	})
}

func updateExecutionStatus(client client.Client, execution *commonv1alpha1.Execution) error {
	return retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		return client.Status().Update(context.TODO(), execution)
	})
}