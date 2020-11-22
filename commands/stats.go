package commands

import (
	"errors"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"strconv"
	"github.com/jfrog/jfrog-cli-core/artifactory/commands"
	"github.com/jfrog/jfrog-cli-core/artifactory/utils"
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
			Name:        "ref-name",
			Description: "Name of a reference. If not specified, it will iterate every reference",
			DefaultValue: "",
		},
	}
}

func getStatsArguments() []components.Argument {
	return []components.Argument{
		{
			Name:         "repo",
			Description:  "Artifactory repository name",
		},
	}
}

func statsCmd(c *components.Context) error {
	if len(c.Arguments) != 1 {
		return errors.New("Wrong number of arguments. Expected: 1, " + "Received: " + strconv.Itoa(len(c.Arguments)))
	}
	repository := c.Arguments[0]
	log.Output("Inspect repository", repository)
	
	rtDetails, err := commands.GetConfig(c.GetStringFlagValue("server-id"), true)
	if err != nil {
		return err
	}
	/*
	servicesManager, err := utils.CreateServiceManager(rtDetails, false)
	if err != nil {
		return err
	}
	*/
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

	log.Output("Done!")
	return nil
}
