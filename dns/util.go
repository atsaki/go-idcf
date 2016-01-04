package dns

import (
	"encoding/json"
	"errors"
	"time"
)

func unmarshal(data []byte, v interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(string(data))
		}
	}()
	err = json.Unmarshal(data, &v)
	return
}

func parseTime(s string) (time.Time, error) {
	layout := "2006-01-02T15:04:05-07:00"
	return time.Parse(layout, s)
}
