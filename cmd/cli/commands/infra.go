package commands

import (
	"github.com/octopipe/cloudx/internal/taskmanager"
	"github.com/spf13/cobra"
)

type infraCmd struct {
}

func (p infraCmd) NewInfraCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "infra",
		Short: "Manage infras",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
}

func (p infraCmd) NewCreateInfraCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "create a infra",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
}

func NewInfraRoot(taskManager taskmanager.Manager) *cobra.Command {
	infraRoot := infraCmd{}

	infraCmd := infraRoot.NewInfraCmd()
	infraCmd.AddCommand(infraRoot.NewCreateInfraCmd())

	return infraCmd
}
