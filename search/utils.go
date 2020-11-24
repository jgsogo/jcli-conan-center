package search

import (
	"encoding/json"
	"io/ioutil"
	"sort"

	"github.com/jfrog/jfrog-cli-core/artifactory/spec"
	"github.com/jfrog/jfrog-cli-core/artifactory/utils"
	"github.com/jfrog/jfrog-client-go/artifactory"
	clientartutils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/jgsogo/jcli-conan-center/types"
)

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

func RunSearch(servicesManager artifactory.ArtifactoryServicesManager, sc spec.SpecFiles) (*content.ContentReader, error) {
	// Most of the implementation taken from https://github.com/jfrog/jfrog-cli-core/blob/master/artifactory/commands/generic/search.go

	// Search Loop
	log.Info("Searching artifacts...")
	var searchResults []*content.ContentReader
	for i := 0; i < len(sc.Files); i++ {
		searchParams, err := utils.GetSearchParams(sc.Get(i))
		if err != nil {
			log.Error(err)
			return nil, err
		}

		reader, err := servicesManager.SearchFiles(searchParams)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		searchResults = append(searchResults, reader)
		if i == 0 {
			defer func() {
				for _, reader := range searchResults {
					reader.Close()
				}
			}()
		}
	}
	reader, err := utils.AqlResultToSearchResult(searchResults)
	if err != nil {
		return nil, err
	}
	length, err := reader.Length()
	clientartutils.LogSearchResults(length)
	return reader, err
}
