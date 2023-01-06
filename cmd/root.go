package cmd

import (
	"fmt"
	"github.com/docktermj/go-logger/logger"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

var (
	configurationFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fileindex",
	Short: "Manage a database of file information",
	Long: `The steps for managing file information:

1) Scan files
2) Populate database
3) Query

For more information, visit https://github.com/docktermj/go-fileindex`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	// Initialize cobra.

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&configurationFile, "config", "", "config file (default is $HOME/.fileindex.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	// Set logging output format.

	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds | log.LUTC)

	// Configure the logger. If not configured, no functions will print.

	logger.SetLevel(logger.LevelInfo)

	// Load configuration file.

	if configurationFile != "" {
		viper.SetConfigFile(configurationFile) // Use config file from the flag.
	} else {

		// Find home directory.

		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Configuration file will be ${HOME}/.fileindex.yaml

		viper.SetConfigName(".fileindex")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(home)

	}

	// Enable automatic inclusing of environment variables.

	viper.SetEnvPrefix("gofileindex")
	viper.BindEnv("DEBUG")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.

	if err := viper.ReadInConfig(); err == nil {
		logger.Info("Using config file:", viper.ConfigFileUsed())
	}
}
