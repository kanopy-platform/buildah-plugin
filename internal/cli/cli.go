package cli

import (
	"fmt"
	"strings"

	"github.com/kanopy-platform/buildah-plugin/internal/cli/drone"
	"github.com/kanopy-platform/buildah-plugin/internal/version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type RootCommand struct{}

func NewRootCommand() *cobra.Command {
	root := &RootCommand{}

	cmd := &cobra.Command{
		Use:               "buildah-plugin",
		Short:             "Plugin for CI/CD tools to run buildah commands.",
		PersistentPreRunE: root.persistentPreRunE,
		RunE:              root.runE,
	}

	cmd.AddCommand(newVersionCommand())

	return cmd
}

func (c *RootCommand) persistentPreRunE(cmd *cobra.Command, args []string) error {
	// add additional settings based on plugin type
	switch version.Get().PluginType {
	case version.PluginTypeDrone:
		drone.AdditionalSetup()
	default:
		return fmt.Errorf("must specify a plugin type")
	}
	// bind flags to viper
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return err
	}

	// set log level
	logLevel, err := log.ParseLevel(viper.GetString("log-level"))
	if err != nil {
		return err
	}
	log.SetLevel(logLevel)

	return nil
}

func (c *RootCommand) runE(cmd *cobra.Command, args []string) error {
	return nil
}
