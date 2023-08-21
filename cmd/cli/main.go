package main

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/octopipe/cloudx/cmd/cli/commands"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	restclient := resty.New()
	// taskManager := taskmanager.NewTaskManager(logger)
	// taskCmd := commands.NewTaskRoot(taskManager)
	// infraCmd := commands.NewInfraRoot(taskManager)
	repositoryCmd := commands.NewRepositoryRoot(logger, restclient)

	commands.RootCmd.AddCommand(repositoryCmd)

	if err := commands.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
