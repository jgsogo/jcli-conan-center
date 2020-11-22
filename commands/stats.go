package commands

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/jfrog/jfrog-cli-core/artifactory/commands"
	"github.com/jfrog/jfrog-cli-core/artifactory/commands/generic"
	"github.com/jfrog/jfrog-cli-core/artifactory/spec"
	"github.com/jfrog/jfrog-cli-core/artifactory/utils"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

type Reference struct {
	Name     string
	Version  string
	User     string
	Channel  string
	Revision string
}

func (ref *Reference) String() string {
	if len(ref.User) > 0 {
		return fmt.Sprintf("%s/%s@%s/%s#%s", ref.Name, ref.Version, ref.User, ref.Channel, ref.Revision)
	} else {
		return fmt.Sprintf("%s/%s#%s", ref.Name, ref.Version, ref.Revision)
	}
}

func GetStatsCommand() components.Command {
	return components.Command{
		Name:        "stats",
		Description: "Print some stats about Conan packages in remote",
		Aliases:     []string{"st"},
		Arguments:   getStatsArguments(),
		Flags:       getStatsFlags(),
		EnvVars:     []components.EnvVar{},
		Action: func(c *components.Context) error {
			return statsCmd(c)
		},
	}
}

func getStatsFlags() []components.Flag {
	return []components.Flag{
		components.StringFlag{
			Name:         "server-id",
			Description:  "Artifactory server ID configured using the config command. If not specified, the default configured Artifactory server is used.",
			DefaultValue: "",
		},
		components.StringFlag{
			Name:         "ref-name", // TODO: Implementation pending (this should act as a filter)
			Description:  "Name of a reference. If not specified, it will iterate every reference",
			DefaultValue: "",
		},
	}
}

func getStatsArguments() []components.Argument {
	return []components.Argument{
		{
			Name:        "repo",
			Description: "Artifactory repository name",
		},
	}
}

func searchReferences(repository string, rtDetails *config.ArtifactoryDetails) ([]Reference, error) {
	// Search all the 'conanfile.py' files inside the repository
	log.Output("Searching references in repository", repository)
	references := []Reference{}

	// Pattern to search
	specFile := spec.NewBuilder().Pattern(repository + "/**/conanfile.py").IncludeDirs(false).BuildSpec()
	validConanChars := "[a-zA-Z0-9_][a-zA-Z0-9_\\+\\.-]"
	referencePattern := regexp.MustCompile(repository + "\\/(?P<user>" + validConanChars + "*)\\/(?P<name>" + validConanChars + "+)\\/(?P<version>" + validConanChars + "+)\\/(?P<channel>" + validConanChars + "*)\\/(?P<revision>[a-z0-9]+)\\/export\\/conanfile\\.py")
	searchCmd := generic.NewSearchCommand()
	searchCmd.SetRtDetails(rtDetails).SetSpec(specFile)
	reader, err := searchCmd.Search()
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	for searchResult := new(utils.SearchResult); reader.NextRecord(searchResult) == nil; searchResult = new(utils.SearchResult) {
		m := referencePattern.FindStringSubmatch(searchResult.Path)
		user := m[1]
		if user == "_" {
			user = ""
		}
		channel := m[4]
		if channel == "_" {
			channel = ""
		}
		references = append(references, Reference{Name: m[2], Version: m[3], User: user, Channel: channel, Revision: m[5]})
	}
	return references, nil
}

func statsCmd(c *components.Context) error {
	if len(c.Arguments) != 1 {
		return errors.New("Wrong number of arguments. Expected: 1, " + "Received: " + strconv.Itoa(len(c.Arguments)))
	}

	rtDetails, err := commands.GetConfig(c.GetStringFlagValue("server-id"), true)
	if err != nil {
		log.Error(err)
		return err
	}

	// Check if repository exists
	repository := c.Arguments[0]
	log.Output("Inspect repository", repository)
	artAuth, err := rtDetails.CreateArtAuthConfig()
	if err != nil {
		log.Error("Error", err)
		return err
	}
	err = utils.CheckIfRepoExists(repository, artAuth)
	if err != nil {
		log.Error("Error2", err)
		return err
	}

	// Return references
	references, err := searchReferences(repository, rtDetails)
	if err != nil {
		return err
	}
	log.Output("Found", len(references), "Conan references")
	for _, ref := range references {
		fmt.Println(ref.String())
	}

	/*
		servicesManager, err := utils.CreateServiceManager(rtDetails, false)
		if err != nil {
			return err
		}
	*/



	log.Output("Done!")
	return nil
}
