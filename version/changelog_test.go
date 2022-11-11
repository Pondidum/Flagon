package version

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestChangelog(t *testing.T) {

	first := ChangeLogEntry{
		Version: "0.0.1",
		When:    time.Date(2022, 11, 02, 0, 0, 0, 0, time.UTC),
		Log: `### Added

- Support S3 backend
- Read and Write metadata
- Fetch and Store artifacts`,
	}

	second := ChangeLogEntry{
		When:    time.Date(2022, 10, 27, 0, 0, 0, 0, time.UTC),
		Version: "0.0.0",
		Log: `### Added

- Initial Version`,
	}

	assert.Equal(t, []ChangeLogEntry{first, second}, process(testChangelog))
}

var testChangelog = `# Changelog

## [0.0.1] - 2022-11-02

### Added

- Support S3 backend
- Read and Write metadata
- Fetch and Store artifacts

## [0.0.0] - 2022-10-27

### Added

- Initial Version
`
