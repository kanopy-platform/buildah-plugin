package manifest

import (
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

func (c *CommandArgs) GetCmds() []*exec.Cmd {
	return []*exec.Cmd{}
}
