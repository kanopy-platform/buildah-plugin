package version

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kanopy-platform/buildah-plugin/pkg/buildah/common"
	"github.com/stretchr/testify/assert"
)

func TestGetCmds(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		args     CommandArgs
		wantCmds []string
	}{
		"printVersion=false": {
			args: CommandArgs{
				Print: false,
			},
			wantCmds: []string{},
		},
		"printVersion=true": {
			args: CommandArgs{
				Print: true,
			},
			wantCmds: []string{fmt.Sprintf("%s %s", common.BuildahCmd, versionCommand)},
		},
	}

	for name, test := range tests {
		t.Log(name)

		result := test.args.GetCmds()
		for index, wantCmd := range test.wantCmds {
			assert.Equal(t, wantCmd, strings.Join(result[index].Args, " "))
		}
	}
}
