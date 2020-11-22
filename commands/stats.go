package commands

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

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

func (ref *Reference) rtPath() string {
	user := ref.User
	if len(user) == 0 {
		user = "_"
	}
	channel := ref.Channel
	if len(channel) == 0 {
		channel = "_"
	}
	str := []string{user, ref.Name, ref.Version, channel, ref.Revision} 
	return strings.Join(str, "/")
}

type Package struct {
	Ref Reference
	PackageId string
	Revision string
}

func (pkg *Package) String() string {
	return fmt.Sprintf("%s:%s#%s", pkg.Ref, pkg.PackageId, pkg.Revision)
}

func (pkg *Package) rtPath() string {
	str := []string{pkg.Ref.rtPath(), "package", pkg.PackageId, pkg.Revision} 
	return strings.Join(str, "/")
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

func searchReferences(rtDetails *config.ArtifactoryDetails, repository string) ([]Reference, error) {
	// Search all the 'conanfile.py' files inside the repository

	specFile := spec.NewBuilder().Pattern(repository + "/**/conanfile.py").IncludeDirs(false).BuildSpec()
	validConanChars := "[a-zA-Z0-9_][a-zA-Z0-9_\\+\\.-]"
	referencePattern := regexp.MustCompile(repository + `\/(?P<user>` + validConanChars + `*)\/(?P<name>` + validConanChars + `+)\/(?P<version>` + validConanChars + `+)\/(?P<channel>` + validConanChars + `*)\/(?P<revision>[a-z0-9]+)\/export\/conanfile\.py`)

	searchCmd := generic.NewSearchCommand()
	searchCmd.SetRtDetails(rtDetails).SetSpec(specFile)
	reader, err := searchCmd.Search()
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	references := []Reference{}
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

func searchPackages(rtDetails *config.ArtifactoryDetails, repository string, ref Reference) ([]Package, error) {
	startsWith := repository + "/" + ref.rtPath() + "/package"
	specFile := spec.NewBuilder().Pattern(startsWith + "/*/*/conaninfo.txt").IncludeDirs(false).BuildSpec()
	
	pkgPattern := regexp.MustCompile(`\/(?P<pkgId>[a-z0-9]*)\/(?P<pkgRev>[a-z0-9]+)\/conaninfo.txt`)

	searchCmd := generic.NewSearchCommand()
	searchCmd.SetRtDetails(rtDetails).SetSpec(specFile)
	reader, err := searchCmd.Search()
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	packages := []Package{}
	for searchResult := new(utils.SearchResult); reader.NextRecord(searchResult) == nil; searchResult = new(utils.SearchResult) {
		m := pkgPattern.FindStringSubmatch(strings.TrimPrefix(searchResult.Path, startsWith))
		packages = append(packages, Package{Ref: ref, PackageId: m[1], Revision: m[2]})
	}
	return packages, nil
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

	// Search packages (first recipes and then packages)
	packages := []Package{}
	references, err := searchReferences(rtDetails, repository)
	if err != nil {
		return err
	}
	log.Output("Found", len(references), "Conan references")
	for _, ref := range references {
		refPackages, err := searchPackages(rtDetails, repository, ref)
		if err != nil {
			return err
		}
		log.Output(" -", ref.String(), ":", len(refPackages), "packages")
		packages = append(packages, refPackages...)
	}
	log.Output("Total found", len(packages),"packages")

	// Group packages


	log.Output("Done!")
	return nil
}
