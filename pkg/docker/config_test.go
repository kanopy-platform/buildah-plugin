package docker

import (
	"encoding/json"
	"testing"
)

func TestConfig(t *testing.T) {
	testRegistry := "public.ecr.aws"

	c := NewConfig()

	c.SetAuth(testRegistry, "test", "password")
	c.SetCredHelper(testRegistry, "ecr-login")

	bytes, err := json.Marshal(c)
	if err != nil {
		t.Error("json marshal failed")
	}

	want := `{"auths":{"public.ecr.aws":{"auth":"dGVzdDpwYXNzd29yZA=="}},"credHelpers":{"public.ecr.aws":"ecr-login"}}`
	got := string(bytes)

	if want != got {
		t.Errorf("unexpected json output:\n  want: %s\n   got: %s", want, got)
	}
}
