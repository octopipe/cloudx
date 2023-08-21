package commands

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/octopipe/cloudx/internal/taskmanager"
	"github.com/spf13/cobra"
)

type taskCmd struct {
	taskManager taskmanager.Manager
}

func (p taskCmd) NewTaskCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "task",
		Short: "Manage tasks",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
}

func (p taskCmd) NewPublishTaskCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "publish",
		Short: "publish a task",
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

			err := p.taskManager.Publish(args[0], contents)
			if err != nil {
				log.Fatalln(err)
			}

		},
	}
}

func (p taskCmd) NewExecutTaskCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "execute",
		Short: "execute task",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			// name := args[0]
			// inputPath := args[1]

			// d, err := os.ReadFile(inputPath)
			// if err != nil {
			// 	log.Fatalln(err)
			// }

			// input := map[string]interface{}{}

			// err = json.Unmarshal(d, &input)
			// if err != nil {
			// 	log.Fatalln(err)
			// }

			// _, _, err = p.terraformProvider.Apply(name, input)
			// if err != nil {
			// 	log.Fatalln(err)
			// }
		},
	}
}

func NewTaskRoot(taskManager taskmanager.Manager) *cobra.Command {
	taskRoot := taskCmd{
		taskManager: taskManager,
	}

	taskCmd := taskRoot.NewTaskCmd()
	taskCmd.AddCommand(taskRoot.NewPublishTaskCmd())
	taskCmd.AddCommand(taskRoot.NewExecutTaskCmd())

	return taskCmd
}
