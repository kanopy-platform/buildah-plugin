package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanEnvVarSlice(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input []string
		want  []string
	}{
		"one-element-multiple-values": {
			input: []string{"value1,value2"},
			want:  []string{"value1", "value2"},
		},
		"multiple-elements": {
			input: []string{"value1", "value2,value3"},
			want:  []string{"value1", "value2", "value3"},
		},
	}

	for name, test := range tests {
		t.Log(name)

		result := cleanEnvVarSlice(test.input)
		assert.Equal(t, test.want, result)
	}
}
