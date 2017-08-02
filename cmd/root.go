package cmd

import (
	"fmt"
	"os"

	"github.com/lytics/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const defaultCfgFileName = "grpcapi"

var (
	cfgFile     string
	composeFile string

	serverConfig map[string]string
	clientConfig map[string]string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "grpcapi",
	Short: "KETI-GRPC API",
	Long:  ``,
	PreRun: func(cmd *cobra.Command, args []string) {
		logrus.SetLevel(logrus.DebugLevel)
	},
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: ./grpcapi.yaml)")
	RootCmd.PersistentFlags().StringVar(&composeFile, "compose", "", "docker compose file (example: docker-compose.yml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	viper.AutomaticEnv() // read in environment variables that match

	if composeFile != "" {
		cfgFile = ""
	} else {

		if cfgFile != "" {
			// Use config file from the flag.
			viper.SetConfigFile(cfgFile)
		} else {
			// Search config in current directory with name ".grpcapi.yaml"
			viper.SetConfigType("yaml")
			viper.AddConfigPath(".")
			viper.SetConfigName(defaultCfgFileName)
		}

		// If a config file is found, read it in.
		viper.ReadInConfig()
		cfgFile = viper.ConfigFileUsed()

		serverConfig = viper.GetStringMapString("api-server")
		clientConfig = viper.GetStringMapString("api-client")
	}

	if cfgFile == "" && composeFile == "" {
		fmt.Println("Flag: Either config or compose is required.")
		os.Exit(1)
	}
}
