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
				Registry: "a",
				Repo:     "b",
				Sources:  []string{"a"},
				Targets:  []string{"b"},
			},
			wantResult: true,
			wantErr:    false,
		},
		"missing sources": {
			args: CommandArgs{
				Registry: "a",
				Repo:     "b",
				Targets:  []string{"b"},
			},
			wantResult: true,
			wantErr:    true,
		},
		"missing targets": {
			args: CommandArgs{
				Registry: "a",
				Repo:     "b",
				Sources:  []string{"a"},
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
		commandArgs CommandArgs
		wantCmds    []*exec.Cmd
	}{
		"multiple sources, multiple targets": {
			commandArgs: CommandArgs{
				Registry: "public.ecr.aws",
				Repo:     "devops",
				Sources:  []string{"a", "b"},
				Targets:  []string{"1", "2"},
			},
			wantCmds: []*exec.Cmd{
				// target 1
				exec.Command(common.BuildahCmd, manifestCommand, "create", storageDriverOptionVfs, "public.ecr.aws/devops:1"),
				exec.Command(common.BuildahCmd, manifestCommand, "add", storageDriverOptionVfs, "public.ecr.aws/devops:1", "public.ecr.aws/devops:a"),
				exec.Command(common.BuildahCmd, manifestCommand, "add", storageDriverOptionVfs, "public.ecr.aws/devops:1", "public.ecr.aws/devops:b"),
				exec.Command(common.BuildahCmd, manifestCommand, "push", storageDriverOptionVfs, "--all", "public.ecr.aws/devops:1"),
				// target 2
				exec.Command(common.BuildahCmd, manifestCommand, "create", storageDriverOptionVfs, "public.ecr.aws/devops:2"),
				exec.Command(common.BuildahCmd, manifestCommand, "add", storageDriverOptionVfs, "public.ecr.aws/devops:2", "public.ecr.aws/devops:a"),
				exec.Command(common.BuildahCmd, manifestCommand, "add", storageDriverOptionVfs, "public.ecr.aws/devops:2", "public.ecr.aws/devops:b"),
				exec.Command(common.BuildahCmd, manifestCommand, "push", storageDriverOptionVfs, "--all", "public.ecr.aws/devops:2"),
			},
		},
	}

	for name, test := range tests {
		t.Log(name)

		result := test.commandArgs.manifestCommands()
		assert.Equal(t, test.wantCmds, result)
	}
}
