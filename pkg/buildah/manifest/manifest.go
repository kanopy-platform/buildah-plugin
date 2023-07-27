package manifest

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/kanopy-platform/buildah-plugin/pkg/buildah/common"
)

const (
	manifestCommand        = "manifest"
	storageDriverOptionVfs = "--storage-driver=vfs"
)

type (
	CommandArgs struct {
		Sources []string `json:"sources"`
		Targets []string `json:"targets"`
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

	return manifestCommands(c.Sources, c.Targets), nil
}

func (c *CommandArgs) isManifestCmd() (bool, error) {
	if len(c.Sources) == 0 && len(c.Targets) == 0 {
		return false, nil
	}

	var err error

	if len(c.Sources) == 0 {
		err = errors.Join(err, fmt.Errorf("manifest command: sources must be specified"))
	}

	if len(c.Targets) == 0 {
		err = errors.Join(err, fmt.Errorf("manifest command: targets must be specified"))
	}

	return true, err
}

func manifestCommands(sources, targets []string) []*exec.Cmd {
	cmds := []*exec.Cmd{}

	for _, target := range targets {
		cmds = append(cmds, exec.Command(common.BuildahCmd, manifestCommand, "create", storageDriverOptionVfs, target))

		for _, source := range sources {
			cmds = append(cmds, exec.Command(common.BuildahCmd, manifestCommand, "add", storageDriverOptionVfs, target, source))
		}

		cmds = append(cmds, exec.Command(common.BuildahCmd, manifestCommand, "push", storageDriverOptionVfs, "--all", target))
	}

	return cmds
}
