package manifest

import (
	"errors"
	"fmt"
	"os/exec"
	"path"

	"github.com/kanopy-platform/buildah-plugin/pkg/buildah/common"
)

const (
	manifestCommand        = "manifest"
	storageDriverOptionVfs = "--storage-driver=vfs"
)

type (
	CommandArgs struct {
		Registry string
		Repo     string
		Sources  []string `json:"sources"`
		Targets  []string `json:"targets"`
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

	return c.manifestCommands(), nil
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

func (c *CommandArgs) manifestCommands() []*exec.Cmd {
	cmds := []*exec.Cmd{}

	registryRepo := path.Join(c.Registry, c.Repo)

	for _, targetTag := range c.Targets {
		target := fmt.Sprintf("%s:%s", registryRepo, targetTag)

		cmds = append(cmds, exec.Command(common.BuildahCmd, manifestCommand, "create", storageDriverOptionVfs, target))

		for _, sourceTag := range c.Sources {
			source := fmt.Sprintf("%s:%s", registryRepo, sourceTag)
			cmds = append(cmds, exec.Command(common.BuildahCmd, manifestCommand, "add", storageDriverOptionVfs, target, source))
		}

		cmds = append(cmds, exec.Command(common.BuildahCmd, manifestCommand, "push", storageDriverOptionVfs, "--all", target))
	}

	return cmds
}
