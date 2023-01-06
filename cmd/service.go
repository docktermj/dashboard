package cmd

import (
	"github.com/docktermj/dashboard/cmd/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Variables that can be set from command-line interface.

var (
	port string
)

// serviceCmd represents the service command.

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		// Pass CLI options into viper.
		// Note: viper variables have underscores.

		viper.BindPFlag("port", cmd.PersistentFlags().Lookup("port"))
		viper.BindPFlag("sqlite_file_name", cmd.PersistentFlags().Lookup("sqlite-file-name"))

		service.Execute()
	},
}

func init() {

	// Define subcommmand command line interface (CLI) parameters using cobra.
	// Note: CLI options have hyphens.

	serviceCmd.PersistentFlags().StringVar(&port, "port", "3000", "GOFILEINDEX_PORT")
	serviceCmd.PersistentFlags().StringVar(&sqliteFileName, "sqlite-file-name", "go-fileindex.db", "GOFILEINDEX_SQLITE_FILE_NAME")

	// Add subcommand to observer pattern.

	rootCmd.AddCommand(serviceCmd)
}
