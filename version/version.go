package version

import (
	"fmt"
)

var (
	GitCommit  string
	Prerelease = "dev"
)

func VersionNumber() string {
	if GitCommit == "" {
		return "local"
	}

	version := GitCommit[0:7]

	if Prerelease != "" {
		version = fmt.Sprintf("%s-%s", version, Prerelease)
	}

	return version
}
