package command

import (
	"context"
	"encoding/json"
	"flagon/backends"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"
)

func TestState(t *testing.T) {

	cases := []struct {
		key          string
		defaultValue string
		flagStates   map[string]bool
		expectedExit int
		expectedFlag backends.Flag
	}{
		{
			key:          "flag-name",
			expectedExit: 1,
			expectedFlag: backends.Flag{Key: "flag-name", DefaultValue: false, Value: false},
		},
		{
			key:          "flag-name",
			defaultValue: "false",
			expectedExit: 1,
			expectedFlag: backends.Flag{Key: "flag-name", DefaultValue: false, Value: false},
		},
		{
			key:          "flag-name",
			defaultValue: "true",
			expectedExit: 0,
			expectedFlag: backends.Flag{Key: "flag-name", DefaultValue: true, Value: true},
		},
		{
			key: "flag-name",
			flagStates: map[string]bool{
				"flag-name": true,
			},
			expectedExit: 0,
			expectedFlag: backends.Flag{Key: "flag-name", DefaultValue: false, Value: true},
		},
	}

	for _, tc := range cases {
		t.Run("", func(t *testing.T) {
			args := []string{tc.key}
			if tc.defaultValue != "" {
				args = append(args, tc.defaultValue)
			}

			ui := cli.NewMockUi()
			cmd, _ := NewStateCommand(ui)
			cmd.Meta.testBackend = &MockBackend{flags: tc.flagStates}

			assert.Equal(t, tc.expectedExit, cmd.Run(args))

			flag := backends.Flag{}
			assert.NoError(t, json.Unmarshal(ui.OutputWriter.Bytes(), &flag))

			assert.Equal(t, tc.expectedFlag, flag)
		})
	}
}

func TestAttributeParsing(t *testing.T) {

	cases := []struct {
		name            string
		userKey         string
		attrs           []string
		attrFile        string
		files           map[string]string
		expectedAttrs   map[string]string
		expectedUserKey string
	}{
		{attrs: []string{}},

		{
			name:  "only cli flags",
			attrs: []string{"from=cli"},
			expectedAttrs: map[string]string{
				"from": "cli",
			},
		},
		{
			name:  "from default flagon attrs file",
			attrs: []string{},
			files: map[string]string{
				"flagon.attrs": "from=default file",
			},
			expectedAttrs: map[string]string{
				"from": "default file",
			},
		},
		{
			name:  "overridden by cli flagon attrs file",
			attrs: []string{"from=cli"},
			files: map[string]string{
				"flagon.attrs": "from=default file",
			},
			expectedAttrs: map[string]string{
				"from": "cli",
			},
		},
		{
			name:     "from specified flagon attrs file",
			attrs:    []string{"from=default file"},
			attrFile: "other.attrs",
			files: map[string]string{
				"flagon.attrs": "from=default file",
				"other.attrs":  "from=specified",
			},
			expectedAttrs: map[string]string{
				"from": "default file",
			},
		},
		{
			name:     "non-existing specified attrs file",
			attrs:    []string{},
			attrFile: "other.attrs",
			files: map[string]string{
				"flagon.attrs": "from=default file",
			},
			expectedAttrs: map[string]string{
				"from": "",
			},
		},
		{
			name:            "user key from cli",
			userKey:         "cli",
			attrs:           []string{"user-key=cli"},
			expectedUserKey: "cli",
			expectedAttrs:   map[string]string{"user-key": ""},
		},
		{
			name:            "user key from cli attr",
			attrs:           []string{"user-key=attr"},
			expectedUserKey: "attr",
			expectedAttrs:   map[string]string{"user-key": ""},
		},
		{
			name:            "user key from cli attr and cli flag",
			userKey:         "flag",
			attrs:           []string{"user-key=attr"},
			expectedUserKey: "flag",
			expectedAttrs:   map[string]string{"user-key": ""},
		},
		{
			name: "user key from attr file",
			files: map[string]string{
				"flagon.attrs": "user-key=file",
			},
			expectedUserKey: "file",
			expectedAttrs:   map[string]string{"user-key": ""},
		},
		{
			name:    "user key from attr file and flag",
			userKey: "flag",
			files: map[string]string{
				"flagon.attrs": "user-key=file",
			},
			expectedUserKey: "flag",
			expectedAttrs:   map[string]string{"user-key": ""},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			args := make([]string, 0, len(tc.attrs)*3)
			args = append(args, "some-flag")

			if tc.userKey != "" {
				args = append(args, "--user", tc.userKey)
			}
			for _, attr := range tc.attrs {
				args = append(args, "--attr", attr)
			}

			if tc.attrFile != "" {
				args = append(args, "--attr-file", tc.attrFile)
			}

			backend := &MockBackend{
				flags: map[string]bool{
					"some-flag": true,
				},
			}

			ui := cli.NewMockUi()
			cmd, _ := NewStateCommand(ui)
			cmd.readFile = func(filePath string) (io.ReadCloser, error) {
				content, found := tc.files[filePath]
				if !found {
					return nil, os.ErrNotExist
				}
				return NewReadCloser(content), nil
			}
			cmd.Meta.testBackend = backend

			assert.Equal(t, 0, cmd.Run(args))

			if tc.expectedUserKey != "" {
				assert.Equal(t, tc.expectedUserKey, backend.users[0].Key)
			}

			for k, v := range tc.expectedAttrs {
				assert.Equal(t, v, backend.users[0].Attributes[k])
			}
		})
	}
}

func TestErrorsExitWithStatus2(t *testing.T) {

	t.Run("bad boolean parse", func(t *testing.T) {

		ui := cli.NewMockUi()
		cmd, _ := NewStateCommand(ui)
		cmd.Meta.testBackend = &MockBackend{flags: map[string]bool{}}

		assert.Equal(t, 2, cmd.Run([]string{"test-flag", "bad-bool"}))
		assert.Contains(t, ui.ErrorWriter.String(), "parsing \"bad-bool\": invalid syntax")
	})

	t.Run("backend fails", func(t *testing.T) {

		ui := cli.NewMockUi()
		cmd, _ := NewStateCommand(ui)

		assert.Equal(t, 2, cmd.Run([]string{"test-flag", "false", "--backend", "non-existing-backend"}))
		assert.Equal(t, "unsupported backend: non-existing-backend\n", ui.ErrorWriter.String())
	})
}

type MockBackend struct {
	flags map[string]bool
	users []backends.User
}

func (m *MockBackend) State(ctx context.Context, flag backends.Flag, user backends.User) (backends.Flag, error) {
	flag.Value = flag.DefaultValue

	if v, found := m.flags[flag.Key]; found {
		flag.Value = v
	}

	m.users = append(m.users, user)

	return flag, nil
}

func (m *MockBackend) Close(ctx context.Context) error {
	return nil
}

func NewReadCloser(content string) io.ReadCloser {
	return &stringrc{wrapped: strings.NewReader(content)}
}

type stringrc struct {
	wrapped io.Reader
}

func (s *stringrc) Read(p []byte) (n int, err error) {
	return s.wrapped.Read(p)
}

func (s *stringrc) Close() error {
	return nil
}
