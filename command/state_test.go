package command

import (
	"context"
	"encoding/json"
	"flagon/backends"
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
			cmd := &StateCommand{}
			cmd.Meta = NewMeta(ui, cmd)
			cmd.Meta.testBackend = &MockBackend{flags: tc.flagStates}

			assert.Equal(t, tc.expectedExit, cmd.Run(args))

			flag := backends.Flag{}
			assert.NoError(t, json.Unmarshal(ui.OutputWriter.Bytes(), &flag))

			assert.Equal(t, tc.expectedFlag, flag)
		})
	}
}

func TestErrorsExitWithStatus2(t *testing.T) {

	t.Run("bad boolean parse", func(t *testing.T) {

		ui := cli.NewMockUi()
		cmd := &StateCommand{}
		cmd.Meta = NewMeta(ui, cmd)
		cmd.Meta.testBackend = &MockBackend{flags: map[string]bool{}}

		assert.Equal(t, 2, cmd.Run([]string{"test-flag", "bad-bool"}))
		assert.Contains(t, ui.ErrorWriter.String(), "parsing \"bad-bool\": invalid syntax")
	})

	t.Run("backend fails", func(t *testing.T) {

		ui := cli.NewMockUi()
		cmd := &StateCommand{}
		cmd.Meta = NewMeta(ui, cmd)

		assert.Equal(t, 2, cmd.Run([]string{"test-flag", "false", "--backend", "non-existing-backend"}))
		assert.Equal(t, "unsupported backend: non-existing-backend\n", ui.ErrorWriter.String())
	})
}

type MockBackend struct {
	flags map[string]bool
}

func (m *MockBackend) State(ctx context.Context, flag backends.Flag, user backends.User) (backends.Flag, error) {
	flag.Value = flag.DefaultValue

	if v, found := m.flags[flag.Key]; found {
		flag.Value = v
	}

	return flag, nil
}

func (m *MockBackend) Close(ctx context.Context) error {
	return nil
}
