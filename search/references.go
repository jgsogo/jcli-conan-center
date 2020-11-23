package search

import (
	"regexp"
	"sort"

	"github.com/jgsogo/jcli-conan-center/types"

	"github.com/jfrog/jfrog-cli-core/artifactory/commands/generic"
	"github.com/jfrog/jfrog-cli-core/artifactory/spec"
	"github.com/jfrog/jfrog-cli-core/artifactory/utils"
	"github.com/jfrog/jfrog-cli-core/utils/config"
)

func SearchReferences(rtDetails *config.ArtifactoryDetails, repository string, onlyLatest bool) ([]types.Reference, error) {
	// Search all references (search for the 'conanfile.py')

	specFile := spec.NewBuilder().Pattern(repository + "/**/conanfile.py").IncludeDirs(false).BuildSpec()
	referencePattern := regexp.MustCompile(repository + `\/(?P<user>` + types.ValidConanChars + `*)\/(?P<name>` + types.ValidConanChars + `+)\/(?P<version>` + types.ValidConanChars + `+)\/(?P<channel>` + types.ValidConanChars + `*)\/(?P<revision>[a-z0-9]+)\/export\/conanfile\.py`)

	searchCmd := generic.NewSearchCommand()
	searchCmd.SetRtDetails(rtDetails).SetSpec(specFile)
	reader, err := searchCmd.Search()
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	references := make(map[string][]types.Reference)
	for searchResult := new(utils.SearchResult); reader.NextRecord(searchResult) == nil; searchResult = new(utils.SearchResult) {
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
			rtRevisions, err := ParseRevisions(rtDetails, repository+"/"+element[0].RtPath(false)+"/index.json")
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
	return retReferences, nil
}
