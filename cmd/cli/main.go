package main

import (
	"fmt"
	"os"

	"github.com/octopipe/cloudx/cmd/cli/commands"
	"github.com/octopipe/cloudx/internal/taskmanager"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	taskManager := taskmanager.NewTaskManager(logger)
	taskCmd := commands.NewTaskRoot(taskManager)
	infraCmd := commands.NewInfraRoot(taskManager)

	commands.RootCmd.AddCommand(taskCmd)
	commands.RootCmd.AddCommand(infraCmd)

	if err := commands.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
