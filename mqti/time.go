package mqti

import (
	"strconv"
	"time"
)

// EndOfTime ...
const EndOfTime string = "9999-12-31T23:59:59"

// ParseEpoch ...
func ParseEpoch(in string) time.Time {
	var err error
	var i int64

	if i, err = strconv.ParseInt(in, 10, 64); err != nil {
		Log.Panic(err)
	}

	return time.Unix(i, 0).UTC()
}

// ParseTime ...
func ParseTime(in string) time.Time {
	var err error
	var t time.Time

	if t, err = time.Parse("2006-01-02T15:04:05", in); err != nil {
		Log.Panic(err)
	}

	return t.UTC()
}
