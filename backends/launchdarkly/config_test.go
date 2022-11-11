package launchdarkly

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReadEnvironment(t *testing.T) {

	os.Setenv(SdkKeyEnvVar, "test-key")
	os.Setenv(TimeoutEnvVar, "17s")
	os.Setenv(DebugEnvVar, "true")

	cfg := ConfigFromEnvironment()

	assert.Equal(t, "test-key", cfg.SdkKey)
	assert.Equal(t, 17*time.Second, cfg.Timeout)
	assert.Equal(t, true, cfg.Debug)
}

func TestFlags(t *testing.T) {

	cfg := LaunchDarklyConfiguration{}
	flags := cfg.Flags()

	assert.NoError(t, flags.Parse([]string{
		"--ld-debug",
		"--ld-sdk-key", "some-key",
		"--ld-timeout", "23s",
	}))

	assert.Equal(t, "some-key", cfg.SdkKey)
	assert.Equal(t, 23*time.Second, cfg.Timeout)
	assert.Equal(t, true, cfg.Debug)
}

func TestOverridingValues(t *testing.T) {

	cases := []struct {
		Override LaunchDarklyConfiguration
		Expected LaunchDarklyConfiguration
	}{
		{
			Override: LaunchDarklyConfiguration{},
			Expected: LaunchDarklyConfiguration{
				SdkKey:  "base-key",
				Timeout: 10 * time.Second,
				Debug:   false,
			},
		},

		{
			Override: LaunchDarklyConfiguration{
				SdkKey: "override-key",
			},
			Expected: LaunchDarklyConfiguration{
				SdkKey:  "override-key",
				Timeout: 10 * time.Second,
				Debug:   false,
			},
		},

		{
			Override: LaunchDarklyConfiguration{
				Debug: true,
			},
			Expected: LaunchDarklyConfiguration{
				SdkKey:  "base-key",
				Timeout: 10 * time.Second,
				Debug:   true,
			},
		},

		{
			Override: LaunchDarklyConfiguration{
				Timeout: 5 * time.Second,
			},
			Expected: LaunchDarklyConfiguration{
				SdkKey:  "base-key",
				Timeout: 5 * time.Second,
				Debug:   false,
			},
		},
	}

	for _, tc := range cases {
		t.Run("", func(t *testing.T) {
			base := LaunchDarklyConfiguration{
				SdkKey:  "base-key",
				Timeout: 10 * time.Second,
				Debug:   false,
			}
			base.OverrideFrom(tc.Override)

			assert.Equal(t, tc.Expected, base)
		})

	}

}
