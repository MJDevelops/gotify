package jsons

import (
	"encoding/json"
)

func ParseJSON(bytes []byte, jsonMap *map[string]string) error {
	err := json.Unmarshal(bytes, jsonMap)
	return err
}
