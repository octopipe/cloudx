package terraform

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"os"
	"path/filepath"

	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/uuid"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/plugin"
	"go.uber.org/zap"
	kubeyaml "k8s.io/apimachinery/pkg/util/yaml"
)

type TerraformProvider interface {
	Apply(pluginRef string, inputs []commonv1alpha1.SharedInfraPluginInput, previousState string, previousLockDeps string) (map[string]any, string, string, error)
	Destroy(pluginRef string, inputs []commonv1alpha1.SharedInfraPluginInput, previousState string, previousLockDeps string) error
}

type terraformProvider struct {
	logger   *zap.Logger
	execPath string
}

func NewTerraformProvider(logger *zap.Logger, terraformVersion string) (TerraformProvider, error) {
	installDirPath := "/tmp/cloudx/terraform-versions"
	err := os.MkdirAll(installDirPath, os.ModePerm)
	if err != nil {
		return nil, err
	}

	if terraformVersion != "" {
		installer := &releases.ExactVersion{
			Product:    product.Terraform,
			Version:    version.Must(version.NewVersion(terraformVersion)),
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

	installer := &releases.LatestVersion{
		Product:    product.Terraform,
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

func (p terraformProvider) prepareExecution(pluginRef string, input []commonv1alpha1.SharedInfraPluginInput) (string, plugin.Plugin, error) {
	p.logger.Info("pulling plugin image", zap.String("image", pluginRef))
	img, err := crane.Pull(pluginRef)
	if err != nil {
		return "", plugin.Plugin{}, err
	}

	var buf bytes.Buffer
	err = crane.Export(img, &buf)
	if err != nil {
		return "", plugin.Plugin{}, err
	}

	tr := tar.NewReader(&buf)
	content := map[string]string{}
	rawPluginConfig := []byte{}
	workdir := fmt.Sprintf("/tmp/cloudx/executions/%s", uuid.New().String())
	err = os.MkdirAll(workdir, os.ModePerm)
	if err != nil {
		return "", plugin.Plugin{}, err
	}

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", plugin.Plugin{}, err
		}

		b, err := io.ReadAll(tr)
		if err != nil {
			return "", plugin.Plugin{}, err
		}

		content[hdr.Name] = string(b)
		f, err := os.Create(fmt.Sprintf("%s/%s", workdir, hdr.Name))
		if err != nil {
			return "", plugin.Plugin{}, err
		}

		_, err = f.Write(b)
		if err != nil {
			return "", plugin.Plugin{}, err
		}

		if hdr.Name == "plugin.yaml" {
			rawPluginConfig = b
		}
	}

	pluginConfig := plugin.Plugin{}
	decoder := kubeyaml.NewYAMLOrJSONDecoder(bytes.NewReader(rawPluginConfig), 4096)
	err = decoder.Decode(&pluginConfig)
	if err != nil {
		return "", plugin.Plugin{}, err
	}

	return workdir, pluginConfig, nil
}

func (p terraformProvider) Apply(pluginRef string, inputs []commonv1alpha1.SharedInfraPluginInput, previousState string, previousLockDeps string) (map[string]any, string, string, error) {
	workdirPath, pluginConfig, err := p.prepareExecution(pluginRef, inputs)
	if err != nil {
		return nil, "", "", err
	}

	execVarsFilePath := filepath.Join(workdirPath, "exec.tfvars")
	f, err := os.Create(execVarsFilePath)
	if err != nil {
		return nil, "", "", err
	}

	p.logger.Info("creating terraform vars file", zap.String("workdir", workdirPath))
	for _, i := range pluginConfig.Spec.Inputs {

		var value commonv1alpha1.SharedInfraPluginInput
		for _, s := range inputs {
			if i.Name == s.Key {
				value = s
				break
			}
		}

		if value.Value == "" && i.Required {
			return nil, "", "", fmt.Errorf("required field: %s", i.Name)
		}

		if value.Value == "" {
			f.WriteString(fmt.Sprintf("%s = \"%s\"\n", i.Name, i.Default))
			continue
		}

		f.WriteString(fmt.Sprintf("%s = \"%s\"\n", i.Name, value.Value))
	}

	tf, err := tfexec.NewTerraform(workdirPath, p.execPath)
	if err != nil {
		return nil, "", "", err
	}

	// if previousLockDeps != "" {
	// 	fmt.Println("LOCK DEPS", previousLockDeps)
	// 	rawPreviousLockDeps, err := base64.StdEncoding.DecodeString(strings.Trim(previousLockDeps, "\""))
	// 	if err != nil {
	// 		return nil, "", "", err
	// 	}

	// 	fmt.Println("RAW", string(rawPreviousLockDeps))

	// 	p.logger.Info("using lock file to increase performance")
	// 	previousLockDepsFilePath := filepath.Join(workdirPath, ".terraform.lock.hcl")
	// 	previousLockDepsFile, err := os.Create(previousLockDepsFilePath)
	// 	if err != nil {
	// 		return nil, "", "", err
	// 	}

	// 	previousLockDepsFile.Write(rawPreviousLockDeps)
	// }

	p.logger.Info("executing terraform init", zap.String("workdir", workdirPath))
	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		return nil, "", "", err
	}

	if previousState != "" {
		rawPreviousState, err := base64.StdEncoding.DecodeString(strings.Trim(previousState, "\""))
		if err != nil {
			return nil, "", "", err
		}

		previousStateFilePath := filepath.Join(workdirPath, "terraform.tfstate")
		previousStateFile, err := os.Create(previousStateFilePath)
		if err != nil {
			return nil, "", "", err
		}

		previousStateFile.Write(rawPreviousState)
	}

	p.logger.Info("executing terraform plan", zap.String("workdir", workdirPath))
	hasModifications, err := tf.Plan(context.Background(), tfexec.VarFile(execVarsFilePath))
	if err != nil {
		return nil, "", "", err
	}

	if hasModifications {
		p.logger.Info("executing terraform apply", zap.String("workdir", workdirPath))
		err = tf.Apply(context.Background(), tfexec.VarFile(execVarsFilePath))
		if err != nil {
			return nil, "", "", err
		}
	}

	out, err := tf.Output(context.Background())
	if err != nil {
		return nil, "", "", err
	}

	p.logger.Info("get terraform state file", zap.String("workdir", workdirPath))
	stateFilePath := fmt.Sprintf("%s/terraform.tfstate", workdirPath)
	stateFile, err := os.ReadFile(stateFilePath)
	if err != nil {
		return nil, "", "", err
	}

	p.logger.Info("get terraform lock deps", zap.String("workdir", workdirPath))
	lockDepsFilePath := fmt.Sprintf("%s/.terraform.lock.hcl", workdirPath)
	lockDepsFile, err := os.ReadFile(lockDepsFilePath)
	if err != nil {
		return nil, "", "", err
	}

	escapedState, err := json.Marshal(stateFile)
	if err != nil {
		return nil, "", "", err
	}

	escapedLockFile, err := json.Marshal(lockDepsFile)
	if err != nil {
		return nil, "", "", err
	}

	outputs := map[string]any{}

	for key, res := range out {
		outputs[key] = string(res.Value)
	}

	return outputs, string(escapedState), string(escapedLockFile), nil
}

func (p terraformProvider) Destroy(pluginRef string, inputs []commonv1alpha1.SharedInfraPluginInput, previousState string, previousLockDeps string) error {
	workdirPath, pluginConfig, err := p.prepareExecution(pluginRef, inputs)
	if err != nil {
		return err
	}

	execVarsFilePath := filepath.Join(workdirPath, "exec.tfvars")
	f, err := os.Create(execVarsFilePath)
	if err != nil {
		return err
	}

	p.logger.Info("creating terraform vars file", zap.String("workdir", workdirPath))
	for _, i := range pluginConfig.Spec.Inputs {

		var value commonv1alpha1.SharedInfraPluginInput
		for _, s := range inputs {
			if i.Name == s.Key {
				value = s
				break
			}
		}

		if value.Value == "" && i.Required {
			return fmt.Errorf("required field: %s", i.Name)
		}

		if value.Value == "" {
			f.WriteString(fmt.Sprintf("%s = \"%s\"\n", i.Name, i.Default))
			continue
		}

		f.WriteString(fmt.Sprintf("%s = \"%s\"\n", i.Name, value.Value))
	}

	tf, err := tfexec.NewTerraform(workdirPath, p.execPath)
	if err != nil {
		return err
	}

	// if previousLockDeps != "" {
	// 	p.logger.Info("using lock file to increase performance")
	// 	previousLockDepsFilePath := filepath.Join(workdirPath, ".terraform.lock.hcl")
	// 	previousLockDepsFile, err := os.Create(previousLockDepsFilePath)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	var unescapedJSON string
	// 	err = json.Unmarshal([]byte(previousLockDeps), &unescapedJSON)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	previousLockDepsFile.Write([]byte(unescapedJSON))
	// }

	p.logger.Info("executing terraform init", zap.String("workdir", workdirPath))
	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {

		return err
	}

	if previousState != "" {
		rawPreviousState, err := base64.StdEncoding.DecodeString(strings.Trim(previousState, "\""))
		if err != nil {
			return err
		}

		previousStateFilePath := filepath.Join(workdirPath, "terraform.tfstate")
		previousStateFile, err := os.Create(previousStateFilePath)
		if err != nil {
			return err
		}

		previousStateFile.Write(rawPreviousState)
	}

	p.logger.Info("executing terraform destroy", zap.String("workdir", workdirPath))
	err = tf.Destroy(context.Background(), tfexec.VarFile(execVarsFilePath))

	return err
}
