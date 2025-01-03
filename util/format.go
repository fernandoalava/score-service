package util

import (
	"math"
	"time"
)

const DateTimeDefaultStringFormat = "2006-01-02T15:04:05"

func TimeToString(value time.Time) string {
	return value.Format(DateTimeDefaultStringFormat)
}

func StringToTime(value string) (t time.Time, err error) {
	t, err = time.Parse(DateTimeDefaultStringFormat, value)
	return
}

func StringToTimeWithFormat(value string, format string) (t time.Time, err error) {
	t, err = time.Parse(format, value)
	return
}

func FormatScore(score float64) float64 {
	return math.Round(score*100) / 100
}
