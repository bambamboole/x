package pkg

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseArgs(t *testing.T) {
	t.Run("It parses arguments as expected", func(t *testing.T) {
		tests := []struct {
			name           string
			passedArgs     []string
			expectedStruct Arguments
		}{
			{
				"Default values",
				[]string{},
				Arguments{
					Taskfiles:   []string{"~/.x/Taskfile"},
					ConfigFiles: []string{"~/.x/config.yml"},
					Shell:       "bash",
				},
			},
			{
				"additional Taskfile",
				[]string{"-t", "~/foo"},
				Arguments{
					Taskfiles:   []string{"~/foo"},
					ConfigFiles: []string{"~/.x/config.yml"},
					Shell:       "bash",
				},
			},
			{
				"Multiple Taskfiles",
				[]string{"-t", "~/foo", "-t", "~/bar"},
				Arguments{
					Taskfiles:   []string{"~/foo", "~/bar"},
					ConfigFiles: []string{"~/.x/config.yml"},
					Shell:       "bash",
				},
			},
			{
				"Custom shell",
				[]string{"-s", "zsh"},
				Arguments{
					Taskfiles:   []string{"~/.x/Taskfile"},
					ConfigFiles: []string{"~/.x/config.yml"},
					Shell:       "zsh",
				},
			},
		}
		for _, test := range tests {
			args, _ := ParseArgs(test.passedArgs, &bytes.Buffer{})
			assert.Equal(t, test.expectedStruct, args)
		}
	})
	t.Run("It writes help to the passed stdout", func(t *testing.T) {
		tests := []struct {
			name       string
			passedArgs []string
		}{
			{"On -h", []string{"-h"}},
			{"On --help", []string{"--help"}},
			{"On nothing", []string{}},
		}
		for _, test := range tests {
			var b bytes.Buffer
			_, _ = ParseArgs(test.passedArgs, &b)
			assert.True(t, b.Len() > 0)
		}
	})
	t.Run("Left over command arguments", func(t *testing.T) {
		tests := []struct {
			name            string
			passedArgs      []string
			expectedCommand []string
		}{
			{"verbose gets passes through", []string{"test", "-v"}, []string{"test", "-v"}},
			{"verbose gets passes through", []string{"test", "-vv"}, []string{"test", "-vv"}},
			{"verbose gets passes through", []string{"test", "-vvv"}, []string{"test", "-vvv"}},
		}
		for _, test := range tests {
			args, _ := ParseArgs(test.passedArgs, &bytes.Buffer{})
			assert.Equal(t, test.expectedCommand, args.Command)
		}
	})
}
