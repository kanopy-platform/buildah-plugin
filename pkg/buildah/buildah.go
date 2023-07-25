package buildah

import (
	"os"
	"os/exec"
	"strings"

	"github.com/kanopy-platform/buildah-plugin/pkg/buildah/manifest"
	"github.com/kanopy-platform/buildah-plugin/pkg/buildah/version"
	log "github.com/sirupsen/logrus"
)

type (
	Buildah struct {
		Repo     string
		Version  version.CommandArgs
		Manifest manifest.CommandArgs
	}
)

func (b *Buildah) Exec() error {
	var cmds []*exec.Cmd

	cmds = append(cmds, b.Version.GetCmds()...)

	manifestCmds, err := b.Manifest.GetCmds()
	if err != nil {
		return err
	}
	cmds = append(cmds, manifestCmds...)

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
