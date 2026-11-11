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
	defer b.Manifest.Cleanup()

	cmds = append(cmds, manifestCmds...)

	for _, cmd := range cmds {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		log.Infof("running: %s", strings.Join(cmd.Args, " "))

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return reportDigests(b.Manifest.Digests())
}

// reportDigests prints a tag reference and a digest reference for each
// pushed target, as independent, individually valid, pullable image
// references:
//
//	tag: repo:tag
//	digest: repo@sha256:...
//
// They are kept separate rather than combined into a single compound
// string since the repo portion may be sourced from a Drone secret and
// therefore partially redacted in plugin output. It is expected to be
// called only after the corresponding "manifest push" commands have run
// successfully, since that is when buildah writes the digest to the
// associated digest file.
func reportDigests(digests []manifest.TargetDigest) error {
	for _, d := range digests {
		data, err := os.ReadFile(d.DigestFile)
		if err != nil {
			return fmt.Errorf("reading digest file for %q: %w", d.Tag, err)
		}

		digest := strings.TrimSpace(string(data))

		fmt.Printf("tag: %s\n", d.Tag)
		fmt.Printf("digest: %s@%s\n", d.Repo, digest)
	}

	return nil
}
