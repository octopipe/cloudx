package commands

import (
	"github.com/octopipe/cloudx/internal/pluginmanager"
	"github.com/spf13/cobra"
)

type stackCmd struct {
}

func (p stackCmd) NewStackCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stack",
		Short: "Manage stacks",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
}

func (p stackCmd) NewCreateStackCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "create a stack",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
}

func NewStackRoot(pluginManager pluginmanager.Manager) *cobra.Command {
	stackRoot := stackCmd{}

	stackCmd := stackRoot.NewStackCmd()
	stackCmd.AddCommand(stackRoot.NewCreateStackCmd())

	return stackCmd
}
