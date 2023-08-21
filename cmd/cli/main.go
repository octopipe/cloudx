package main

import (
	"fmt"
	"os"

	"github.com/octopipe/cloudx/cmd/cli/commands"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	// taskManager := taskmanager.NewTaskManager(logger)
	// taskCmd := commands.NewTaskRoot(taskManager)
	// infraCmd := commands.NewInfraRoot(taskManager)
	repositoryCmd := commands.NewRepositoryRoot(logger)

	commands.RootCmd.AddCommand(repositoryCmd)

	if err := commands.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
