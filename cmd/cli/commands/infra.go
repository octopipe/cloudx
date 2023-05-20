package commands

import (
	"github.com/octopipe/cloudx/internal/pluginmanager"
	"github.com/spf13/cobra"
)

type stackSetCmd struct {
}

func (p stackSetCmd) NewStackSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stackSet",
		Short: "Manage stackSets",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
}

func (p stackSetCmd) NewCreateStackSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "create a stackSet",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
}

func NewStackSetRoot(pluginManager pluginmanager.Manager) *cobra.Command {
	stackSetRoot := stackSetCmd{}

	stackSetCmd := stackSetRoot.NewStackSetCmd()
	stackSetCmd.AddCommand(stackSetRoot.NewCreateStackSetCmd())

	return stackSetCmd
}
