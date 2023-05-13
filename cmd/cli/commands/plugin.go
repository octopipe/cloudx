package commands

import (
	"encoding/json"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/octopipe/cloudx/internal/pluginmanager"
	"github.com/spf13/cobra"
)

type pluginCmd struct {
	pluginManager pluginmanager.Manager
}

func (p pluginCmd) NewPluginCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "plugin",
		Short: "Manage plugins",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
}

func (p pluginCmd) NewPublishPluginCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "publish",
		Short: "publish a plugin",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			contents := map[string][]byte{}
			filepath.Walk(args[1], func(path string, info fs.FileInfo, err error) error {
				if info.IsDir() {
					return nil
				}

				file, err := os.ReadFile(path)
				if err != nil {
					return err
				}

				contents[info.Name()] = file

				return nil
			})

			err := p.pluginManager.Publish(args[0], contents)
			if err != nil {
				log.Fatalln(err)
			}

		},
	}
}

func (p pluginCmd) NewExecutPluginCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "execute",
		Short: "execute plugin",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			inputPath := args[1]

			d, err := os.ReadFile(inputPath)
			if err != nil {
				log.Fatalln(err)
			}

			input := map[string]interface{}{}

			err = json.Unmarshal(d, &input)
			if err != nil {
				log.Fatalln(err)
			}

			err = p.pluginManager.Execute(name, input)
			if err != nil {
				log.Fatalln(err)
			}
		},
	}
}

func NewPluginRoot(pluginManager pluginmanager.Manager) *cobra.Command {
	pluginRoot := pluginCmd{
		pluginManager: pluginManager,
	}

	pluginCmd := pluginRoot.NewPluginCmd()
	pluginCmd.AddCommand(pluginRoot.NewPublishPluginCmd())
	pluginCmd.AddCommand(pluginRoot.NewExecutPluginCmd())

	return pluginCmd
}
