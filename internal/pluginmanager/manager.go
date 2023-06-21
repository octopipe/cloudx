package pluginmanager

import (
	"fmt"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/octopipe/cloudx/internal/terraform"
	"go.uber.org/zap"
)

type Manager interface {
	Publish(pluginName string, filecontents map[string][]byte) error
}

type manager struct {
	logger            *zap.Logger
	terraformProvider terraform.TerraformProvider
}

func NewPluginManager(logger *zap.Logger, terraformProvider terraform.TerraformProvider) Manager {
	return manager{
		logger:            logger,
		terraformProvider: terraformProvider,
	}
}

func (m manager) Publish(pluginName string, filecontents map[string][]byte) error {
	tag, err := name.NewTag(fmt.Sprintf("mayconjrpacheco/plugin:%s", pluginName))
	if err != nil {
		return err
	}

	newImage, err := crane.Image(filecontents)
	if err != nil {
		return err
	}

	err = crane.Push(newImage, tag.String())
	if err != nil {
		return err
	}

	return nil
}
