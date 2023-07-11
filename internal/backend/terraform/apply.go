package terraform

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-exec/tfexec"
	"go.uber.org/zap"
)

func persistDependenciesLock(previousDependenciesLock string, workdirPath string) error {
	rawPreviousLockDeps, err := base64.StdEncoding.DecodeString(strings.Trim(previousDependenciesLock, "\""))
	if err != nil {
		return err
	}

	previousLockDepsFilePath := filepath.Join(workdirPath, ".terraform.lock.hcl")
	previousLockDepsFile, err := os.Create(previousLockDepsFilePath)
	if err != nil {
		return err
	}

	previousLockDepsFile.Write(rawPreviousLockDeps)
	return nil
}

func persistPreviousState(previousState string, workdirPath string) error {
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
	return nil
}

func (t terraformBackend) Apply(input TerraformApplyInput) (TerraformApplyResult, error) {
	workdirPath, err := t.dowloadSource(input.Source)
	if err != nil {
		return TerraformApplyResult{}, err
	}

	terraformPath, err := t.install(input.Version)
	if err != nil {
		return TerraformApplyResult{}, err
	}

	execVarsFilePath := filepath.Join(workdirPath, "exec.tfvars")
	f, err := os.Create(execVarsFilePath)
	if err != nil {
		return TerraformApplyResult{}, err
	}

	t.logger.Info("creating terraform vars file", zap.String("workdir", workdirPath))
	for _, i := range input.TaskInputs {
		f.WriteString(fmt.Sprintf("%s = \"%s\"\n", i.Key, i.Value))
	}

	tf, err := tfexec.NewTerraform(workdirPath, terraformPath)
	if err != nil {
		return TerraformApplyResult{}, err
	}

	if input.PreviousLockDeps != "" {
		err = persistDependenciesLock(input.PreviousLockDeps, workdirPath)
		if err != nil {
			return TerraformApplyResult{}, err
		}
	}

	t.logger.Info("executing terraform init", zap.String("workdir", workdirPath))
	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		return TerraformApplyResult{}, err
	}

	if input.PreviousState != "" {
		err := persistPreviousState(input.PreviousState, workdirPath)
		if err != nil {
			return TerraformApplyResult{}, err
		}
	}

	t.logger.Info("executing terraform plan", zap.String("workdir", workdirPath))
	hasModifications, err := tf.Plan(context.Background(), tfexec.VarFile(execVarsFilePath))
	if err != nil {
		return TerraformApplyResult{}, err
	}

	if hasModifications {
		t.logger.Info("executing terraform apply", zap.String("workdir", workdirPath))
		err = tf.Apply(context.Background(), tfexec.VarFile(execVarsFilePath))
		if err != nil {
			return TerraformApplyResult{}, err
		}
	}

	out, err := tf.Output(context.Background())
	if err != nil {
		return TerraformApplyResult{}, err
	}

	t.logger.Info("get terraform state file", zap.String("workdir", workdirPath))
	stateFilePath := fmt.Sprintf("%s/terraform.tfstate", workdirPath)
	stateFile, err := os.ReadFile(stateFilePath)
	if err != nil {
		return TerraformApplyResult{}, err
	}

	t.logger.Info("get terraform lock deps", zap.String("workdir", workdirPath))
	lockDepsFilePath := fmt.Sprintf("%s/.terraform.lock.hcl", workdirPath)
	lockDepsFile, err := os.ReadFile(lockDepsFilePath)
	if err != nil {
		return TerraformApplyResult{}, err
	}

	return TerraformApplyResult{
		Outputs:          out,
		State:            base64.StdEncoding.EncodeToString(stateFile),
		DependenciesLock: base64.StdEncoding.EncodeToString(lockDepsFile),
	}, nil
}
