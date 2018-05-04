package executable

import (
	"github.com/egjiri/cliff/cliff"
)

func init() {
	cliff.AddRunToCommand("bash-completion", func(cmd cliff.Command, args []string) {
		outputPath := cmd.Flag("output").Value.String()
		cliff.ConfigureAndGenerateBashCompletionFile(outputPath)
	})
}
