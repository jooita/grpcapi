package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"beji/grpcapi/compose"
)

var (
	debug   bool
	nr      bool
	fr      bool
	nb      bool
	build   bool
	timeout int

	upCmd *cobra.Command
)

func init() {

	// upCmd represents the up command
	upCmd = &cobra.Command{
		Use:   "up",
		Short: "Create and start API Server/Client",
		Long:  `Create and start API Server/Client`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("up called")
			return runUp()
		},
	}

	RootCmd.AddCommand(upCmd)
	upCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Do not block and log")
	upCmd.Flags().BoolVar(&nr, "no-recreate", false, "If containers already exist, don't recreate them. Incompatible with --force-recreate")
	upCmd.Flags().BoolVar(&fr, "force-recreate", false, "Recreate containers even if their configuration and image haven't changed. Incompatible with --no-recreate")
	upCmd.Flags().BoolVar(&nb, "no-build", false, "Don't build an image, even if it's missing")
	upCmd.Flags().BoolVar(&build, "build", false, "Build images before starting containers.")
	upCmd.Flags().IntVarP(&timeout, "timeout", "t", 0, "Specify a shutdown timeout in seconds.")
}

func runUp() error {
	composefile := RootCmd.Flag("compose").Value.String()
	flags := upCmd.Flags()

	if cfgFile != "" {
		p := compose.NewProject(serverConfig, clientConfig)
		if err := p.Up(flags); err != nil {
			return err
		}
	} else {
		if err := compose.ComposeFileUp(composefile); err != nil {
			return err
		}
	}
	return nil
}
