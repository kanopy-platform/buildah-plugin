package buildah

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/kanopy-platform/buildah-plugin/pkg/buildah/manifest"
	"github.com/kanopy-platform/buildah-plugin/pkg/buildah/version"
	log "github.com/sirupsen/logrus"
)

type (
	Buildah struct {
		Login    Login // configuration for "buildah login"
		Repo     string
		Command  string
		Manifest manifest.CommandArgs
	}

	Login struct {
		Registry string
		Username string
		Password string
	}
)

func (b *Buildah) Exec() error {
	log.Infof("Buildah struct: %+v", b) // TODO remove

	var cmds []*exec.Cmd

	switch b.Command {
	case version.Command:
		v := version.CommandArgs{}
		cmds = append(cmds, v.GetCmds()...)
	case manifest.Command:
		cmds = append(cmds, b.Manifest.GetCmds()...)
	default:
		return fmt.Errorf("invalid command: %q", b.Command)
	}

	for _, cmd := range cmds {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		log.Infof("running: %s", strings.Join(cmd.Args, " "))

		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
