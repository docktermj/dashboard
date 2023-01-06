package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// addFilesCmd represents the addFiles command
var addFilesCmd = &cobra.Command{
	Use:   "addFiles",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples

to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("addFiles called")
	},
}

func init() {

	// FIXME:  Work in progress.

	rootCmd.AddCommand(addFilesCmd)

}
