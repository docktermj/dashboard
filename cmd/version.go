package cmd

import (
	"github.com/docktermj/dashboard/cmd/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Values updated via "go install -ldflags" parameters.

var (
	programName    string = "unknown"
	buildVersion   string = "0.0.0"
	buildIteration string = "0"
)

// versionCmd represents the version command.

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version of fileindex",
	Long:  `Show the version of the fileindex program`,
	Run: func(cmd *cobra.Command, args []string) {

		// Set build parameters in viper.

		viper.Set("program_name", programName)
		viper.Set("build_version", buildVersion)
		viper.Set("build_iteration", buildIteration)

		version.Execute()
	},
}

// As part of "observer pattern", init() is run to register the observer.
// It is alway run, even if the subcommand is not specified.
func init() {

	// Register subcommand ("observer pattern").

	rootCmd.AddCommand(versionCmd)
}
