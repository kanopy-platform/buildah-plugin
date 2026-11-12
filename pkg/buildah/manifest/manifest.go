package manifest

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/kanopy-platform/buildah-plugin/pkg/buildah/common"
)

const (
	manifestCommand        = "manifest"
	storageDriverOptionVfs = "--storage-driver=vfs"
	digestFileFlag         = "--digestfile"
	digestFileNamePattern  = "buildah-plugin-digest-"
)

type (
	CommandArgs struct {
		Registry string
		Repo     string
		Sources  []string `json:"sources"`
		Targets  []string `json:"targets"`

		// Workspace is the directory digest files are created in. It
		// should be set to the plugin's workspace directory (e.g. Drone's
		// workspace mount) rather than left to default to the OS temp
		// directory, since the plugin's runtime environment is not
		// guaranteed to have a writable system tmp dir. If unset, it
		// falls back to the current working directory.
		Workspace string

		digests []TargetDigest
	}

	// TargetDigest associates a pushed manifest list's target tag with the
	// file buildah was asked to write that manifest's digest to.
	TargetDigest struct {
		// Tag is the full "repo:tag" image reference that was pushed.
		Tag string
		// Repo is the image reference without a tag or digest
		// (i.e. "registry/repo"), used to build a valid "repo@digest"
		// reference once the digest is known.
		Repo       string
		DigestFile string
	}
)

func (c *CommandArgs) GetCmds() ([]*exec.Cmd, error) {
	addCmds, err := c.isManifestCmd()
	if err != nil {
		return nil, err
	}

	if !addCmds {
		return nil, nil
	}

	return c.manifestCommands()
}

// Digests returns the target tag/digest-file pairs recorded while building
// the "manifest push" commands. It is only populated after GetCmds has been
// called, and the digest files it references are only populated once the
// corresponding push command has run successfully.
func (c *CommandArgs) Digests() []TargetDigest {
	return c.digests
}

// Cleanup removes any temporary digest files created while building push
// commands. It should be called once the plugin is done consuming the
// results of Digests().
func (c *CommandArgs) Cleanup() {
	for _, d := range c.digests {
		os.Remove(d.DigestFile)
	}
}

func (c *CommandArgs) isManifestCmd() (bool, error) {
	if c.Registry == "" &&
		c.Repo == "" &&
		len(c.Sources) == 0 &&
		len(c.Targets) == 0 {
		return false, nil
	}

	var err error

	if c.Registry == "" {
		err = errors.Join(err, fmt.Errorf("manifest command: registry must be specified"))
	}

	if c.Repo == "" {
		err = errors.Join(err, fmt.Errorf("manifest command: repo must be specified"))
	}

	if len(c.Sources) == 0 {
		err = errors.Join(err, fmt.Errorf("manifest command: sources must be specified"))
	}

	if len(c.Targets) == 0 {
		err = errors.Join(err, fmt.Errorf("manifest command: targets must be specified"))
	}

	return true, err
}

func (c *CommandArgs) manifestCommands() ([]*exec.Cmd, error) {
	cmds := []*exec.Cmd{}

	registryRepo := path.Join(c.Registry, c.Repo)

	digestDir, err := c.digestDir()
	if err != nil {
		return nil, fmt.Errorf("manifest command: determining digest file directory: %w", err)
	}

	for _, targetTag := range c.Targets {
		target := fmt.Sprintf("%s:%s", registryRepo, targetTag)

		cmds = append(cmds, exec.Command(common.BuildahCmd, manifestCommand, "create", storageDriverOptionVfs, target))

		for _, sourceTag := range c.Sources {
			source := fmt.Sprintf("%s:%s", registryRepo, sourceTag)
			cmds = append(cmds, exec.Command(common.BuildahCmd, manifestCommand, "add", storageDriverOptionVfs, target, source))
		}

		digestFile, err := os.CreateTemp(digestDir, digestFileNamePattern)
		if err != nil {
			return nil, fmt.Errorf("manifest command: creating digest file for %q: %w", target, err)
		}

		if err := digestFile.Close(); err != nil {
			return nil, fmt.Errorf("manifest command: closing digest file for %q: %w", target, err)
		}

		c.digests = append(c.digests, TargetDigest{Tag: target, Repo: registryRepo, DigestFile: digestFile.Name()})

		cmds = append(cmds, exec.Command(common.BuildahCmd, manifestCommand, "push", storageDriverOptionVfs, "--all", digestFileFlag, digestFile.Name(), target))
	}

	return cmds, nil
}

// digestDir returns the directory digest files should be created in: the
// configured Workspace, or the current working directory if unset.
func (c *CommandArgs) digestDir() (string, error) {
	if c.Workspace != "" {
		return c.Workspace, nil
	}

	return os.Getwd()
}
