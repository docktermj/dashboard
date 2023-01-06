package cmd

import (
	"github.com/docktermj/dashboard/cmd/scan"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Variables that can be set from command-line interface.

var (
	outputFileName string
	rootPath       string
	volumeName     string
)

// scanCmd represents the scan command.

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan volume for file metadata",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		// Pass CLI options into viper.
		// Note: viper variables have underscores.

		viper.BindPFlag("output_file_name", cmd.PersistentFlags().Lookup("output-file-name"))
		viper.BindPFlag("root_path", cmd.PersistentFlags().Lookup("root-path"))
		viper.BindPFlag("volume_name", cmd.PersistentFlags().Lookup("volume-name"))

		scan.Execute()
	},
}

func init() {

	// Define subcommmand command line interface (CLI) parameters using cobra.
	// Note: CLI options have hyphens.

	scanCmd.PersistentFlags().StringVar(&outputFileName, "output-file-name", "go-fileindex.json", "GOFILEINDEX_OUTPUT_FILE_NAME")
	scanCmd.PersistentFlags().StringVar(&rootPath, "root-path", "", "GOFILEINDEX_ROOT_PATH")
	scanCmd.PersistentFlags().StringVar(&volumeName, "volume-name", "", "GOFILEINDEX_VOLUME_NAME")

	// Add subcommand to observer pattern.

	rootCmd.AddCommand(scanCmd)
}
