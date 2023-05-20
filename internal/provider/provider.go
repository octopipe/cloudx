package provider

import commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"

type Provider interface {
	Apply(workdirPath string, input map[string]interface{}) ([]commonv1alpha1.StackSetPluginOutput, string, error)
}
