package pluginmanager

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/plugin"
	"github.com/octopipe/cloudx/internal/provider"
	kubeyaml "k8s.io/apimachinery/pkg/util/yaml"
)

type Manager interface {
	Publish(pluginName string, filecontents map[string][]byte) error
	Execute(pluginRef string, input map[string]interface{}) ([]commonv1alpha1.StackSetPluginOutput, error)
}

type manager struct {
	terraformProvider provider.Provider
}

func NewPluginManager(terraformProvider provider.Provider) Manager {
	return manager{
		terraformProvider: terraformProvider,
	}
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

func (m manager) Execute(pluginRef string, input map[string]interface{}) ([]commonv1alpha1.StackSetPluginOutput, error) {
	fmt.Println("Pulling plugin")
	img, err := crane.Pull(pluginRef)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = crane.Export(img, &buf)
	if err != nil {
		return nil, err
	}

	tr := tar.NewReader(&buf)
	content := map[string]string{}

	rawPluginConfig := []byte{}

	workdir := fmt.Sprintf("/tmp/cloudx/executions/%s", strconv.Itoa(int(time.Now().Unix())))

	err = os.MkdirAll(workdir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		b, err := io.ReadAll(tr)
		if err != nil {
			return nil, err
		}

		content[hdr.Name] = string(b)
		f, err := os.Create(fmt.Sprintf("%s/%s", workdir, hdr.Name))
		if err != nil {
			return nil, err
		}

		_, err = f.Write(b)
		if err != nil {
			return nil, err
		}

		if hdr.Name == "plugin.yaml" {
			rawPluginConfig = b
		}
	}

	pluginConfig := plugin.Plugin{}
	decoder := kubeyaml.NewYAMLOrJSONDecoder(bytes.NewReader(rawPluginConfig), 4096)
	err = decoder.Decode(&pluginConfig)
	if err != nil {
		return nil, err
	}

	parsedInput := map[string]interface{}{}
	for _, i := range pluginConfig.Spec.Inputs {
		value, ok := input[i.Name]
		if !ok && i.Required {
			return nil, fmt.Errorf("required field: %s", i.Name)
		}

		if !ok {
			parsedInput[i.Name] = i.Default
			continue
		}

		parsedInput[i.Name] = value
	}

	fmt.Printf("Executing plugin on %s\n", workdir)
	return m.terraformProvider.Apply(workdir, parsedInput)
}
