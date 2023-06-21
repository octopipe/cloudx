package main

import (
	"fmt"
	"os"

	"github.com/octopipe/cloudx/cmd/cli/commands"
	"github.com/octopipe/cloudx/internal/pluginmanager"
	"github.com/octopipe/cloudx/internal/terraform"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	terraformProvider, err := terraform.NewTerraformProvider(logger)
	if err != nil {
		panic(err)
	}
	pluginManager := pluginmanager.NewPluginManager(logger, terraformProvider)
	pluginCmd := commands.NewPluginRoot(pluginManager, terraformProvider)
	sharedInfraCmd := commands.NewSharedInfraRoot(pluginManager)

	commands.RootCmd.AddCommand(pluginCmd)
	commands.RootCmd.AddCommand(sharedInfraCmd)

	if err := commands.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
