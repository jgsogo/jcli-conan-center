// Package types provides some abstraction over some data returned by JFrog Artifactory.
package types

import (
	"strings"
	"time"
)

const (
	timeLayout = "2006-01-02T15:04:05.999+0000"
)

// RtTimestamp represents a custom timestamp using format '2006-01-02T15:04:05.999+0000'. It allows 
// serializing and deserializing using the representation used by Artifactory.
type RtTimestamp struct {
	time.Time
}

// UnmarshalJSON overrides parsing from JSON a timestamp using the RtTimestamp type.
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

// RtRevisionsData represents the data associated to a Conan revision in Artifactory.
type RtRevisionsData struct {
	Revision string
	Time     RtTimestamp
}

// ByTime is a helper operator to order revisions by date.
type ByTime []RtRevisionsData

func (a ByTime) Len() int           { return len(a) }
func (a ByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool { return a[i].Time.Before(a[j].Time.Time) }

// RtIndexJSON represents the JSON where Artifactory stores Conan revisions (using 'index.json' files).
type RtIndexJSON struct {
	//Reference string
	Revisions []RtRevisionsData
}
