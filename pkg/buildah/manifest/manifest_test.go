package manifest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	workspace := t.TempDir()

	tests := map[string]struct {
		commandArgs CommandArgs
		wantRepo    string
		wantTags    []string
	}{
		"multiple sources, multiple targets": {
			commandArgs: CommandArgs{
				Registry:  "public.ecr.aws",
				Repo:      "devops",
				Sources:   []string{"a", "b"},
				Targets:   []string{"1", "2"},
				Workspace: workspace,
			},
			wantRepo: "public.ecr.aws/devops",
			wantTags: []string{
				"public.ecr.aws/devops:1",
				"public.ecr.aws/devops:2",
			},
		},
	}

	for name, test := range tests {
		t.Log(name)

		args := test.commandArgs
		cmds, err := args.manifestCommands()
		require.NoError(t, err)
		defer args.Cleanup()

		// 1 create + 2 add + 1 push per target
		wantLen := len(test.wantTags) * (1 + len(args.Sources) + 1)
		assert.Len(t, cmds, wantLen)

		digests := args.Digests()
		require.Len(t, digests, len(test.wantTags))

		cmdsPerTarget := 1 + len(args.Sources) + 1

		for i, tag := range test.wantTags {
			assert.Equal(t, tag, digests[i].Tag)
			// Repo is the image reference without a tag or digest, used to
			// build a valid "repo@digest" reference for the pushed digest.
			assert.Equal(t, test.wantRepo, digests[i].Repo)
			assert.NotEmpty(t, digests[i].DigestFile)

			// the digest file must be created inside the configured
			// workspace directory rather than the OS temp directory, since
			// the plugin's runtime environment is not guaranteed to have a
			// writable system tmp dir.
			assert.Equal(t, workspace, filepath.Dir(digests[i].DigestFile))

			// the digest file should have been created on disk
			_, err := os.Stat(digests[i].DigestFile)
			assert.NoError(t, err)

			// the push command for this target should reference the same
			// digest file via --digestfile
			pushCmd := cmds[(i+1)*cmdsPerTarget-1]
			assert.Contains(t, pushCmd.Args, "--digestfile")
			assert.Contains(t, pushCmd.Args, digests[i].DigestFile)
			assert.Contains(t, pushCmd.Args, tag)
		}
	}
}

func TestManifestCommandsDefaultsToWorkingDirectory(t *testing.T) {
	// not run in parallel: changes the process-wide working directory
	workspace := t.TempDir()

	origWd, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(workspace))
	defer func() { require.NoError(t, os.Chdir(origWd)) }()

	args := CommandArgs{
		Registry: "public.ecr.aws",
		Repo:     "devops",
		Sources:  []string{"a"},
		Targets:  []string{"1"},
		// Workspace intentionally left unset
	}

	_, err = args.manifestCommands()
	require.NoError(t, err)
	defer args.Cleanup()

	digests := args.Digests()
	require.Len(t, digests, 1)

	// with no Workspace configured, the digest file must still land in the
	// current working directory rather than the OS temp directory.
	// filepath.EvalSymlinks normalizes any symlinked tmp dirs (e.g. macOS's
	// /var -> /private/var) so the comparison isn't platform-dependent.
	wantDir, err := filepath.EvalSymlinks(workspace)
	require.NoError(t, err)
	gotDir, err := filepath.EvalSymlinks(filepath.Dir(digests[0].DigestFile))
	require.NoError(t, err)
	assert.Equal(t, wantDir, gotDir)
}

func TestDigestsCleanup(t *testing.T) {
	t.Parallel()

	args := CommandArgs{
		Registry:  "public.ecr.aws",
		Repo:      "devops",
		Sources:   []string{"a"},
		Targets:   []string{"1"},
		Workspace: t.TempDir(),
	}

	_, err := args.manifestCommands()
	require.NoError(t, err)

	digests := args.Digests()
	require.Len(t, digests, 1)

	// digest file exists prior to cleanup
	_, err = os.Stat(digests[0].DigestFile)
	require.NoError(t, err)

	args.Cleanup()

	// digest file is removed after cleanup
	_, err = os.Stat(digests[0].DigestFile)
	assert.True(t, os.IsNotExist(err))
}
