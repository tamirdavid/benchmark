package main

import (
	"benchmark/commands"

	"github.com/jfrog/jfrog-cli-core/v2/plugins"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
)

func main() {
	plugins.PluginMain(getApp())
}

func getApp() components.App {
	app := components.App{}
	app.Name = "benchmark"
	app.Description = "Easily test uploads/downloads"
	app.Version = "v0.1.4"
	app.Commands = getCommands()
	return app
}

func getCommands() []components.Command {
	return []components.Command{
		commands.DownloadCommand(),
		commands.UploadCommand(),
	}

}
