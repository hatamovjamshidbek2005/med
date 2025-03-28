package utils

import (
	"fmt"
	"time"
)

const layoutDateOnly = "2006-01-02"
const layoutDateTime = "2006-01-02 15:04:05"
const defaultTime = "00:00:00"

func DateValidate(dateStr string) (time.Time, error) {
	if len(dateStr) == len(layoutDateOnly) {
		dateStr = fmt.Sprintf("%s %s", dateStr, defaultTime)
	}

	parsedTime, err := time.Parse(layoutDateTime, dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("time parse error")
	}
	return parsedTime, nil
}
func TimeValidate(timeStr string) (time.Time, error) {
	const timeFormat = "15:04:05"

	parsedTime, err := time.Parse(timeFormat, timeStr)
	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
}
