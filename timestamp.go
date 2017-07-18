package main

import (
	"fmt"
	"strconv"
	"time"
)

// Timestamp is defined here as per
// https://gist.github.com/bsphere/8369aca6dde3e7b4392c#gistcomment-1413740
type Timestamp struct {
	time.Time
}

// MarshalJSON time2json
func (t *Timestamp) MarshalJSON() ([]byte, error) {
	ts := t.Time.Unix()
	stamp := fmt.Sprint(ts)
	return []byte(stamp), nil
}

// UnmarshalJSON json2time
func (t *Timestamp) UnmarshalJSON(b []byte) error {
	ts, err := strconv.Atoi(string(b))
	if err != nil {
		return err
	}
	t.Time = time.Unix(int64(ts), 0)
	return nil
}
