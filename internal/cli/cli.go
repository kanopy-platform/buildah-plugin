package cli

import (
	"fmt"
	"strings"

	buildversion "github.com/kanopy-platform/buildah-plugin/internal/version"
	"github.com/kanopy-platform/buildah-plugin/pkg/buildah"
	"github.com/kanopy-platform/buildah-plugin/pkg/buildah/manifest"
	"github.com/kanopy-platform/buildah-plugin/pkg/buildah/version"
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
	cmd.PersistentFlags().String("repo", "", "The repository in the ECR registry")
	cmd.PersistentFlags().Bool("print-version", false, "Prints plugin and buildah version for debugging")
	addManifestCommandFlags(cmd)

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
	if viper.GetBool("print-version") {
		log.Infof("%#v\n", buildversion.Get())
	}

	// TODO get password from AWS ECR provider

	buildah := buildah.Buildah{
		Login: buildah.Login{
			Registry: viper.GetString("registry"),
			Username: "AWS",      // TODO use output from AWS ECR provider
			Password: "password", // TODO use output from AWS ECR provider
		},
		Repo: viper.GetString("repo"),
		Version: version.CommandArgs{
			PrintVersion: viper.GetBool("print-version"),
		},
		Manifest: manifest.CommandArgs{
			Sources: cleanEnvVarSlice(viper.GetStringSlice("manifest-sources")),
			Targets: cleanEnvVarSlice(viper.GetStringSlice("manifest-targets")),
		},
	}

	return buildah.Exec()
}

func pluginTypeSetup() error {
	pluginType := buildversion.Get().PluginType

	switch pluginType {
	case buildversion.PluginTypeDrone:
		viper.SetEnvPrefix("PLUGIN")
	default:
		return fmt.Errorf("invalid plugin type: %q", pluginType)
	}

	return nil
}

// pflag StringSlice() delimits strings by comma but viper GetStringSlice() does not.
// https://github.com/spf13/viper/issues/380
// This causes issues for environment variables, example: ENV=test1,test2.
// Delimit the input by commas and return the result.
func cleanEnvVarSlice(stringSlice []string) []string {
	var result []string

	for _, s := range stringSlice {
		result = append(result, strings.Split(s, ",")...)
	}

	return result
}
