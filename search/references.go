package search

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"

	"github.com/jgsogo/jcli-conan-center/types"

	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	servicesUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

func SearchReferences(serviceManager artifactory.ArtifactoryServicesManager, repository string, referenceName string, onlyLatest bool) ([]types.Reference, error) {
	log.Info("Searching references...")

	// Search all references (search for the 'conanfile.py')
	specSearchPattern := repository + "/*/"
	if len(referenceName) > 0 {
		specSearchPattern = specSearchPattern + referenceName
	} else {
		specSearchPattern = specSearchPattern + "*"
	}
	specSearchPattern = specSearchPattern + "/*/*/conanfile.py"
	log.Debug(fmt.Sprintf("Search references using specPattern '%s'", specSearchPattern))
	referencePattern := regexp.MustCompile(`(?P<user>` + types.ValidConanChars + `*)\/(?P<name>` + types.ValidConanChars + `+)\/(?P<version>` + types.ValidConanChars + `+)\/(?P<channel>` + types.ValidConanChars + `*)\/(?P<revision>[a-z0-9]+)\/export`)

	params := services.NewSearchParams()
	params.Pattern = specSearchPattern
	params.Recursive = false
	params.IncludeDirs = false

	reader, err := RunSearch(serviceManager, params)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	references := make(map[string][]types.Reference)
	for searchResult := new(servicesUtils.ResultItem); reader.NextRecord(searchResult) == nil; searchResult = new(servicesUtils.ResultItem) {
		m := referencePattern.FindStringSubmatch(searchResult.Path)
		var reference types.Reference
		user := m[1]
		channel := m[4]
		if user == types.FilesystemPlaceHolder {
			reference = types.Reference{Name: m[2], Version: m[3], User: nil, Channel: nil, Revision: m[5]}
		} else {
			reference = types.Reference{Name: m[2], Version: m[3], User: &user, Channel: &channel, Revision: m[5]}
		}
		references[reference.ToString(false)] = append(references[reference.ToString(false)], reference)
	}

	// Filter duplicated references using 'index.json' (if onlyLatest)
	retReferences := []types.Reference{}
	for _, element := range references {
		if onlyLatest && len(element) > 1 {
			rtRevisions, err := ParseRevisions(serviceManager, repository+"/"+element[0].RtPath(false)+"/index.json")
			if err != nil {
				return nil, err
			}
			latestRevision := rtRevisions[len(rtRevisions)-1]
			i := sort.Search(len(element), func(i int) bool { return latestRevision.Revision == element[i].Revision })
			retReferences = append(retReferences, element[i])
		} else {
			retReferences = append(retReferences, element...)
		}
	}
	log.Info("Found", strconv.Itoa(len(retReferences)), "references.")
	return retReferences, nil
}
