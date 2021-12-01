package utils

import (
	"time"
)

const (
	DateFormat = "2006-01-02"
)

func GetDateStringNow() string {
	return time.Now().Format(DateFormat)
}

func IsFutureDateString(dateString string) bool {
	now := time.Now()
	t, err := time.Parse(DateFormat, dateString)
	if err != nil {
		return false
	}
	return !now.After(t)
}
