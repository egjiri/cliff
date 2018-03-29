package executable

import (
	"github.com/egjiri/cliff/cli"
)

// Execute configures the CLI and executes the root command
func Execute() {
	cli.ConfigureFromFile("cli.yml")
	cli.ConfigureSubcommandFromFile("cli.yml")
	cli.Execute()
}
