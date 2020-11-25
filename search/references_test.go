package search

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
	"github.com/stretchr/testify/assert"
)

type MockRtServicesManager struct {
	artifactory.EmptyArtifactoryServicesManager
}

func (esm *MockRtServicesManager) SearchFiles(params services.SearchParams) (*content.ContentReader, error) {
	wd, _ := os.Getwd()
	filePath := filepath.Join(wd, "testdata/search_references.json")

	tmpFile, _ := ioutil.TempFile(os.TempDir(), "prefix-")
	fileContent, _ := ioutil.ReadFile(filePath)
	_, _ = tmpFile.Write(fileContent)
	tmpFile.Close()

	reader := content.NewContentReader(tmpFile.Name(), "results")
	return reader, nil
}

func (esm *MockRtServicesManager) ReadRemoteFile(readPath string) (io.ReadCloser, error) {
	if readPath == "repository/_/b2/4.0.0/_/index.json" {
		return ioutil.NopCloser(strings.NewReader(`{
			"reference": "b2/4.0.0@_/_",
			"revisions": [{
				"revision": "3c07b6a54477e856d429493d01c85636",
				"time": "2020-09-16T14:05:05.965+0000"
			}, {
				"revision": "5918010f58ef4294511ff176ccc236b0",
				"time": "2020-08-17T15:20:47.871+0000"
			}]
		}`)), nil
	} else if readPath == "repository/_/b2/4.0.1/_/index.json" {
		return ioutil.NopCloser(strings.NewReader(`{
			"reference": "b2/4.0.1@_/_",
			"revisions": [{
				"revision": "fe103dcc7b9fa2226d82f5fb43af1d09",
				"time": "2020-09-16T14:06:23.885+0000"
			}, {
				"revision": "64a94a3e9fe90b33033ec9e00eb036e6",
				"time": "2020-08-17T15:20:52.616+0000"
			}]
		}`)), nil
	} else if readPath == "repository/_/b2/4.2.0/_/index.json" {
		return ioutil.NopCloser(strings.NewReader(`{
			"reference": "b2/4.2.0@_/_",
			"revisions": [{
				"revision": "efacbfac6ee3561ff07968a372b940af",
				"time": "2020-09-16T14:08:54.728+0000"
			}, {
				"revision": "7987eb34c600c944d8a30ffe090fd013",
				"time": "2020-08-17T15:21:01.757+0000"
			}]
		}`)), nil
	}
	return nil, nil
}

func TestSearchReferences(t *testing.T) {
	servicesManager := MockRtServicesManager{}
	params := services.NewSearchParams()
	params.Pattern = "the/pattern/to/search/for"
	params.Recursive = false
	params.IncludeDirs = false

	references, err := SearchReferences(&servicesManager, "repository", "name/version", false)
	assert.Nil(t, err)
	assert.Equal(t, 8, len(references))
}

func TestSearchReferencesLatest(t *testing.T) {
	servicesManager := MockRtServicesManager{}
	params := services.NewSearchParams()
	params.Pattern = "the/pattern/to/search/for"
	params.Recursive = false
	params.IncludeDirs = false

	references, err := SearchReferences(&servicesManager, "repository", "name/version", true)
	assert.Nil(t, err)
	sort.Slice(references, func(i, j int) bool {
		return references[i].String() < references[j].String()
	})
	assert.Equal(t, 5, len(references))
	assert.Equal(t, "b2/4.0.0#3c07b6a54477e856d429493d01c85636", references[0].String())
	assert.Equal(t, "b2/4.0.1#fe103dcc7b9fa2226d82f5fb43af1d09", references[1].String())
	assert.Equal(t, "b2/4.1.0#151655c3ac57c4adcc3681a2bf44e0af", references[2].String())
	assert.Equal(t, "b2/4.2.0#efacbfac6ee3561ff07968a372b940af", references[3].String())
	assert.Equal(t, "b2/4.3.0#ec8af29b790f5745890470ce4220ed50", references[4].String())
}
