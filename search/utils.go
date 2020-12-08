package search

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"

	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	servicesUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/jgsogo/jcli-conan-center/types"
)

// Search function invokes the given function with the natural numbers from zero to the input `length`. Useful to search into an array.
func Search(length int, f func(index int) bool) int {
	for index := 0; index < length; index++ {
		if f(index) {
			return index
		}
	}
	return -1
}

// ParseRevisions parses and 'index.json' file stored in Artifactory and returns a sorted list of revisions.
func ParseRevisions(serviceManager artifactory.ArtifactoryServicesManager, indexPath string) ([]types.RtRevisionsData, error) {
	ioReaderCloser, err := serviceManager.ReadRemoteFile(indexPath)
	if err != nil {
		return nil, err
	}
	defer ioReaderCloser.Close()
	content, err := ioutil.ReadAll(ioReaderCloser)
	if err != nil {
		return nil, err
	}
	var revisions types.RtIndexJSON
	err = json.Unmarshal(content, &revisions)
	if err != nil {
		return nil, err
	}
	sort.Sort(types.ByTime(revisions.Revisions))
	return revisions.Revisions, nil
}

// RunSearch return the content according to the given `searchParams`.
func RunSearch(servicesManager artifactory.ArtifactoryServicesManager, searchParams services.SearchParams) (*content.ContentReader, error) {
	reader, err := servicesManager.SearchFiles(searchParams)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return reader, err
}

func readProperties(serviceManager artifactory.ArtifactoryServicesManager, repository string, path string) ([]servicesUtils.Property, error) {
	params := services.NewSearchParams()
	params.Pattern = repository + "/" + path
	params.Recursive = false
	params.IncludeDirs = true

	reader, err := RunSearch(serviceManager, params)
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	
	for resultItem := new(servicesUtils.ResultItem); reader.NextRecord(resultItem) == nil; resultItem = new(servicesUtils.ResultItem) {
		return resultItem.Properties, nil
	}

	return nil, fmt.Errorf("Properties for path '%s' not found", path)
}

// ReadReferenceProperties returns the properties for a given Conan reference `ref` in the given `repository`.
func ReadReferenceProperties(serviceManager artifactory.ArtifactoryServicesManager, repository string, ref types.Reference) ([]servicesUtils.Property, error) {
	props, err := readProperties(serviceManager, repository, ref.RtPath(true))
	if err != nil {
		return nil, fmt.Errorf("Properties for reference '%s' not found", ref.RtPath(true))
	}
	return props, nil
}

// ReadPackageProperties returns the properties for a given Conan package `pkg` in the given `repository`.
func ReadPackageProperties(serviceManager artifactory.ArtifactoryServicesManager, repository string, pkg types.Package) ([]servicesUtils.Property, error) {
	props, err := readProperties(serviceManager, repository, pkg.RtPath(true))
	if err != nil {
		return nil, fmt.Errorf("Properties for package '%s' not found", pkg.RtPath(true))
	}
	return props, nil
}
