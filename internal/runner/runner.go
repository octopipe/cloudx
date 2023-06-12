package runner

import (
	"fmt"
	"time"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Runner struct {
	Object *v1.Pod
}

func NewRunner(executionId string, sharedInfraRef string, sharedInfra commonv1alpha1.SharedInfra) Runner {
	newRunnerObject := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-runner-%d", sharedInfra.GetName(), time.Now().Unix()),
			Namespace: "cloudx-system",
			Labels: map[string]string{
				"commons.cloudx.io/sharedinfra-name":      sharedInfra.GetName(),
				"commons.cloudx.io/sharedinfra-namespace": sharedInfra.GetNamespace(),
				"commons.cloudx.io/execution-id":          executionId,
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
					Args:            []string{sharedInfraRef, executionId},
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

	return Runner{
		Object: newRunnerObject,
	}
}
