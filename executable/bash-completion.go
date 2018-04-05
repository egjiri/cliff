package executable

import (
	"io/ioutil"
	"log"

	"github.com/egjiri/cliff/cli"
)

func init() {
	cli.AddRunToCommand("bash-completion", func(cmd cli.Command, args []string) {
		yamlConfigFilePath := cmd.Flag("file").Value.String()
		yamlConfigContent, err := ioutil.ReadFile(yamlConfigFilePath)
		if err != nil {
			log.Fatal(err)
		}
		cli.Configure(yamlConfigContent)

		outputPath := cmd.Flag("output").Value.String()
		cli.GenerateBashCompletionFile(outputPath)
	})
}
