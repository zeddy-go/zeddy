package json

import (
	"strings"
	"time"
)

var (
	TimeFormat = time.DateTime
)

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 2 {
		t.Time = time.Time{}
		return
	}

	t.Time, err = time.ParseInLocation(TimeFormat, strings.Trim(string(data), "\""), time.Local)
	if err != nil {
		return
	}

	return
}

func (t Time) MarshalJSON() (data []byte, err error) {
	return []byte("\"" + t.Time.In(time.Local).Format(TimeFormat) + "\""), nil
}
