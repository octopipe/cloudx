package main

import (
	"fmt"
	"os"

	"github.com/octopipe/cloudx/cmd/cli/commands"
	"github.com/octopipe/cloudx/internal/pluginmanager"
)

func main() {
	pluginManager := pluginmanager.NewPluginManager()
	pluginCmd := commands.NewPluginRoot(pluginManager)

	commands.RootCmd.AddCommand(pluginCmd)

	if err := commands.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}