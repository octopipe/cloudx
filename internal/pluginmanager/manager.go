package pluginmanager

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/octopipe/cloudx/internal/execution"
)

type Manager interface {
	Publish(pluginName string, filecontents map[string][]byte) error
	Execute(pluginName string) error
}

type manager struct {
}

func NewPluginManager() Manager {
	return manager{}
}

func (m manager) Publish(pluginName string, filecontents map[string][]byte) error {
	tag, err := name.NewTag("mayconjrpacheco/ec2-plugin:latest")
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

func (m manager) Execute(pluginName string) error {
	fmt.Println("Pulling plugin")
	img, err := crane.Pull("mayconjrpacheco/ec2-plugin:latest")
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = crane.Export(img, &buf)
	if err != nil {
		return err
	}

	tr := tar.NewReader(&buf)
	content := map[string]string{}

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		b, err := io.ReadAll(tr)
		if err != nil {
			return err
		}

		content[hdr.Name] = string(b)
		f, err := os.Create(fmt.Sprintf("./tmp/%s", hdr.Name))
		if err != nil {
			return err
		}

		_, err = f.Write(b)
		if err != nil {
			return err
		}
	}

	fmt.Println("Executing plugin")
	out, err := execution.NewTerraformExecution("./tmp")
	if err != nil {
		return err
	}

	fmt.Println(out)

	return nil
}
