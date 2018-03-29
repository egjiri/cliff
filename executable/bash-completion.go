package executable

import (
	"github.com/egjiri/cliff/cli"
)

func init() {
	cli.AddRunToCommand("bash-completion", func(cmd cli.Command, args []string) {
		path := cmd.Flag("output").Value.String()
		cli.GenerateBashCompletionFile(path)
	})
}
