// Package search contains functionality to search references and packages in Artifactory.
package search

import (
	"regexp"

	"github.com/jgsogo/jcli-conan-center/types"

	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	servicesUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

// SearchPackages returns a list of packages matching the `referenceName` in the given `repository`. Use the argument
// `onlyLatestRecipe` to retrieve only packages that belong to the latest revision for each reference, and argument
// `onlyLatestPackage` to retrieve only the latest revision for each package.
func SearchPackages(serviceManager artifactory.ArtifactoryServicesManager, repository string, referenceName string, onlyLatestRecipe bool, onlyLatestPackage bool) ([]types.Package, error) {
	log.Info("Searching packages...")

	// Search all packages (search for the 'conaninfo.txt')
	specSearchPattern := repository + "/*/"
	if len(referenceName) > 0 {
		specSearchPattern = specSearchPattern + referenceName
	} else {
		specSearchPattern = specSearchPattern + "*"
	}
	specSearchPattern = specSearchPattern + "/*/*/package/*/*/conaninfo.txt"
	pkgPattern := regexp.MustCompile(`(?P<user>` + types.ValidConanChars + `*)\/(?P<name>` + types.ValidConanChars + `+)\/(?P<version>` + types.ValidConanChars + `+)\/(?P<channel>` + types.ValidConanChars + `*)\/(?P<revision>[a-z0-9]+)\/package\/(?P<pkgId>[a-z0-9]*)\/(?P<pkgRev>[a-z0-9]+)`)

	params := services.NewSearchParams()
	params.Pattern = specSearchPattern
	params.Recursive = false
	params.IncludeDirs = false

	reader, err := RunSearch(serviceManager, params)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	//
	allPackages := make(map[string]map[string]map[string][]types.Package)
	for resultItem := new(servicesUtils.ResultItem); reader.NextRecord(resultItem) == nil; resultItem = new(servicesUtils.ResultItem) {
		m := pkgPattern.FindStringSubmatch(resultItem.Path)
		var reference types.Reference
		user := m[1]
		channel := m[4]
		if user == types.FilesystemPlaceHolder {
			reference = types.Reference{Name: m[2], Version: m[3], User: nil, Channel: nil, Revision: m[5]}
		} else {
			reference = types.Reference{Name: m[2], Version: m[3], User: &user, Channel: &channel, Revision: m[5]}
		}

		if len(referenceName) > 0 && referenceName != reference.Name {
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
			rtRevisions, err := ParseRevisions(serviceManager, repository+"/"+key+"/index.json")
			if err != nil {
				return nil, err
			}
			latestRevision := rtRevisions[len(rtRevisions)-1]
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
				rtRevisions, err := ParseRevisions(serviceManager, repository+"/"+key+"/package/"+keyId+"/index.json")
				if err != nil {
					return nil, err
				}
				latestRevision := rtRevisions[len(rtRevisions)-1]
				i := Search(len(elementId), func(i int) bool {
					return latestRevision.Revision == elementId[i].Revision
				})
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
