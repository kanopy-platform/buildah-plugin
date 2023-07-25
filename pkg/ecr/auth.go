package ecr

import (
	"errors"
	"fmt"
	"os"

	"github.com/kanopy-platform/buildah-plugin/pkg/docker"
)

const (
	accessKeyEnv string = "AWS_ACCESS_KEY_ID"
	secretKeyEnv string = "AWS_SECRET_ACCESS_KEY"
)

func CreateDockerConfig(accessKey, secretKey, registry string) (*docker.Config, error) {
	var errs []error

	if accessKey == "" {
		errs = append(errs, fmt.Errorf("access_key must be specified"))
	}

	if secretKey == "" {
		errs = append(errs, fmt.Errorf("secret_key must be specified"))
	}

	if registry == "" {
		errs = append(errs, fmt.Errorf("registry must be specified"))
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	dockerConfig := docker.NewConfig()

	err := os.Setenv(accessKeyEnv, accessKey)
	if err != nil {
		return nil, fmt.Errorf("failed to set %s environment variable: %v", accessKeyEnv, err)
	}

	err = os.Setenv(secretKeyEnv, secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to set %s environment variable: %v", secretKeyEnv, err)
	}

	// uses the amazon-ecr-credential-helper credential helper
	dockerConfig.SetCredHelper(registry, "ecr-login")

	return dockerConfig, nil
}
