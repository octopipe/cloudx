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
	vFalse := false
	vTrue := true
	vUser := int64(65532)
	securityContext := &v1.SecurityContext{
		Capabilities: &v1.Capabilities{
			Drop: []v1.Capability{"ALL"},
		},
		AllowPrivilegeEscalation: &vFalse,
		RunAsNonRoot:             &vTrue,
		RunAsUser:                &vUser,
		SeccompProfile: &v1.SeccompProfile{
			Type: v1.SeccompProfileTypeRuntimeDefault,
		},
		ReadOnlyRootFilesystem: &vTrue,
	}

	podVolumes := []v1.Volume{
		{
			Name: "temp",
			VolumeSource: v1.VolumeSource{
				EmptyDir: &v1.EmptyDirVolumeSource{},
			},
		},
		{
			Name: "home",
			VolumeSource: v1.VolumeSource{
				EmptyDir: &v1.EmptyDirVolumeSource{},
			},
		},
	}

	podVolumeMounts := []v1.VolumeMount{
		{
			Name:      "temp",
			MountPath: "/tmp",
		},
		{
			Name:      "home",
			MountPath: "/home/runner",
		},
	}
	serviceAccount := "controller-sa"

	if sharedInfra.Spec.RunnerConfig.ServiceAccount != "" {
		serviceAccount = sharedInfra.Spec.RunnerConfig.ServiceAccount
	}

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
			ServiceAccountName: serviceAccount,
			RestartPolicy:      v1.RestartPolicyNever,
			Containers: []v1.Container{
				{
					Name:            "runner",
					Image:           "mayconjrpacheco/cloudx-runner:latest",
					Args:            []string{sharedInfraRef, executionId},
					ImagePullPolicy: v1.PullAlways,
					SecurityContext: securityContext,
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
					VolumeMounts: podVolumeMounts,
				},
			},
			Volumes: podVolumes,
		},
	}

	return Runner{
		Object: newRunnerObject,
	}
}