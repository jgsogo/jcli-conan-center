package search

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
	"github.com/jgsogo/jcli-conan-center/types"
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
		}, {
			"revision": "7777777777",
			"time": "2020-08-15T15:20:47.871+0000"
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
	if params.Pattern == "the/pattern/to/search/for" {
		return content.NewContentReader("filePath", "arrayKey"), nil
	} else if params.Pattern == "repository/_/name/version/_/rrev" {
		wd, _ := os.Getwd()
		filePath := filepath.Join(wd, "testdata/search_utils_props_reference.json")

		tmpFile, _ := ioutil.TempFile(os.TempDir(), "prefix-")
		fileContent, _ := ioutil.ReadFile(filePath)
		_, _ = tmpFile.Write(fileContent)
		tmpFile.Close()

		reader := content.NewContentReader(tmpFile.Name(), "results")
		return reader, nil
	} else if params.Pattern == "repository/_/name/version/_/rrev/package/pkgID/prev" {
		wd, _ := os.Getwd()
		filePath := filepath.Join(wd, "testdata/search_utils_props_package.json")

		tmpFile, _ := ioutil.TempFile(os.TempDir(), "prefix-")
		fileContent, _ := ioutil.ReadFile(filePath)
		_, _ = tmpFile.Write(fileContent)
		tmpFile.Close()

		reader := content.NewContentReader(tmpFile.Name(), "results")
		return reader, nil
	} else {
		wd, _ := os.Getwd()
		filePath := filepath.Join(wd, "testdata/not_found.json")

		tmpFile, _ := ioutil.TempFile(os.TempDir(), "prefix-")
		fileContent, _ := ioutil.ReadFile(filePath)
		_, _ = tmpFile.Write(fileContent)
		tmpFile.Close()

		reader := content.NewContentReader(tmpFile.Name(), "results")
		return reader, nil
	}
}

func TestParseRevisions(t *testing.T) {
	servicesManager := MockArtifactoryServicesManager{}
	revisions, err := ParseRevisions(&servicesManager, "indexPath")
	assert.Nil(t, err)
	assert.Equal(t, len(revisions), 3)
	assert.Equal(t, revisions[0].Revision, "7777777777")
	assert.Equal(t, revisions[1].Revision, "5918010f58ef4294511ff176ccc236b0")
	assert.Equal(t, revisions[2].Revision, "3c07b6a54477e856d429493d01c85636")
}

func TestRunSearch(t *testing.T) {
	servicesManager := MockArtifactoryServicesManager{}
	params := services.NewSearchParams()
	params.Pattern = "the/pattern/to/search/for"
	params.Recursive = false
	params.IncludeDirs = false

	reader, err := RunSearch(&servicesManager, params)
	assert.Nil(t, err)
	assert.Equal(t, reader.GetFilePath(), "filePath")

}

func TestReadReferenceProperties(t *testing.T) {
	servicesManager := MockArtifactoryServicesManager{}
	reference := types.Reference{Name: "name", Version: "version", User: nil, Channel: nil, Revision: "rrev"}
	props, err := ReadReferenceProperties(&servicesManager, "repository", reference)
	assert.Nil(t, err)
	assert.Equal(t, 17, len(props))
	assert.Equal(t, "topics", props[0].Key)
	assert.Equal(t, "conan", props[0].Value)

	otherRef := types.Reference{Name: "other", Version: "version", User: nil, Channel: nil, Revision: "rrev"}
	props, err = ReadReferenceProperties(&servicesManager, "repository", otherRef)
	assert.NotNil(t, err)
	assert.Equal(t, "Properties for reference '_/other/version/_/rrev' not found", err.Error())
}

func TestReadPackageProperties(t *testing.T) {
	servicesManager := MockArtifactoryServicesManager{}
	reference := types.Reference{Name: "name", Version: "version", User: nil, Channel: nil, Revision: "rrev"}
	pkg := types.Package{Ref: reference, PackageId: "pkgID", Revision: "prev"}
	props, err := ReadPackageProperties(&servicesManager, "repository", pkg)
	assert.Nil(t, err)
	assert.Equal(t, 17, len(props))
	assert.Equal(t, "license", props[0].Key)
	assert.Equal(t, "BSL-1.0", props[0].Value)

	otherpkg := types.Package{Ref: reference, PackageId: "otherID", Revision: "prev"}
	props, err = ReadPackageProperties(&servicesManager, "repository", otherpkg)
	assert.NotNil(t, err)
	assert.Equal(t, "Properties for package '_/name/version/_/rrev/package/otherID/prev' not found", err.Error())
}
