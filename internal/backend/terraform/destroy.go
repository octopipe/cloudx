package terraform

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-exec/tfexec"
	"go.uber.org/zap"
)

func (t terraformBackend) Destroy(input TerraformDestroyInput) error {
	workdirPath, err := t.dowloadSource(input.Source)
	if err != nil {
		return err
	}

	terraformPath, err := t.install(input.Version)
	if err != nil {
		return err
	}

	execVarsFilePath := filepath.Join(workdirPath, "exec.tfvars")
	f, err := os.Create(execVarsFilePath)
	if err != nil {
		return err
	}

	t.logger.Info("creating terraform vars file", zap.String("workdir", workdirPath))
	for _, i := range input.TaskInputs {
		f.WriteString(fmt.Sprintf("%s = \"%s\"\n", i.Key, i.Value))
	}

	tf, err := tfexec.NewTerraform(workdirPath, terraformPath)
	if err != nil {
		return err
	}

	if input.PreviousLockDeps != "" {
		err = persistDependenciesLock(input.PreviousLockDeps, workdirPath)
		if err != nil {
			return err
		}
	}

	t.logger.Info("executing terraform init", zap.String("workdir", workdirPath))
	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {

		return err
	}

	if input.PreviousState != "" {
		err := persistPreviousState(input.PreviousState, workdirPath)
		if err != nil {
			return err
		}
	}

	t.logger.Info("executing terraform destroy", zap.String("workdir", workdirPath))
	err = tf.Destroy(context.Background(), tfexec.VarFile(execVarsFilePath))

	return err
}
