package cli

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalIfExists(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		Age        int    `json:"age"`
		Occupation string `json:"occupation"`
	}

	tests := map[string]struct {
		input   string
		want    testStruct
		wantErr bool
	}{
		"successfully populate struct": {
			input: "{\"age\":40,\"occupation\":\"chef\"}",
			want: testStruct{
				Age:        40,
				Occupation: "chef",
			},
			wantErr: false,
		},
		"input is empty": {
			input:   "",
			want:    testStruct{},
			wantErr: false,
		},
	}

	for name, test := range tests {
		t.Log(name)

		s := testStruct{}
		viper.Set("testKey", test.input)

		err := unmarshalIfExists("testKey", &s)
		t.Log(err)
		assert.Equal(t, test.wantErr, err != nil)
		assert.Equal(t, test.want, s)
	}
}
