package cli

import (
	"fmt"
	"strings"

	"github.com/kanopy-platform/buildah-plugin/internal/version"
	"github.com/kanopy-platform/buildah-plugin/pkg/buildah"
	"github.com/kanopy-platform/buildah-plugin/pkg/buildah/manifest"
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

	cmd.PersistentFlags().String("log-level", "info", "Configure log level")
	cmd.PersistentFlags().String("access-key", "", "AWS Access Key for ECR authentication")
	cmd.PersistentFlags().String("secret-key", "", "AWS Secret Key for ECR authentication")
	cmd.PersistentFlags().String("registry", "", "ECR registry")
	cmd.PersistentFlags().Bool("public-registry", false, "Public ECR registry")
	cmd.PersistentFlags().String("repo", "", "The repository in the ECR registry")
	cmd.PersistentFlags().String("command", "version", "The buildah command to run")
	addManifestCommandFlags(cmd)

	cmd.AddCommand(newVersionCommand())

	return cmd
}

func addManifestCommandFlags(cmd *cobra.Command) {
	// flags specific to "manifest" command
	cmd.PersistentFlags().StringSlice("manifest-sources", []string{}, "List of image tags to add to final manifest")
	cmd.PersistentFlags().StringSlice("manifest-targets", []string{}, "List of tags to associate with final manifest")
}

func (c *RootCommand) persistentPreRunE(cmd *cobra.Command, args []string) error {
	// additional settings based on plugin type
	if err := pluginTypeSetup(); err != nil {
		return err
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
	// TODO populate AWS ECR provider with auth credentials

	buildah := buildah.Buildah{
		Login: buildah.Login{
			Registry: viper.GetString("registry"),
			Username: "AWS",      // TODO use output from AWS ECR provider
			Password: "password", // TODO use output from AWS ECR provider
		},
		Repo:    viper.GetString("repo"),
		Command: viper.GetString("command"),
		Manifest: manifest.CommandArgs{
			Sources: viper.GetStringSlice("manifest-sources"),
			Targets: viper.GetStringSlice("manifest-targets"),
		},
	}

	return buildah.Exec()
}

func pluginTypeSetup() error {
	pluginType := version.Get().PluginType

	switch pluginType {
	case version.PluginTypeDrone:
		viper.SetEnvPrefix("PLUGIN")
	default:
		return fmt.Errorf("invalid plugin type: %q", pluginType)
	}

	return nil
}
