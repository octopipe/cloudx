package io

import (
	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
)

type ProviderInputMetadata struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

type ProviderInput map[string]ProviderInputMetadata

type ProviderOutputMetadata struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	Sensitive bool   `json:"sensitive"`
	Type      string `json:"type"`
}

type ProviderOutput map[string]ProviderOutputMetadata

func ToProviderInput(pluginInputs []commonv1alpha1.SharedInfraPluginInput) ProviderInput {
	i := ProviderInput{}

	for _, p := range pluginInputs {
		i[p.Key] = ProviderInputMetadata{
			Value: p.Value,
			Type:  "text",
		}
	}

	return i
}
