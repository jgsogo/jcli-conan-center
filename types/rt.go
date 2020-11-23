package types

import (
	"strings"
	"time"
)

const (
	timeLayout = "2006-01-02T15:04:05.999+0000"
)

type RtTimestamp struct {
	time.Time
}

func (ct *RtTimestamp) UnmarshalJSON(b []byte) (err error) {
	// +info: https://stackoverflow.com/questions/25087960/json-unmarshal-time-that-isnt-in-rfc-3339-format
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse(timeLayout, s)
	return
}

type RtRevisionsData struct {
	Revision string
	Time     RtTimestamp
}

type ByTime []RtRevisionsData

func (a ByTime) Len() int           { return len(a) }
func (a ByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool { return a[i].Time.Before(a[j].Time.Time) }

type RtIndexJSON struct {
	//Reference string
	Revisions []RtRevisionsData
}
