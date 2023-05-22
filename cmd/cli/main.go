package main

import (
	"fmt"
	"os"

	"github.com/octopipe/cloudx/cmd/cli/commands"
	"github.com/octopipe/cloudx/internal/pluginmanager"
	"github.com/octopipe/cloudx/internal/terraform"
)

func main() {
	terraformProvider, err := terraform.NewTerraformProvider()
	if err != nil {
		panic(err)
	}
	pluginManager := pluginmanager.NewPluginManager(terraformProvider)
	pluginCmd := commands.NewPluginRoot(pluginManager)
	sharedInfraCmd := commands.NewSharedInfraRoot(pluginManager)

	commands.RootCmd.AddCommand(pluginCmd)
	commands.RootCmd.AddCommand(sharedInfraCmd)

	if err := commands.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
