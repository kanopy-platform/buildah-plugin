package manifest

import (
	"errors"
	"fmt"
	"os/exec"
)

const (
	Command = "manifest"
)

type (
	CommandArgs struct {
		Sources []string
		Targets []string
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

	// TODO add commands to run

	return cmds, nil
}

func (c *CommandArgs) isManifestCmd() (bool, error) {
	if len(c.Sources) == 0 && len(c.Targets) == 0 {
		return false, nil
	}

	var err error

	if len(c.Sources) == 0 {
		err = errors.Join(err, fmt.Errorf("manifest command: manifest_sources must be specified"))
	}

	if len(c.Targets) == 0 {
		err = errors.Join(err, fmt.Errorf("manifest command: manifest_targets must be specified"))
	}

	return true, err
}
