package main

import (
	"fmt"
	"os"

	"github.com/octopipe/cloudx/cmd/cli/commands"
	"github.com/octopipe/cloudx/internal/taskmanager"
	"github.com/octopipe/cloudx/internal/terraform"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	terraformProvider, err := terraform.NewTerraformProvider(logger)
	if err != nil {
		panic(err)
	}
	taskManager := taskmanager.NewTaskManager(logger, terraformProvider)
	taskCmd := commands.NewTaskRoot(taskManager, terraformProvider)
	infraCmd := commands.NewInfraRoot(taskManager)

	commands.RootCmd.AddCommand(taskCmd)
	commands.RootCmd.AddCommand(infraCmd)

	if err := commands.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
