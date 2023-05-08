package terraform

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/hashicorp/terraform-exec/tfexec"
	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
)

type execution struct {
	sync.RWMutex
	execPath        string
	executedPlugins map[string][]commonv1alpha1.StackPluginOutput
}

func NewExecution(execPath string) execution {
	return execution{
		execPath:        execPath,
		executedPlugins: make(map[string][]commonv1alpha1.StackPluginOutput),
	}
}

func (c *execution) isEnabled(dependencies []string) bool {
	if len(dependencies) <= 0 {
		return true
	}

	isEnabled := true
	for _, dependency := range dependencies {
		if _, ok := c.executedPlugins[dependency]; !ok {
			isEnabled = false
			break
		}
	}

	return isEnabled
}

func (c *execution) executeNextSteps(stackPlugins []commonv1alpha1.StackPlugin) {
	var wg sync.WaitGroup

	for _, plugin := range stackPlugins {
		if _, ok := c.executedPlugins[plugin.Name]; !ok && c.isEnabled(plugin.Depends) {
			wg.Add(1)

			go func(p commonv1alpha1.StackPlugin) {
				defer wg.Done()

				outputs, err := c.execute(p)
				if err != nil {
					panic(err)
				}

				c.Lock()
				c.executedPlugins[p.Name] = outputs
				c.Unlock()
			}(plugin)

		}
	}

	wg.Wait()

	fmt.Println("FINISH")

	if len(c.executedPlugins) == len(stackPlugins) {
		return
	}

	c.executeNextSteps(stackPlugins)
}

func (c *execution) execute(plugin commonv1alpha1.StackPlugin) (output []commonv1alpha1.StackPluginOutput, err error) {
	workdirPath := plugin.Ref.Path
	if plugin.Ref.Type != "local" {
		workdirPath = ""
	}

	tf, err := tfexec.NewTerraform(workdirPath, c.execPath)
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
