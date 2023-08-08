package ecr

import (
	"os"
	"testing"

	"github.com/kanopy-platform/buildah-plugin/pkg/docker"
	"github.com/stretchr/testify/assert"
)

func TestCreateDockerConfig(t *testing.T) {
	t.Parallel()

	testDockerConfig := docker.NewConfig()
	testDockerConfig.SetCredHelper("hello.com", "ecr-login")

	publicDockerConfig := docker.NewConfig()
	publicDockerConfig.SetCredHelper("public.ecr.aws", "ecr-login")

	tests := map[string]struct {
		accessKey string
		secretKey string
		registry  string
		want      *docker.Config
		wantErr   bool
	}{
		"missing accessKey, secretKey, registry": {
			wantErr: true,
		},
		"successful": {
			accessKey: "access",
			secretKey: "secret",
			registry:  "hello.com",
			want:      testDockerConfig,
		},
		"clean public registry": {
			accessKey: "access",
			secretKey: "secret",
			registry:  "public.ecr.aws/kanopy",
			want:      publicDockerConfig,
		},
	}

	for name, test := range tests {
		t.Log(name)

		result, err := CreateDockerConfig(test.accessKey, test.secretKey, test.registry)
		assert.Equal(t, test.want, result)
		assert.Equal(t, test.wantErr, err != nil)

		if !test.wantErr {
			assert.Equal(t, test.accessKey, os.Getenv(accessKeyEnv))
			assert.Equal(t, test.secretKey, os.Getenv(secretKeyEnv))
		}
	}
}

func TestCleanRegistryName(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		registry string
		want     string
	}{
		"public registry, remove suffix": {
			registry: "public.ecr.aws/kanopy",
			want:     "public.ecr.aws",
		},
		"public registry, no suffix": {
			registry: "public.ecr.aws",
			want:     "public.ecr.aws",
		},
		"private registry, leave as is": {
			registry: "hello.com/kanopy",
			want:     "hello.com/kanopy",
		},
	}

	for name, test := range tests {
		t.Log(name)

		result := cleanRegistryName(test.registry)
		assert.Equal(t, test.want, result)
	}
}
