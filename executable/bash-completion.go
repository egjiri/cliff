package executable

import (
	"io/ioutil"
	"log"

	"github.com/egjiri/cliff/cliff"
)

func init() {
	cliff.AddRunToCommand("bash-completion", func(cmd cliff.Command, args []string) {
		yamlConfigFilePath := cmd.Flag("file").Value.String()
		yamlConfigContent, err := ioutil.ReadFile(yamlConfigFilePath)
		if err != nil {
			log.Fatal(err)
		}
		cliff.Configure(yamlConfigContent)

		outputPath := cmd.Flag("output").Value.String()
		cliff.GenerateBashCompletionFile(outputPath)
	})
}
