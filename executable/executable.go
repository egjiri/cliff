package executable

import "github.com/egjiri/cliff/cliff"

// Execute configures the CLI and executes the root command
func Execute() {
	cliff.ConfigureFromFile("cli.yml")
	cliff.ConfigureSubcommandFromFile("cli.yml")
	cliff.Execute()
}
