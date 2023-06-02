package terraform

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"

	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/octopipe/cloudx/internal/plugin"
	providerIO "github.com/octopipe/cloudx/internal/provider/io"
	"go.uber.org/zap"
	kubeyaml "k8s.io/apimachinery/pkg/util/yaml"
)

type TerraformProvider interface {
	Apply(pluginRef string, input providerIO.ProviderInput) (providerIO.ProviderOutput, string, error)
	Destroy(workdirPath string, state string) error
}

type terraformProvider struct {
	logger   *zap.Logger
	execPath string
}

func NewTerraformProvider(logger *zap.Logger) (TerraformProvider, error) {
	installDirPath := "/tmp/terraform-ins"
	err := os.MkdirAll(installDirPath, os.ModePerm)
	if err != nil {
		return nil, err
	}

	installer := &releases.ExactVersion{
		Product:    product.Terraform,
		Version:    version.Must(version.NewVersion("1.0.6")),
		InstallDir: installDirPath,
	}

	execPath, err := installer.Install(context.Background())
	if err != nil {
		return nil, err
	}

	return terraformProvider{
		logger:   logger,
		execPath: execPath,
	}, nil
}

func (p terraformProvider) prepareExecution(pluginRef string, input providerIO.ProviderInput) (string, map[string]interface{}, error) {
	p.logger.Info("pulling plugin image", zap.String("image", pluginRef))
	img, err := crane.Pull(pluginRef)
	if err != nil {
		return "", nil, err
	}

	var buf bytes.Buffer
	err = crane.Export(img, &buf)
	if err != nil {
		return "", nil, err
	}

	tr := tar.NewReader(&buf)
	content := map[string]string{}

	rawPluginConfig := []byte{}

	workdir := fmt.Sprintf("/tmp/cloudx/executions/%s", strconv.Itoa(int(time.Now().UnixNano())))

	err = os.MkdirAll(workdir, os.ModePerm)
	if err != nil {
		return "", nil, err
	}

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", nil, err
		}

		b, err := io.ReadAll(tr)
		if err != nil {
			return "", nil, err
		}

		content[hdr.Name] = string(b)
		f, err := os.Create(fmt.Sprintf("%s/%s", workdir, hdr.Name))
		if err != nil {
			return "", nil, err
		}

		_, err = f.Write(b)
		if err != nil {
			return "", nil, err
		}

		if hdr.Name == "plugin.yaml" {
			rawPluginConfig = b
		}
	}

	pluginConfig := plugin.Plugin{}
	decoder := kubeyaml.NewYAMLOrJSONDecoder(bytes.NewReader(rawPluginConfig), 4096)
	err = decoder.Decode(&pluginConfig)
	if err != nil {
		return "", nil, err
	}

	parsedInput := map[string]interface{}{}
	for _, i := range pluginConfig.Spec.Inputs {
		value, ok := input[i.Name]
		if !ok && i.Required {
			return "", nil, fmt.Errorf("required field: %s", i.Name)
		}

		if !ok {
			parsedInput[i.Name] = i.Default
			continue
		}

		parsedInput[i.Name] = value
	}

	return workdir, parsedInput, nil
}

func (p terraformProvider) Apply(pluginRef string, executionInput providerIO.ProviderInput) (providerIO.ProviderOutput, string, error) {
	workdirPath, input, err := p.prepareExecution(pluginRef, executionInput)
	if err != nil {
		return nil, "", err
	}

	tf, err := tfexec.NewTerraform(workdirPath, p.execPath)
	if err != nil {
		return nil, "", err
	}

	p.logger.Info("executing terraform init", zap.String("workdir", workdirPath))
	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {

		return nil, "", err
	}

	execVarsFilePath := filepath.Join(workdirPath, "exec.tfvars")
	f, err := os.Create(execVarsFilePath)
	if err != nil {
		return nil, "", err
	}

	p.logger.Info("creating terraform vars file", zap.String("workdir", workdirPath))
	for key, val := range input {
		f.WriteString(fmt.Sprintf("%s = \"%s\"\n", key, val))
	}

	p.logger.Info("executing terraform apply", zap.String("workdir", workdirPath))
	err = tf.Apply(context.Background(), tfexec.VarFile(execVarsFilePath))
	if err != nil {
		return nil, "", err
	}

	out, err := tf.Output(context.Background())
	if err != nil {
		return nil, "", err
	}

	outputs := providerIO.ProviderOutput{}

	for key, res := range out {
		outputs[key] = providerIO.ProviderOutputMetadata{
			Value:     string(res.Value),
			Sensitive: res.Sensitive,
		}
	}

	p.logger.Info("get terraform state file", zap.String("workdir", workdirPath))
	stateFilePath := fmt.Sprintf("%s/terraform.tfstate", workdirPath)
	stateFile, err := os.ReadFile(stateFilePath)
	if err != nil {
		return nil, "", err
	}

	return outputs, string(stateFile), nil
}

func (p terraformProvider) Destroy(workdirPath string, state string) error {
	tf, err := tfexec.NewTerraform(workdirPath, p.execPath)
	if err != nil {
		return err
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		return err
	}

	execVarsFilePath := filepath.Join(workdirPath, "exec.tfvars")
	_, err = os.Create(execVarsFilePath)
	if err != nil {
		return err
	}

	err = tf.Destroy(context.Background(), tfexec.StateOut(""))
	if err != nil {
		return err
	}

	return nil
}
