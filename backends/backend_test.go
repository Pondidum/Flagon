package backends

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshaling(t *testing.T) {

	f := Flag{
		Key:          "the-key",
		DefaultValue: true,
		Value:        false,
	}

	bytes, _ := json.Marshal(f)
	plain := map[string]interface{}{}
	json.Unmarshal(bytes, &plain)

	assert.Equal(t, f.Key, plain["key"])
	assert.Equal(t, f.DefaultValue, plain["defaultValue"])
	assert.Equal(t, f.Value, plain["value"])
}
