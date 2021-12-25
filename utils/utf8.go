package utils

import (
	"encoding/json"
	"strings"
)

func UnescapeUTF8(inStr string) (outStr string, err error) {
	// Use json.Unmarshal to replace e.g. "\u0026" with "&"
	jsonStr := `"` + strings.ReplaceAll(inStr, `"`, `\"`) + `"`
	err = json.Unmarshal([]byte(jsonStr), &outStr)
	return
}
