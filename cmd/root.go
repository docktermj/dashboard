/*
 */
package cmd

import (
	"context"
	"os"

	"github.com/docktermj/dashboard/dashboard"
	"github.com/senzing/go-cmdhelping/cmdhelper"
	"github.com/senzing/go-cmdhelping/option"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	Short string = "Bring up the Senzing dashboard"
	Use   string = "dashboard"
	Long  string = `
For more information, visit https://github.com/Senzing/dashboard
    `
)

// ----------------------------------------------------------------------------
// Context variables
// ----------------------------------------------------------------------------

var ContextVariablesForMultiPlatform = []option.ContextVariable{
	option.Configuration,
	option.EngineConfigurationJson,
	option.HttpPort,
	option.LogLevel,
}

var ContextVariables = append(ContextVariablesForMultiPlatform, ContextVariablesForOsArch...)

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

// Since init() is always invoked, define command line parameters.
func init() {
	cmdhelper.Init(RootCmd, ContextVariables)
}

// ----------------------------------------------------------------------------
// Public functions
// ----------------------------------------------------------------------------

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// Used in construction of cobra.Command
func PreRun(cobraCommand *cobra.Command, args []string) {
	cmdhelper.PreRun(cobraCommand, args, Use, ContextVariables)
}

// Used in construction of cobra.Command
func RunE(_ *cobra.Command, _ []string) error {
	ctx := context.TODO()

	httpServer := &dashboard.DashboardImpl{
		ServerPort: viper.GetInt(option.HttpPort.Arg),
	}
	return httpServer.Serve(ctx)
}

// Used in construction of cobra.Command
func Version() string {
	return cmdhelper.Version(githubVersion, githubIteration)
}

// ----------------------------------------------------------------------------
// Command
// ----------------------------------------------------------------------------

// RootCmd represents the command.
var RootCmd = &cobra.Command{
	Use:     Use,
	Short:   Short,
	Long:    Long,
	PreRun:  PreRun,
	RunE:    RunE,
	Version: Version(),
}
