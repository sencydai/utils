package systime

import (
	"time"
)

const (
	DATETIME_FORMAT = "2006-01-02 15:04:05"
	DATE_FORMAT     = "2006-01-02"
	TIME_FORMAT     = "15:04:05"
)

func Format(t time.Time, layout string) string {
	return t.Format(layout)
}

func FormatDateTime(t time.Time) string {
	return t.Format(DATETIME_FORMAT)
}

func FormatDate(t time.Time) string {
	return t.Format(DATE_FORMAT)
}

func FormatTime(t time.Time) string {
	return t.Format(TIME_FORMAT)
}

func Parse(layout string, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, time.Local)
}

func ParseDateTime(value string) (time.Time, error) {
	return time.ParseInLocation(DATETIME_FORMAT, value, time.Local)
}

func ParseDate(value string) (time.Time, error) {
	return time.ParseInLocation(DATE_FORMAT, value, time.Local)
}

func ParseTime(value string) (time.Time, error) {
	return time.ParseInLocation(TIME_FORMAT, value, time.Local)
}
