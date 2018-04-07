package executable

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"

	"github.com/egjiri/cliff/cli"
	ex "github.com/egjiri/go-utils/exec"
	yaml "gopkg.in/yaml.v2"
)

type config struct {
	Name string
}

func init() {
	cli.AddRunToCommand("build", func(cmd cli.Command, args []string) {
		exPath, err := os.Getwd()
		if err != nil {
			log.Fatal("Error: ", err)
		}

		command := fmt.Sprintf("docker run --rm -v %s:/data -e GOOS=%s egjiri/cliff:0.0.1", exPath, runtime.GOOS)
		ex.Execute(command)

		newName := fmt.Sprintf("%s/%s", cmd.Flag("output").Value.String(), name())
		if err := os.Rename("cliff", newName); err != nil {
			log.Fatal("Error: ", err)
		}
		fmt.Println("Built binary:", newName)
	})
}

func name() string {
	yamlConfigContent, err := ioutil.ReadFile("cli.yml")
	if err != nil {
		log.Fatal(err)
	}

	var c config
	if err := yaml.Unmarshal(yamlConfigContent, &c); err != nil {
		log.Fatal(err)
	}
	return c.Name
}
