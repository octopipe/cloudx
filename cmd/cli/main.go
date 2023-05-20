package main

import (
	"fmt"
	"os"

	"github.com/octopipe/cloudx/cmd/cli/commands"
	"github.com/octopipe/cloudx/internal/pluginmanager"
	"github.com/octopipe/cloudx/internal/provider"
)

func main() {
	terraformProvider, err := provider.NewTerraformProvider()
	if err != nil {
		panic(err)
	}
	pluginManager := pluginmanager.NewPluginManager(terraformProvider)
	pluginCmd := commands.NewPluginRoot(pluginManager)
	stackSetCmd := commands.NewStackSetRoot(pluginManager)

	commands.RootCmd.AddCommand(pluginCmd)
	commands.RootCmd.AddCommand(stackSetCmd)

	if err := commands.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
