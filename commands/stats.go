package commands

import (
	"errors"
	"strconv"

	"github.com/jgsogo/jcli-conan-center/search"
	"github.com/jgsogo/jcli-conan-center/types"

	"github.com/jfrog/jfrog-cli-core/artifactory/commands"
	"github.com/jfrog/jfrog-cli-core/artifactory/utils"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"github.com/jfrog/jfrog-client-go/utils/log"
)

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

	// Search packages (first recipes and then packages)
	packages := []types.Package{}
	references, err := search.SearchReferences(serviceManager, repository, "", false)
	if err != nil {
		return err
	}
	log.Output("Found", len(references), "Conan references")
	for _, ref := range references {
		refPackages, err := search.SearchPackages(serviceManager, repository, ref.Name, false, false)
		if err != nil {
			return err
		}
		log.Output(" -", ref.String(), ":", len(refPackages), "packages")
		packages = append(packages, refPackages...)
	}
	log.Output("Total found", len(packages), "packages")

	// Search packages (all at once)
	allPackages, err := search.SearchPackages(serviceManager, repository, "", false, false)
	if err != nil {
		return err
	}
	log.Output("Total found", len(allPackages), "packages")

	log.Output("Done!")
	return nil
}
