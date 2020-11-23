package types

import (
	"encoding/json"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	content = `{
		"reference": "b2/4.0.0@_/_",
		"revisions": [{
			"revision": "3c07b6a54477e856d429493d01c85636",
			"time": "2020-09-16T14:05:05.965+0000"
		}, {
			"revision": "5918010f58ef4294511ff176ccc236b0",
			"time": "2020-08-17T15:20:47.871+0000"
		}]
	}`
)

func TestParseJSON(t *testing.T) {
	var revisions RtIndexJSON
	err := json.Unmarshal([]byte(content), &revisions)

	assert.Nil(t, err)
	assert.Equal(t, len(revisions.Revisions), 2)
	assert.Equal(t, revisions.Revisions[0].Revision, "3c07b6a54477e856d429493d01c85636")
	assert.Equal(t, revisions.Revisions[1].Revision, "5918010f58ef4294511ff176ccc236b0")
}

func TestOrderByTime(t *testing.T) {
	var revisions RtIndexJSON
	err := json.Unmarshal([]byte(content), &revisions)
	assert.Nil(t, err)

	sort.Sort(ByTime(revisions.Revisions))
	assert.Equal(t, len(revisions.Revisions), 2)
	assert.Equal(t, revisions.Revisions[0].Revision, "5918010f58ef4294511ff176ccc236b0")
	assert.Equal(t, revisions.Revisions[1].Revision, "3c07b6a54477e856d429493d01c85636")

}
