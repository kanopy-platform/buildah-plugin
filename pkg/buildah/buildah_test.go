package buildah

import (
	"io"
	"os"
	"testing"

	"github.com/kanopy-platform/buildah-plugin/pkg/buildah/manifest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// captureStdout runs fn with os.Stdout redirected to a pipe and returns
// whatever was written to it.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()

	r, w, err := os.Pipe()
	require.NoError(t, err)

	orig := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = orig }()

	fn()

	require.NoError(t, w.Close())

	out, err := io.ReadAll(r)
	require.NoError(t, err)

	return string(out)
}

func TestReportDigests(t *testing.T) {
	// not run in parallel: captureStdout mutates the package-global os.Stdout
	digestFile, err := os.CreateTemp("", "buildah-plugin-test-digest-")
	require.NoError(t, err)
	defer os.Remove(digestFile.Name())

	_, err = digestFile.WriteString("sha256:deadbeef\n")
	require.NoError(t, err)
	require.NoError(t, digestFile.Close())

	digests := []manifest.TargetDigest{
		{Tag: "public.ecr.aws/devops:1", Repo: "public.ecr.aws/devops", DigestFile: digestFile.Name()},
	}

	var reportErr error
	out := captureStdout(t, func() {
		reportErr = reportDigests(digests)
	})

	assert.NoError(t, reportErr)

	// the tag and digest must each be emitted as their own valid, full
	// image reference rather than combined into a single "tag@digest"
	// string: "repo:tag" for the tag reference, "repo@digest" for the
	// digest reference. The repo portion may be sourced from a redacted
	// Drone secret, so it must not be mixed into a single compound value.
	assert.Contains(t, out, "tag: public.ecr.aws/devops:1")
	assert.Contains(t, out, "digest: public.ecr.aws/devops@sha256:deadbeef")
}

func TestReportDigestsMultipleTargets(t *testing.T) {
	// not run in parallel: captureStdout mutates the package-global os.Stdout
	makeDigestFile := func(digest string) string {
		f, err := os.CreateTemp("", "buildah-plugin-test-digest-")
		require.NoError(t, err)
		defer f.Close()

		_, err = f.WriteString(digest + "\n")
		require.NoError(t, err)

		t.Cleanup(func() { os.Remove(f.Name()) })

		return f.Name()
	}

	digests := []manifest.TargetDigest{
		{Tag: "public.ecr.aws/devops:1", Repo: "public.ecr.aws/devops", DigestFile: makeDigestFile("sha256:aaaa")},
		{Tag: "public.ecr.aws/devops:2", Repo: "public.ecr.aws/devops", DigestFile: makeDigestFile("sha256:bbbb")},
	}

	var reportErr error
	out := captureStdout(t, func() {
		reportErr = reportDigests(digests)
	})

	require.NoError(t, reportErr)
	assert.Contains(t, out, "tag: public.ecr.aws/devops:1")
	assert.Contains(t, out, "digest: public.ecr.aws/devops@sha256:aaaa")
	assert.Contains(t, out, "tag: public.ecr.aws/devops:2")
	assert.Contains(t, out, "digest: public.ecr.aws/devops@sha256:bbbb")
}

func TestReportDigestsMissingFile(t *testing.T) {
	t.Parallel()

	digests := []manifest.TargetDigest{
		{Tag: "public.ecr.aws/devops:1", DigestFile: "/nonexistent/digestfile"},
	}

	err := reportDigests(digests)
	assert.Error(t, err)
}
