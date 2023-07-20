package version

import (
	"os/exec"

	"github.com/kanopy-platform/buildah-plugin/pkg/buildah/common"
)

const (
	Command = "version"
)

type (
	CommandArgs struct{}
)

func (c *CommandArgs) GetCmds() []*exec.Cmd {
	return []*exec.Cmd{
		exec.Command(common.BuildahCmd, Command),
	}
}
