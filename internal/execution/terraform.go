package execution

import (
	"context"
	"log"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
)

func NewTerraformExecution(workdirPath string) (output []commonv1alpha1.StackPluginOutput, err error) {
	installer := &releases.ExactVersion{
		Product: product.Terraform,
		Version: version.Must(version.NewVersion("1.0.6")),
	}

	execPath, err := installer.Install(context.Background())
	if err != nil {
		log.Fatalf("error installing Terraform: %s", err)
	}

	tf, err := tfexec.NewTerraform(workdirPath, execPath)
	if err != nil {
		log.Fatalf("error running NewTerraform: %s", err)
	}

	err = tf.Init(context.Background(), tfexec.Upgrade(true))
	if err != nil {
		log.Fatalf("error running Init: %s", err)
	}

	err = tf.Apply(context.Background())
	if err != nil {
		log.Fatalf("error running apply: %s", err)
	}

	out, err := tf.Output(context.Background())
	if err != nil {
		log.Fatalf("error running output: %s", err)
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
