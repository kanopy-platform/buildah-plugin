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
		Commands []string
		Manifest manifest.CommandArgs
	}

	Login struct {
		Registry string
		Username string
		Password string
	}
)

func (b *Buildah) Exec() error {
	var cmds []*exec.Cmd

	for _, command := range b.Commands {
		switch command {
		case version.Command:
			v := version.CommandArgs{}
			cmds = append(cmds, v.GetCmds()...)
		case manifest.Command:
			cmds = append(cmds, b.Manifest.GetCmds()...)
		default:
			return fmt.Errorf("invalid command: %q", command)
		}
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
