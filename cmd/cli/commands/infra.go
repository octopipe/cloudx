package commands

import (
	"github.com/octopipe/cloudx/internal/pluginmanager"
	"github.com/spf13/cobra"
)

type sharedInfraCmd struct {
}

func (p sharedInfraCmd) NewSharedInfraCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sharedInfra",
		Short: "Manage sharedInfras",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
}

func (p sharedInfraCmd) NewCreateSharedInfraCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "create a sharedInfra",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
}

func NewSharedInfraRoot(pluginManager pluginmanager.Manager) *cobra.Command {
	sharedInfraRoot := sharedInfraCmd{}

	sharedInfraCmd := sharedInfraRoot.NewSharedInfraCmd()
	sharedInfraCmd.AddCommand(sharedInfraRoot.NewCreateSharedInfraCmd())

	return sharedInfraCmd
}
