package command

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseKeyValuePairs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    []string
		expected map[string]string
		err      error
	}{
		{
			input:    []string{"key=value"},
			expected: map[string]string{"key": "value"},
			err:      nil,
		},
		{
			input:    []string{"key=value", "other=yes"},
			expected: map[string]string{"key": "value", "other": "yes"},
			err:      nil,
		},
		{
			input: []string{"key"},
			err:   fmt.Errorf("must be in the format key=value"),
		},
		{
			input: []string{"=value"},
			err:   fmt.Errorf("no key specified (must be in the format key=value)"),
		},
		{
			input: []string{"key="},
			err:   fmt.Errorf("no value specified (must be in the format key=value)"),
		},
		{
			input:    []string{"key=value=test"},
			expected: map[string]string{"key": "value=test"},
			err:      nil,
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {

			m, err := parseKeyValuePairs(tc.input)
			if tc.err == nil {
				assert.NoError(t, err)

				assert.Equal(t, tc.expected, m)
			} else {
				assert.Equal(t, tc.err, err)
			}
		})
	}

}
