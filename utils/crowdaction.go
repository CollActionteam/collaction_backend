package utils

import (
	"time"
)

func IsFutureDateString(dateFormat string, dateString string) bool {
	now := time.Now()
	t, err := time.Parse(dateFormat, dateString)
	if err != nil {
		return false
	}
	return now.Before(t)
}
