package executable

//go:generate go-bindata -pkg data -o ../data/go-bindata.go ../cli.yml

import (
	"github.com/egjiri/cliff/cli"
	"github.com/egjiri/cliff/data"
)

func init() {
	if err := cli.ConfigureFromFile("cli.yml"); err != nil {
		cli.Configure(data.MustAsset("../cli.yml"))
	}
	setupCustomCommands()
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
