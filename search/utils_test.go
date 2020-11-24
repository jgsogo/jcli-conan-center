package search

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory"
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

type MockArtifactoryServicesManager struct {
	artifactory.EmptyArtifactoryServicesManager
}

func (esm *MockArtifactoryServicesManager) ReadRemoteFile(readPath string) (io.ReadCloser, error) {
	r := ioutil.NopCloser(strings.NewReader(content))
	return r, nil
}

func TestParseRevisions(t *testing.T) {
	servicesManager := MockArtifactoryServicesManager{}
	revisions, err := ParseRevisions(&servicesManager, "indexPath")
	assert.Nil(t, err)
	assert.Equal(t, len(revisions), 2)
	assert.Equal(t, revisions[0].Revision, "5918010f58ef4294511ff176ccc236b0")
	assert.Equal(t, revisions[1].Revision, "3c07b6a54477e856d429493d01c85636")
}
