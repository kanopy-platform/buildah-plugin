package manifest

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/kanopy-platform/buildah-plugin/pkg/buildah/common"
)

const (
	Command = "manifest"
)

type (
	CommandArgs struct {
		Sources []string `json:"sources"`
		Targets []string `json:"targets"`
	}
)

func (c *CommandArgs) GetCmds() ([]*exec.Cmd, error) {
	cmds := []*exec.Cmd{}

	addCmds, err := c.isManifestCmd()
	if err != nil {
		return cmds, err
	}

	if !addCmds {
		return cmds, nil
	}

	// TODO replace with actual manifest commands. Currently is just for testing credentials work.
	cmds = append(cmds,
		exec.Command(common.BuildahCmd, "pull", c.Sources[0]),
	)

	return cmds, nil
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
