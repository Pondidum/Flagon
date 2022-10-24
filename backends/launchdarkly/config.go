package launchdarkly

import (
	"os"
	"strconv"
	"time"
)

type LaunchDarklyConfiguration struct {
	ApiKey  string
	Timeout time.Duration
}

const ApiKeyEnvVar = "FLAGON_LD_APIKEY"
const TimeoutSecondsEnvVar = "FLAGON_LD_TIMEOUT_SECONDS"

func ConfigFromEnvironment() (LaunchDarklyConfiguration, error) {

	cfg := LaunchDarklyConfiguration{}
	cfg.ApiKey = os.Getenv(ApiKeyEnvVar)
	cfg.Timeout = 10 * time.Second

	if val := os.Getenv(TimeoutSecondsEnvVar); val != "" {
		i, err := strconv.Atoi(val)
		if err != nil {
			return cfg, err
		}
		cfg.Timeout = time.Second * time.Duration(i)
	}

	return cfg, nil
}
