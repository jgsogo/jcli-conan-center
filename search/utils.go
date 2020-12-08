package search

import (
	"encoding/json"
	"io/ioutil"
	"sort"

	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/jgsogo/jcli-conan-center/types"
)

func Search(length int, f func(index int) bool) int {
    for index := 0; index < length; index++ {
        if f(index) {
            return index
        }
    }
    return -1
}

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

func RunSearch(servicesManager artifactory.ArtifactoryServicesManager, searchParams services.SearchParams) (*content.ContentReader, error) {
	reader, err := servicesManager.SearchFiles(searchParams)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return reader, err
}
