package commands

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/jfrog/jfrog-cli-core/artifactory/commands"
	"github.com/jfrog/jfrog-cli-core/artifactory/utils"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	servicesUtils "github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/jgsogo/jcli-conan-center/search"
	"github.com/jgsogo/jcli-conan-center/types"
)

// GetPropertiesGetCommand returns object description for the command 'properties'
func GetPropertiesGetCommand() components.Command {
	return components.Command{
		Name:        "properties",
		Description: "Return properties for a given Conan reference",
		Aliases:     []string{"p"},
		Arguments:   getPropertiesGetArguments(),
		Flags:       getPropertiesGetFlags(),
		EnvVars:     []components.EnvVar{},
		Action: func(c *components.Context) error {
			return propertiesGetCmd(c)
		},
	}
}

func getPropertiesGetFlags() []components.Flag {
	return []components.Flag{
		components.StringFlag{
			Name:         "server-id",
			Description:  "Artifactory server ID configured using the config command. If not specified, the default configured Artifactory server is used.",
			DefaultValue: "",
		},
		components.BoolFlag{
			Name:         "packages",
			Description:  "If specified, it will retrieve also packages",
			DefaultValue: false,
		},
	}
}

func getPropertiesGetArguments() []components.Argument {
	return []components.Argument{
		{
			Name:        "repo",
			Description: "Artifactory repository name",
		},
		{
			Name:        "reference",
			Description: "Conan reference to work with (use v2 style, without trailing @). If no revision is given, it will use latest one",
		},
	}
}

func propertiesGetCmd(c *components.Context) error {
	if len(c.Arguments) != 2 {
		return errors.New("Wrong number of arguments. Expected: 2, " + "Received: " + strconv.Itoa(len(c.Arguments)))
	}

	rtDetails, err := commands.GetConfig(c.GetStringFlagValue("server-id"), true)
	if err != nil {
		log.Error(err)
		return err
	}

	// Check if repository exists
	repository := c.Arguments[0]
	log.Output("Work on repository", repository)
	artAuth, err := rtDetails.CreateArtAuthConfig()
	if err != nil {
		return err
	}
	err = utils.CheckIfRepoExists(repository, artAuth)
	if err != nil {
		return err
	}

	// Create services manager
	serviceManager, err := utils.CreateServiceManager(rtDetails, false)
	if err != nil {
		return err
	}

	log.Info("Command properties-get")
	reference := c.Arguments[1]
	log.Info(fmt.Sprintf(" - input reference: %s", reference))

	// Search for the specific revision in the repository
	rtReference, err := types.ParseStringReference(reference)
	if err != nil {
		return err
	}
	if rtReference.Revision == "" { // Search for the latest revision
		rtRevisions, err := search.ParseRevisions(serviceManager, repository+"/"+rtReference.RtPath(false)+"/index.json")
		if err != nil {
			return err
		}
		rtReference.Revision = rtRevisions[len(rtRevisions)-1].Revision
	}
	log.Info(" - working reference:", rtReference.ToString(true))

	// Get properties for the given reference
	properties, err := search.ReadReferenceProperties(serviceManager, repository, *rtReference)
	if err != nil {
		return err
	}
	log.Output(fmt.Sprintf("Reference '%s':", rtReference.ToString(true)))
	for i := range properties {
		prop := properties[i]
		log.Output(fmt.Sprintf("  %s: %s", prop.Key, prop.Value))
	}

	if c.GetBoolFlagValue("packages") {
		// Get all packages for the given reference
		specSearchPattern := repository + "/" + rtReference.RtPath(true) + "/package/*/*/conaninfo.txt"
		params := services.NewSearchParams()
		params.Pattern = specSearchPattern
		params.Recursive = false
		params.IncludeDirs = false

		reader, err := search.RunSearch(serviceManager, params)
		if err != nil {
			return err
		}
		defer reader.Close()

		//
		pkgPattern := regexp.MustCompile(rtReference.RtPath(true) + "/package/" + `(?P<pkgId>[a-z0-9]*)\/(?P<pkgRev>[a-z0-9]+)`)
		for resultItem := new(servicesUtils.ResultItem); reader.NextRecord(resultItem) == nil; resultItem = new(servicesUtils.ResultItem) {
			m := pkgPattern.FindStringSubmatch(resultItem.Path)
			pkgReference := types.Package{Ref: *rtReference, PackageId: m[1], Revision: m[2]}
			properties, err := search.ReadPackageProperties(serviceManager, repository, pkgReference)
			if err != nil {
				return err
			}
			log.Output(fmt.Sprintf("Package '%s':", pkgReference.ToString(true)))
			for i := range properties {
				prop := properties[i]
				log.Output(fmt.Sprintf("  %s: %s", prop.Key, prop.Value))
			}
		}
	}

	return nil
}
