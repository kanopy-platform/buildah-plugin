// Drone specific command line parameters or environment variable settings
package drone

import (
	"github.com/spf13/viper"
)

func AdditionalSetup() {
	// all settings passed to plugin from Drone are prefixed with PLUGIN_
	viper.SetEnvPrefix("PLUGIN")
}
