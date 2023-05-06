package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
	"gopkg.in/yaml.v2"
)

type StackPluginRef struct {
	Type string `yaml:"type"`
	Path string `yaml:"path"`
}

type StackPluginOutput struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type StackPluginInput struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type StackPlugin struct {
	Name    string             `yaml:"name"`
	Ref     StackPluginRef     `yaml:"ref"`
	Depends []string           `yaml:"depends"`
	Inputs  []StackPluginInput `yaml:"inputs"`
}

type StackSpec struct {
	Plugins []StackPlugin `yaml:"plugins"`
}

type Stack struct {
	Kind string    `yaml:"kind"`
	Spec StackSpec `yaml:"spec"`
}

func main() {
	installer := &releases.ExactVersion{
		Product: product.Terraform,
		Version: version.Must(version.NewVersion("1.0.6")),
	}

	execPath, err := installer.Install(context.Background())
	if err != nil {
		log.Fatalf("error installing Terraform: %s", err)
	}

	stkValue, err := os.ReadFile("./workdir/stk.yaml")
	if err != nil {
		panic(err)
	}
	newStack := Stack{}
	err = yaml.Unmarshal(stkValue, &newStack)
	if err != nil {
		panic(err)
	}

	executionContext := ExecutionContext{
		execPath:        execPath,
		executedPlugins: make(map[string][]StackPluginOutput),
	}

	executionContext.executeNextSteps(newStack.Spec.Plugins)

}

type ExecutionContext struct {
	sync.RWMutex
	execPath        string
	executedPlugins map[string][]StackPluginOutput
}

func (c *ExecutionContext) isEnabled(dependencies []string) bool {
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

func (c *ExecutionContext) executeNextSteps(stackPlugins []StackPlugin) {
	var wg sync.WaitGroup

	for _, plugin := range stackPlugins {
		if _, ok := c.executedPlugins[plugin.Name]; !ok && c.isEnabled(plugin.Depends) {
			wg.Add(1)

			go func(p StackPlugin) {
				defer wg.Done()

				outputs, err := c.execute(p)
				if err != nil {
					panic(err)
				}

				fmt.Println("OUTPUTS", p.Name, outputs)

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

func (c *ExecutionContext) execute(plugin StackPlugin) (output []StackPluginOutput, err error) {
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

	outputs := []StackPluginOutput{}

	for key, value := range out {
		outputs = append(outputs, StackPluginOutput{
			Key:   key,
			Value: string(value.Value),
		})
	}

	return outputs, nil
}
