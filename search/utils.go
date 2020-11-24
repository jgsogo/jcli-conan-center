package search

import (
	"encoding/json"
	"io/ioutil"
	"sort"

	"github.com/jfrog/jfrog-client-go/artifactory"
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
