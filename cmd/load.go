package cmd

import (
	"github.com/docktermj/dashboard/cmd/load"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Variables that can be set from command-line interface.

var (
	jsonFileName   string
	sqliteFileName string
)

// loadCmd represents the load command.

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load the database",
	Long:  `Load the database.`,
	Run: func(cmd *cobra.Command, args []string) {

		// Pass CLI options into viper.
		// Note: viper variables have underscores.

		viper.BindPFlag("json_file_name", cmd.PersistentFlags().Lookup("json-file-name"))
		viper.BindPFlag("sqlite_file_name", cmd.PersistentFlags().Lookup("sqlite-file-name"))

		load.Execute()
	},
}

func init() {

	// Define subcommmand command line interface (CLI) parameters using cobra.
	// Note: CLI options have hyphens.

	loadCmd.PersistentFlags().StringVar(&jsonFileName, "json-file-name", "go-fileindex.json", "GOFILEINDEX_JSON_FILE_NAME")
	loadCmd.PersistentFlags().StringVar(&sqliteFileName, "sqlite-file-name", "go-fileindex.db", "GOFILEINDEX_SQLITE_FILE_NAME")

	// Add subcommand to observer pattern.

	rootCmd.AddCommand(loadCmd)
}
