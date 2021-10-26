package utils

import (
	"os"
	"time"
)

const DateTimeLayout = `2006-01-02 15:04:05`

func Now(zone ...string) time.Time {
	zone = append(zone, GetTimeZone())
	var local, _ = time.LoadLocation(zone[0])
	if local != nil {
		return time.Now().In(local)
	}
	return time.Now()
}

func DateTime(zone ...time.Time) string {
	zone = append(zone, time.Now())
	return TimeLocal(zone[0]).Format(DateTimeLayout)
}

func Timestamp(zone ...string) int64 {
	return Now(zone...).Unix()
}

func TimeLocal(t time.Time) time.Time {
	var local, _ = time.LoadLocation(GetTimeZone())
	if local != nil {
		return t.Local()
	}
	return t.In(local)
}

func GetTimeZone() string {
	var t = os.Getenv("TZ")
	if t == "" {
		return "UTC"
	}
	return t
}
