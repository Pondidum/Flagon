package launchdarkly

import (
	"os"
	"strconv"
	"time"

	"github.com/spf13/pflag"
)

const SdkKeyEnvVar = "FLAGON_LD_SDKKEY"
const TimeoutEnvVar = "FLAGON_LD_TIMEOUT"
const DebugEnvVar = "FLAGON_LD_DEBUG"

type LaunchDarklyConfiguration struct {
	SdkKey  string
	Timeout time.Duration
	Debug   bool
}

func (cfg *LaunchDarklyConfiguration) OverrideFrom(other LaunchDarklyConfiguration) {
	if other.Debug {
		cfg.Debug = other.Debug
	}

	if other.SdkKey != "" {
		cfg.SdkKey = other.SdkKey
	}

	if other.Timeout > 0 {
		cfg.Timeout = other.Timeout
	}
}

func (cfg *LaunchDarklyConfiguration) Flags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("LaunchDarkly Backend", pflag.ContinueOnError)

	flags.BoolVar(&cfg.Debug, "ld-debug", false, "enable debug logging for launchdarkly")
	flags.StringVar(&cfg.SdkKey, "ld-sdk-key", "", "the sdk-key to use")
	flags.DurationVar(&cfg.Timeout, "ld-timeout", 0, "timeout before failing to communicate with launchdarkly")

	return flags
}

func ConfigFromEnvironment() LaunchDarklyConfiguration {

	cfg := LaunchDarklyConfiguration{}
	cfg.SdkKey = os.Getenv(SdkKeyEnvVar)

	if val := os.Getenv(TimeoutEnvVar); val != "" {
		if timeout, err := time.ParseDuration(val); err == nil {
			cfg.Timeout = timeout
		}

	}

	if val := os.Getenv(DebugEnvVar); val != "" {
		b, err := strconv.ParseBool(val)
		cfg.Debug = err == nil && b
	}

	return cfg
}

func DefaultConfig() LaunchDarklyConfiguration {
	return LaunchDarklyConfiguration{
		SdkKey:  "",
		Timeout: 2 * time.Second,
		Debug:   false,
	}
}
