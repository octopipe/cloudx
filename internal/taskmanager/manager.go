package taskmanager

import (
	"fmt"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	"go.uber.org/zap"
)

type Manager interface {
	Publish(taskName string, filecontents map[string][]byte) error
}

type manager struct {
	logger *zap.Logger
}

func NewTaskManager(logger *zap.Logger) Manager {
	return manager{
		logger: logger,
	}
}

func (m manager) Publish(taskName string, filecontents map[string][]byte) error {
	tag, err := name.NewTag(fmt.Sprintf("mayconjrpacheco/task:%s", taskName))
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
