package launchdarkly

import (
	"os"
	"strconv"
	"time"
)

type LaunchDarklyConfiguration struct {
	SdkKey  string
	Timeout time.Duration
	Debug   bool
}

const SdkKeyEnvVar = "FLAGON_LD_SDKKEY"
const TimeoutSecondsEnvVar = "FLAGON_LD_TIMEOUT_SECONDS"
const DebugEnvVar = "FLAGON_LD_DEBUG"

func ConfigFromEnvironment() (LaunchDarklyConfiguration, error) {

	cfg := LaunchDarklyConfiguration{}
	cfg.SdkKey = os.Getenv(SdkKeyEnvVar)
	cfg.Timeout = 10 * time.Second

	if val := os.Getenv(TimeoutSecondsEnvVar); val != "" {
		i, err := strconv.Atoi(val)
		if err != nil {
			return cfg, err
		}
		cfg.Timeout = time.Second * time.Duration(i)
	}

	if val := os.Getenv(DebugEnvVar); val != "" {
		b, err := strconv.ParseBool(val)
		cfg.Debug = err == nil && b
	}

	return cfg, nil
}
