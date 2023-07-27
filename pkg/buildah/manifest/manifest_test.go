package manifest

import (
	"os/exec"
	"testing"

	"github.com/kanopy-platform/buildah-plugin/pkg/buildah/common"
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

func TestManifestCommands(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		sources  []string
		targets  []string
		wantCmds []*exec.Cmd
	}{
		"multiple sources, multiple targets": {
			sources: []string{"a", "b"},
			targets: []string{"1", "2"},
			wantCmds: []*exec.Cmd{
				// target 1
				exec.Command(common.BuildahCmd, manifestCommand, "create", storageDriverOptionVfs, "1"),
				exec.Command(common.BuildahCmd, manifestCommand, "add", storageDriverOptionVfs, "1", "a"),
				exec.Command(common.BuildahCmd, manifestCommand, "add", storageDriverOptionVfs, "1", "b"),
				exec.Command(common.BuildahCmd, manifestCommand, "push", storageDriverOptionVfs, "--all", "1"),
				// target 2
				exec.Command(common.BuildahCmd, manifestCommand, "create", storageDriverOptionVfs, "2"),
				exec.Command(common.BuildahCmd, manifestCommand, "add", storageDriverOptionVfs, "2", "a"),
				exec.Command(common.BuildahCmd, manifestCommand, "add", storageDriverOptionVfs, "2", "b"),
				exec.Command(common.BuildahCmd, manifestCommand, "push", storageDriverOptionVfs, "--all", "2"),
			},
		},
	}

	for name, test := range tests {
		t.Log(name)

		result := manifestCommands(test.sources, test.targets)
		assert.Equal(t, test.wantCmds, result)
	}
}
