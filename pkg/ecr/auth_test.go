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
