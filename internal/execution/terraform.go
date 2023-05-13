package execution

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
)

func NewTerraformExecution(workdirPath string, input map[string]interface{}) ([]commonv1alpha1.StackPluginOutput, error) {
	installer := &releases.ExactVersion{
		Product:    product.Terraform,
		Version:    version.Must(version.NewVersion("1.0.6")),
		InstallDir: "/tmp/terraform-ins",
	}

	execPath, err := installer.Install(context.Background())
	if err != nil {
		return nil, err
	}

	tf, err := tfexec.NewTerraform(workdirPath, execPath)
	if err != nil {
		return nil, err
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		return nil, err
	}

	execVarsFile := filepath.Join(workdirPath, "exec.tfvars")
	f, err := os.Create(execVarsFile)
	if err != nil {
		return nil, err
	}

	for key, val := range input {
		f.WriteString(fmt.Sprintf("%s = \"%s\"\n", key, val))
	}

	err = tf.Apply(context.Background(), tfexec.VarFile(execVarsFile))
	if err != nil {
		return nil, err
	}

	out, err := tf.Output(context.Background())
	if err != nil {
		return nil, err
	}

	outputs := []commonv1alpha1.StackPluginOutput{}

	for key, value := range out {
		outputs = append(outputs, commonv1alpha1.StackPluginOutput{
			Key:   key,
			Value: string(value.Value),
		})
	}

	return outputs, nil
}
