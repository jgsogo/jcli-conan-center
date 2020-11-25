package search

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
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
	tmpFile.Write(fileContent)
	tmpFile.Close()

	reader := content.NewContentReader(tmpFile.Name(), "results")
	return reader, nil
}

func (esm *MockRtServicesManager) ReadRemoteFile(readPath string) (io.ReadCloser, error) {
	fmt.Println(">>>> ", readPath)
	r := ioutil.NopCloser(strings.NewReader(contentRevisions))
	return r, nil
}

func TestSearchReferences(t *testing.T) {
	servicesManager := MockRtServicesManager{}
	params := services.NewSearchParams()
	params.Pattern = "the/pattern/to/search/for"
	params.Recursive = false
	params.IncludeDirs = false

	references, err := SearchReferences(&servicesManager, "repository", "name/version", false)
	assert.Nil(t, err)
	assert.Equal(t, 9, len(references))
}

func TestSearchReferencesLatest(t *testing.T) {
	servicesManager := MockRtServicesManager{}
	params := services.NewSearchParams()
	params.Pattern = "the/pattern/to/search/for"
	params.Recursive = false
	params.IncludeDirs = false

	references, err := SearchReferences(&servicesManager, "repository", "name/version", true)
	assert.Nil(t, err)
	assert.Equal(t, 9, len(references))
}
