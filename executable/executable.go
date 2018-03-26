package executable

//go:generate go-bindata -pkg data -o ../data/go-bindata.go ../cli.yml

import (
	"github.com/egjiri/cliff/cli"
	"github.com/egjiri/cliff/data"
)

func init() {
	setupCustomCommands() // This must be run before Configure
	cli.Configure(data.MustAsset("../cli.yml"))
	cli.ConfigureSubcommandFromFile("cli.yml")
}

// Execute configures the CLI and executes the root command
func Execute() {
	cli.Execute()
}

func setupCustomCommands() {
	cli.AddRunToCommand("bash-completion", func(cmd cli.Command, args []string) {
		path := cmd.Flag("output").Value.String()
		cli.GenerateBashCompletionFile(path)
	})
}
