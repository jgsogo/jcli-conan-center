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
			Name:         "only-latest",
			Description:  "If specified, it will retrieve only the latest revision",
			DefaultValue: true,
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
			Description: "Conan reference to work with",
		},
	}
}

func propertiesGetCmd(c *components.Context) error {
	if len(c.Arguments) != 1 {
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

	// Search
	log.Info("Command properties-get -")
	reference := c.Arguments[1]
	log.Info(fmt.Sprintf(" - reference: %s", reference))
	references, err := search.SearchReferences(serviceManager, repository, reference, c.GetBoolFlagValue("only-latest"))
	if err != nil {
		return err
	}
	if len(references) > 0 {
		//log.Output(fmt.Sprintf("Found %d references:", len(references)))
		for _, ref := range references {
			log.Output(ref.String())
		}
	}

	return nil
}
