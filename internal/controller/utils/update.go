package utils

import (
	"context"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func UpdateInfraStatus(client client.Client, infra commonv1alpha1.Infra) error {
	return retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		return client.Status().Update(context.TODO(), &infra)
	})
}
