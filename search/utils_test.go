package search

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/jfrog/jfrog-cli-core/artifactory/spec"
	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
	"github.com/stretchr/testify/assert"
)

const (
	contentRevisions = `{
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
	r := ioutil.NopCloser(strings.NewReader(contentRevisions))
	return r, nil
}

func (esm *MockArtifactoryServicesManager) SearchFiles(params services.SearchParams) (*content.ContentReader, error) {
	fmt.Println(params.GetPattern())
	return content.NewContentReader("filePath", "arrayKey"), nil
}

func TestParseRevisions(t *testing.T) {
	servicesManager := MockArtifactoryServicesManager{}
	revisions, err := ParseRevisions(&servicesManager, "indexPath")
	assert.Nil(t, err)
	assert.Equal(t, len(revisions), 2)
	assert.Equal(t, revisions[0].Revision, "5918010f58ef4294511ff176ccc236b0")
	assert.Equal(t, revisions[1].Revision, "3c07b6a54477e856d429493d01c85636")
}

func TestRunSearch(t *testing.T) {
	servicesManager := MockArtifactoryServicesManager{}

	specFile := *spec.NewBuilder().Pattern("a/pattern").IncludeDirs(false).BuildSpec()
	revisions, err := RunSearch(&servicesManager, specFile)
	assert.Nil(t, err)
	assert.Nil(t, revisions)

}
