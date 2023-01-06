package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/docktermj/dashboard/service"
	"github.com/senzing/go-logging/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configurationFile string
	buildVersion      string = "0.0.0"
	buildIteration    string = "0"
)

func makeVersion(version string, iteration string) string {
	var result string = ""
	if buildIteration == "0" {
		result = version
	} else {
		result = fmt.Sprintf("%s-%s", buildVersion, buildIteration)
	}
	return result
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "Bring up the Senzing dashboard",
	Long:  `For more information, visit https://github.com/Senzing/dashboard`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error = nil
		ctx := context.TODO()

		logLevel, ok := logger.TextToLevelMap[viper.GetString("log-level")]
		if !ok {
			logLevel = logger.LevelInfo
		}

		httpServer := &service.HttpServerImpl{
			Port:     viper.GetInt("dashboard-port"),
			LogLevel: logLevel,
		}
		httpServer.Serve(ctx)

		return err
	},
	Version: makeVersion(buildVersion, buildIteration),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Define flags for command.

	RootCmd.Flags().Int("dashboard-port", 8258, "port used to serve HTTP [SENZING_TOOLS_DASHBOARD_PORT]")

	// Integrate with Viper.

	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("SENZING_TOOLS")

	// Define flags in Viper.

	viper.SetDefault("dashboard-port", 8258)
	viper.BindPFlag("dashboard-port", RootCmd.Flags().Lookup("dashboard-port"))

	viper.SetDefault("log-level", "INFO")
	viper.BindPFlag("log-level", RootCmd.Flags().Lookup("log-level"))

	// Set version template.

	versionTemplate := `{{printf "%s: %s - version %s\n" .Name .Short .Version}}`
	RootCmd.SetVersionTemplate(versionTemplate)

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	// Set logging output format.

	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds | log.LUTC)

	// Load configuration file.

	if configurationFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configurationFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".servegrpc" (without extension).
		viper.AddConfigPath(home + "/.senzing-tools")
		viper.AddConfigPath(home)
		viper.AddConfigPath("/etc/senzing-tools")
		viper.SetConfigType("yaml")
		viper.SetConfigName("dashboard")
	}

	// Read in environment variables that match "SENZING_TOOLS_*" pattern.

	viper.AutomaticEnv()

	// If a config file is found, read it in.

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
