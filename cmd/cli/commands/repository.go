package commands

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/go-resty/resty/v2"
	commonv1alpha1 "github.com/octopipe/cloudx/apis/common/v1alpha1"
	"github.com/octopipe/cloudx/internal/repository"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type repositoryCmd struct {
	logger     *zap.Logger
	restclient *resty.Client
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
			var commonQuestions = []*survey.Question{
				{
					Name:      "name",
					Prompt:    &survey.Input{Message: "What is the name of repository?"},
					Validate:  survey.Required,
					Transform: survey.Title,
				},
				{
					Name:      "namespace",
					Prompt:    &survey.Input{Message: "What is the namespace of repository?", Default: "default"},
					Validate:  survey.Required,
					Transform: survey.Title,
				},
				{
					Name:      "url",
					Prompt:    &survey.Input{Message: "What is the url of repository?"},
					Validate:  survey.Required,
					Transform: survey.Title,
				},
				{
					Name:      "path",
					Prompt:    &survey.Input{Message: "What is the path inside of repository?", Default: "."},
					Validate:  survey.Required,
					Transform: survey.Title,
				},
				{
					Name:      "branch",
					Prompt:    &survey.Input{Message: "What is the branch of repository?", Default: "main"},
					Validate:  survey.Required,
					Transform: survey.Title,
				},
				{
					Name: "public",
					Prompt: &survey.Confirm{
						Message: "This repository is public?",
						Default: true,
					},
				},
				// {
				// 	Name: "authMethod",
				// 	Prompt: &survey.Select{
				// 		Message: "Choose a auth method:",
				// 		Options: []string{"username/password", "username/accessKey", "sshPublicKey"},
				// 		Default: "username/accessKey",
				// 	},
				// },
			}

			commonAnswers := struct {
				Name      string
				Namespace string
				Url       string
				Path      string
				Branch    string
				Public    bool
			}{}

			// perform the questions
			err := survey.Ask(commonQuestions, &commonAnswers)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			_, err = p.restclient.R().
				SetBody(repository.Repository{
					Name:      commonAnswers.Name,
					Namespace: commonAnswers.Namespace,
					RepositorySpec: commonv1alpha1.RepositorySpec{
						Sync: commonv1alpha1.RepositorySync{
							Auto: true,
						},
						Url:    commonAnswers.Url,
						Path:   commonAnswers.Path,
						Branch: commonAnswers.Branch,
					},
				}).
				Post(fmt.Sprintf("%s/repositories", os.Getenv("APISERVER_BASE_PATH")))

			if err != nil {
				fmt.Println(err.Error())
				return
			}

		},
	}
}

func NewRepositoryRoot(logger *zap.Logger, restclient *resty.Client) *cobra.Command {
	repositoryRoot := repositoryCmd{
		logger:     logger,
		restclient: restclient,
	}

	repositoryCmd := repositoryRoot.NewRepositoryCmd()
	repositoryCmd.AddCommand(repositoryRoot.NewCreateRepositoryCmd())

	return repositoryCmd
}
