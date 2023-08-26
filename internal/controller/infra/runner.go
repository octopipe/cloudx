package infra

import (
	"fmt"
	"os"
	"time"

	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

type Runner struct {
	Pod     *v1.Pod
	Service *v1.Service
}

func (c *controller) NewRunner(action string, infra commonv1alpha1.Infra, providerConfig commonv1alpha1.ProviderConfig, varsCreds []v1.EnvVar) (Runner, error) {
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
	serviceAccount := "cloudx-controller"

	if infra.Spec.RunnerConfig.ServiceAccount != "" {
		serviceAccount = infra.Spec.RunnerConfig.ServiceAccount
	}

	defaultVars := []v1.EnvVar{
		{
			Name:  "TF_VERSION",
			Value: "latest",
		},
		{
			Name:  "RPC_SERVER_ADDRESS",
			Value: os.Getenv("RPC_SERVER_ADDRESS"),
		},
	}

	infraRef := types.NamespacedName{
		Name:      infra.GetName(),
		Namespace: infra.GetNamespace(),
	}

	defaultVars = append(defaultVars, varsCreds...)

	args := []string{"/usr/local/bin/runner", action, infraRef.String()}

	newRunnerObject := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-runner-%d", infra.GetName(), time.Now().Unix()),
			Namespace: "default",
			Labels: map[string]string{
				"commons.cloudx.io/infra-name":      infra.GetName(),
				"commons.cloudx.io/infra-namespace": infra.GetNamespace(),
				"app.kubernetes.io/managed-by":      "cloudx",
			},
		},
		Spec: v1.PodSpec{
			RestartPolicy: v1.RestartPolicyNever,
			Containers: []v1.Container{
				{
					Name:            "runner",
					Image:           "mayconjrpacheco/cloudx:latest",
					Args:            args,
					ImagePullPolicy: v1.PullAlways,
					SecurityContext: securityContext,
					Env:             defaultVars,
					VolumeMounts:    podVolumeMounts,
				},
			},
			Volumes: podVolumes,
		},
	}

	if os.Getenv("ENV") != "local" {
		newRunnerObject.Namespace = "cloudx-system"
		newRunnerObject.Spec.ServiceAccountName = serviceAccount
	}

	return Runner{
		Pod: newRunnerObject,
	}, nil
}
