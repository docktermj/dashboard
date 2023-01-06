// The version package simply prints the version of the go-fileindex binary file.
package version

import (
	"fmt"
	"github.com/spf13/viper"
)

// ----------------------------------------------------------------------------
// Version
// ----------------------------------------------------------------------------

// Print a version string.
func Version(programName string, buildVersion string, buildIteration string) {
	fmt.Printf("%s version %s-%s\n", programName, buildVersion, buildIteration)
}

// ----------------------------------------------------------------------------
// Command pattern "Execute" function.
// ----------------------------------------------------------------------------

// The Command sofware design pattern's Execute() method.
func Execute() {

	// Get parameters from viper.

	var programName string = viper.GetString("program_name")
	var buildVersion string = viper.GetString("build_version")
	var buildIteration string = viper.GetString("build_iteration")

	// Perform command.

	Version(programName, buildVersion, buildIteration)
}
