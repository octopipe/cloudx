package commands

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type repositoryCmd struct {
}

func (p repositoryCmd) NewRepositoryCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "repository",
		Short: "Manage repository",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
}

func (p repositoryCmd) NewCreateRepositoryCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create",
		Short: "create a repository",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
}

func NewRepositoryRoot(logger *zap.Logger) *cobra.Command {
	repositoryRoot := repositoryCmd{}

	repositoryCmd := repositoryRoot.NewRepositoryCmd()
	repositoryCmd.AddCommand(repositoryRoot.NewCreateRepositoryCmd())

	return repositoryCmd
}
