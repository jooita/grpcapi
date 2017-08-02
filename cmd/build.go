package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "API Server && Client Dockerfile Build",
	Long:  `API Server && Client Dockerfile Build`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("build called")
	},
}

func init() {
	RootCmd.AddCommand(buildCmd)
}
