// Package main constains a plugin for JFrog CLI, the command line interface
// to JFrog Artifactory. This plugin implements several commands realated to
// Conan (https://conan.io) and ConanCenter (https://conan.io/center). Most
// of these commands are intended to manage the Artifactory repository
// where the packages in ConanCenter are stored.
package main

import (
	"github.com/jfrog/jfrog-cli-core/plugins"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"github.com/jgsogo/jcli-conan-center/commands"
)

func main() {
	plugins.PluginMain(getApp())
}

func getApp() components.App {
	app := components.App{}
	app.Name = "conan-center"
	app.Description = "Manage Conan repository (ConanCenter like)."
	app.Version = "v0.1.0"
	app.Commands = getCommands()
	return app
}

func getCommands() []components.Command {
	return []components.Command{
		commands.GetStatsCommand(),
		commands.GetSearchCommand(),
		commands.GetPropertiesGetCommand(),
		commands.GetIndexReferenceCommand(),
	}
}
