package search

import (
	"regexp"
	"sort"
	"strings"

	"github.com/jgsogo/jcli-conan-center/types"

	"github.com/jfrog/jfrog-cli-core/artifactory/commands/generic"
	"github.com/jfrog/jfrog-cli-core/artifactory/spec"
	"github.com/jfrog/jfrog-cli-core/artifactory/utils"
	"github.com/jfrog/jfrog-cli-core/utils/config"
)

func SearchPackages(rtDetails *config.ArtifactoryDetails, repository string, ref *types.Reference, onlyLatestRecipe bool, onlyLatestPackage bool) ([]types.Package, error) {
	// Search all packages (search for the 'conaninfo.txt')
	if ref != nil && onlyLatestRecipe {
		panic("Incompatible input arguments, do not request to filter by latest recipe if a reference is provided")
	}

	startsWith := repository + "/"
	if ref != nil {
		startsWith = startsWith + ref.RtPath(true)
	} else {
		startsWith = startsWith + "*/*/*/*"
	}
	startsWith = startsWith + "/package/*/*/conaninfo.txt"
	specFile := spec.NewBuilder().Pattern(startsWith).IncludeDirs(false).BuildSpec()
	pkgPattern := regexp.MustCompile(repository + `\/(?P<user>` + types.ValidConanChars + `*)\/(?P<name>` + types.ValidConanChars + `+)\/(?P<version>` + types.ValidConanChars + `+)\/(?P<channel>` + types.ValidConanChars + `*)\/(?P<revision>[a-z0-9]+)\/package\/(?P<pkgId>[a-z0-9]*)\/(?P<pkgRev>[a-z0-9]+)\/conaninfo.txt`)

	searchCmd := generic.NewSearchCommand()
	searchCmd.SetRtDetails(rtDetails).SetSpec(specFile)
	reader, err := searchCmd.Search()
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	//
	allPackages := make(map[string]map[string]map[string][]types.Package)
	for searchResult := new(utils.SearchResult); reader.NextRecord(searchResult) == nil; searchResult = new(utils.SearchResult) {
		m := pkgPattern.FindStringSubmatch(strings.TrimPrefix(searchResult.Path, startsWith))
		user := m[1]
		if user == "_" {
			user = ""
		}
		channel := m[4]
		if channel == "_" {
			channel = ""
		}
		reference := types.Reference{Name: m[2], Version: m[3], User: user, Channel: channel, Revision: m[5]}
		if ref != nil && *ref != reference {
			panic("Mismatch references!")
		}
		conanPackage := types.Package{Ref: reference, PackageId: m[6], Revision: m[7]}
		inner, ok := allPackages[conanPackage.Ref.RtPath(false)]
		if !ok {
			inner = make(map[string]map[string][]types.Package)
			allPackages[conanPackage.Ref.RtPath(false)] = inner
		}
		inner2, ok := inner[conanPackage.Ref.Revision]
		if !ok {
			inner2 = make(map[string][]types.Package)
			inner[conanPackage.Ref.Revision] = inner2
		}
		inner2[conanPackage.PackageId] = append(inner2[conanPackage.PackageId], conanPackage)
	}

	// Filter recipes using 'index.json' (if onlyLatestRecipe)
	filteredPackages := make(map[string]map[string][]types.Package)
	for key, element := range allPackages {
		if onlyLatestRecipe && len(element) > 1 {
			rtRevisions, err := ParseRevisions(rtDetails, repository+"/"+key+"/index.json")
			if err != nil {
				return nil, err
			}
			latestRevision := rtRevisions[len(rtRevisions)-1]
			//i := sort.Search(len(element), func(i int) bool { return latestRevision.Revision == element[i].Revision })
			for k, v := range element[latestRevision.Revision] {
				inner, ok := filteredPackages[key+"/"+latestRevision.Revision]
				if !ok {
					inner = make(map[string][]types.Package)
					filteredPackages[key+"/"+latestRevision.Revision] = inner
				}
				inner[k] = v
			}
		} else {
			for rrev, elements := range element {
				for k, v := range elements {
					inner, ok := filteredPackages[key+"/"+rrev]
					if !ok {
						inner = make(map[string][]types.Package)
						filteredPackages[key+"/"+rrev] = inner
					}
					inner[k] = v
				}
			}
		}
	}

	// Filter packages using 'index.json' (if onlyLatestPackages)
	packages := []types.Package{}
	for key, element := range filteredPackages {
		if onlyLatestPackage && len(element) > 1 {
			for keyId, elementId := range element {
				rtRevisions, err := ParseRevisions(rtDetails, repository+"/"+key+"/package/"+keyId+"/index.json")
				if err != nil {
					return nil, err
				}
				latestRevision := rtRevisions[len(rtRevisions)-1]
				i := sort.Search(len(elementId), func(i int) bool { return latestRevision.Revision == elementId[i].Revision })
				packages = append(packages, elementId[i])
			}
		} else {
			for _, elementId := range element {
				packages = append(packages, elementId...)
			}
		}
	}
	return packages, nil
}
