package utils

import (
	"encoding/json"
	"strings"
)

func UnescapeUTF8(inStr string) (outStr string, err error) {
	jsonStr := `"` + strings.ReplaceAll(inStr, `"`, `\"`) + `"`
	err = json.Unmarshal([]byte(jsonStr), &outStr)
	return
}
