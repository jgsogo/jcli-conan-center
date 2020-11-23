package search

import (
	"encoding/json"
	"io/ioutil"
	"sort"

	"github.com/jgsogo/jcli-conan-center/types"

	"github.com/jfrog/jfrog-cli-core/artifactory/utils"
	"github.com/jfrog/jfrog-cli-core/utils/config"
)

func ParseRevisions(rtDetails *config.ArtifactoryDetails, indexPath string) ([]types.RtRevisionsData, error) {
	serviceManager, err := utils.CreateServiceManager(rtDetails, false)
	if err != nil {
		return nil, err
	}
	// https://github.com/jfrog/jfrog-cli-core/blob/8a53bb7180151cf4093f714f2bc7949029f48e18/artifactory/utils/docker/buildinfo.go
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
