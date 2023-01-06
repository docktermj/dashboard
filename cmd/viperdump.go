package cmd

import (
	"github.com/docktermj/dashboard/cmd/viperdump"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Variables that can be set from command-line interface.

var (
	flagKeyOnly string
	flagKey     string
	setKey      string
)

// viperdumpCmd represents the viperdump command.

var viperdumpCmd = &cobra.Command{
	Use:   "viperdump",
	Short: "Dump the contents of viper",
	Long: `viper holds the context of the running program.
Dump those contents.`,
	Run: func(cmd *cobra.Command, args []string) {

		// Priority 0. Set explicitly.  User cannot modify.

		viper.Set("set_key_only", "From Set()")

		// Priority 1. Set from CLI.  Note: CLI variables have hyphens.

		viper.BindPFlag("flag_key", cmd.PersistentFlags().Lookup("flag-key"))
		viper.BindPFlag("flag_key_only", cmd.PersistentFlags().Lookup("flag-key-only"))
		viper.BindPFlag("set_key", cmd.PersistentFlags().Lookup("set-key"))

		// Priority 2. Set from environment variables.

		viper.SetEnvPrefix("gofileindex")
		viper.BindEnv("OS_KEY")
		viper.BindEnv("OS_KEY_ONLY")
		viper.BindEnv("FLAG_KEY")
		viper.BindEnv("SET_KEY")
		viper.AutomaticEnv() // read in environment variables that match

		// Priority 3. Set from config file.  Done in root.go

		// Priority 4. Set from default values.

		viper.SetDefault("flag_key", "From SetDefault()")
		viper.SetDefault("os_key", "From SetDefault()")
		viper.SetDefault("configuration_key", "From SetDefault()")
		viper.SetDefault("default_key", "From SetDefault()")
		viper.SetDefault("set_key", "From SetDefault()")

		// Command pattern software design pattern. Execute command.

		viperdump.Execute()
	},
}

// As part of "observer pattern", init() is run to register the observer.
// It is alway run, even if the subcommand is not specified.
func init() {

	// Define subcommmand command line interface (CLI) parameters using cobra.
	// Note: CLI options have hyphens.

	viperdumpCmd.PersistentFlags().StringVar(&flagKey, "flag-key", "From SetDefault()", "GOFILEINDEX_FLAG_KEY")
	viperdumpCmd.PersistentFlags().StringVar(&flagKeyOnly, "flag-key-only", "none", "flagKeyOnly")
	viperdumpCmd.PersistentFlags().StringVar(&setKey, "set-key", "From SetDefault()", "GOFILEINDEX_SET_KEY")

	// Add subcommand to observer pattern.

	rootCmd.AddCommand(viperdumpCmd)
}
