package version

import (
	_ "embed"
	"regexp"
	"strings"
	"time"
)

//go:generate cp ../changelog.md ./changelog.md
//go:embed changelog.md
var changelog string
var sectionRegex = regexp.MustCompile(`## \[(?P<version>.*)\]\s*-\s*(?P<date>.*)`)

var versionGroup = sectionRegex.SubexpIndex("version")
var dateGroup = sectionRegex.SubexpIndex("date")

func Changelog() []ChangeLogEntry {
	return process(changelog)
}

func process(content string) []ChangeLogEntry {
	versions := sectionRegex.FindAllStringIndex(content, -1)
	log := make([]ChangeLogEntry, len(versions))

	versions = append(versions, []int{len(content)})
	for i := 0; i < len(versions)-1; i++ {
		headerStart := versions[i][0]
		logStart := versions[i][1]
		logFinish := versions[i+1][0]

		header := sectionRegex.FindStringSubmatch(content[headerStart:logStart])
		version := header[versionGroup]
		date := header[dateGroup]

		when, _ := time.Parse("2006-01-02", date)

		log[i] = ChangeLogEntry{
			Version: version,
			When:    when,
			Log:     strings.TrimSpace(content[logStart:logFinish]),
		}
	}

	return log
}

type ChangeLogEntry struct {
	Version string
	When    time.Time

	Log string
}
