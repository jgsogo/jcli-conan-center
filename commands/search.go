package commands

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/jfrog/jfrog-cli-core/artifactory/commands"
	"github.com/jfrog/jfrog-cli-core/artifactory/utils"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/jgsogo/jcli-conan-center/search"
)

func GetSearchCommand() components.Command {
	return components.Command{
		Name:        "search",
		Description: "Return references to Conan packages found",
		Aliases:     []string{"s"},
		Arguments:   getSearchArguments(),
		Flags:       getSearchFlags(),
		EnvVars:     []components.EnvVar{},
		Action: func(c *components.Context) error {
			return searchCmd(c)
		},
	}
}

func getSearchFlags() []components.Flag {
	return []components.Flag{
		components.StringFlag{
			Name:         "server-id",
			Description:  "Artifactory server ID configured using the config command. If not specified, the default configured Artifactory server is used.",
			DefaultValue: "",
		},
		components.StringFlag{
			Name:         "ref-name",
			Description:  "Name of the reference to search (only the name). If not set, it will search for all references",
			DefaultValue: "",
		},
		components.BoolFlag{
			Name:         "packages",
			Description:  "If specified, it will retrieve also packages",
			DefaultValue: false,
		},
		components.BoolFlag{
			Name:         "only-latest",
			Description:  "If specified, it will retrieve only the latest revision",
			DefaultValue: false,
		},
	}
}

func getSearchArguments() []components.Argument {
	return []components.Argument{
		{
			Name:        "repo",
			Description: "Artifactory repository name",
		},
	}
}

func searchCmd(c *components.Context) error {
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

	// Search
	if c.GetBoolFlagValue("packages") {
		log.Info("Command search - retrieve packages")
		referenceName := c.GetStringFlagValue("ref-name")
		log.Info(fmt.Sprintf(" - ref-name: %s", referenceName))
		onlyLatest := c.GetBoolFlagValue("only-latest")
		packages, err := search.SearchPackages(rtDetails, repository, referenceName, onlyLatest, onlyLatest)
		if err != nil {
			return err
		}
		if len(packages) > 0 {
			log.Output(fmt.Sprintf("Found %d packages:", len(packages)))
			for _, pkg := range packages {
				log.Output(pkg.String())
			}
		}
	} else {
		log.Info("Command search - retrieve recipes")
		referenceName := c.GetStringFlagValue("ref-name")
		log.Info(fmt.Sprintf(" - ref-name: %s", referenceName))
		references, err := search.SearchReferences(serviceManager, repository, referenceName, c.GetBoolFlagValue("only-latest"))
		if err != nil {
			return err
		}
		if len(references) > 0 {
			//log.Output(fmt.Sprintf("Found %d references:", len(references)))
			for _, ref := range references {
				log.Output(ref.String())
			}
		}
	}
	return nil
}
