package version

import (
	"os/exec"

	"github.com/kanopy-platform/buildah-plugin/pkg/buildah/common"
)

const (
	Command = "version"
)

type (
	CommandArgs struct {
		PrintVersion bool
	}
)

func (c *CommandArgs) GetCmds() []*exec.Cmd {
	if c.PrintVersion {
		return []*exec.Cmd{
			exec.Command(common.BuildahCmd, Command),
		}
	}

	return []*exec.Cmd{}
}
