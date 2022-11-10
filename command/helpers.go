package command

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
)

func parseKeyValuePairs(tags []string) (map[string]string, error) {

	m := map[string]string{}

	for _, pair := range tags {
		index := strings.Index(pair, "=")
		if index == -1 {
			return nil, fmt.Errorf("must be in the format key=value")
		}

		key := strings.TrimSpace(pair[:index])
		val := strings.TrimSpace(pair[index+1:])

		if key == "" {
			return nil, fmt.Errorf("no key specified (must be in the format key=value)")
		}
		if val == "" {
			return nil, fmt.Errorf("no value specified (must be in the format key=value)")
		}

		m[key] = val
	}

	return m, nil
}

func print(ui cli.Ui, format string, vals map[string]interface{}) error {

	switch format {
	case "json":
		b, err := json.Marshal(vals)
		if err != nil {
			return err
		}
		ui.Output(string(b))

	}

	return nil
}
