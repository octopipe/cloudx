package terraform

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"go.uber.org/zap"
)

type TerraformProvider interface {
	Apply(workdirPath string, input map[string]interface{}) ([]commonv1alpha1.SharedInfraPluginOutput, string, error)
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

func (p terraformProvider) Apply(workdirPath string, input map[string]interface{}) ([]commonv1alpha1.SharedInfraPluginOutput, string, error) {
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

	outputs := []commonv1alpha1.SharedInfraPluginOutput{}

	for key, value := range out {
		outputs = append(outputs, commonv1alpha1.SharedInfraPluginOutput{
			Key:   key,
			Value: string(value.Value),
		})
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
