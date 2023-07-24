package manifest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsManifestCmd(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		args       CommandArgs
		wantResult bool
		wantErr    bool
	}{
		"not manifest command": {
			args:       CommandArgs{},
			wantResult: false,
			wantErr:    false,
		},
		"all required args exist": {
			args: CommandArgs{
				Sources: []string{"a"},
				Targets: []string{"b"},
			},
			wantResult: true,
			wantErr:    false,
		},
		"missing sources": {
			args: CommandArgs{
				Targets: []string{"b"},
			},
			wantResult: true,
			wantErr:    true,
		},
		"missing targets": {
			args: CommandArgs{
				Sources: []string{"a"},
			},
			wantResult: true,
			wantErr:    true,
		},
	}

	for name, test := range tests {
		t.Log(name)

		result, err := test.args.isManifestCmd()
		assert.Equal(t, test.wantResult, result)
		assert.Equal(t, test.wantErr, err != nil)
	}
}
