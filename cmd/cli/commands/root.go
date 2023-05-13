package commands

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "cloudx",
	Short: "Cloudx is a developer platform",
	Long:  `Cloudx is a developer platform`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}
