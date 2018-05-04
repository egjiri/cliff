package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	prefix := strings.Replace(os.Args[1], "github.com/egjiri/cliff/vendor/", "", 1)
	modifyCliCode(prefix + "github.com/egjiri/cliff/cliff/cliff.go")
}

func modifyCliCode(filePath string) {
	code, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	newCode := strings.Replace(string(code), "import (", "import (\n	\"github.com/egjiri/cliff/data\"", 1)
	newCode = strings.Replace(newCode, "yamlConfigContent, err := ioutil.ReadFile(path)", "yamlConfigContent, err := data.Asset(\"cli.yml\")", -1)
	if !strings.Contains(newCode, "ioutil.") {
		newCode = strings.Replace(newCode, "\"io/ioutil\"", "// \"io/ioutil\"", 1)
	}

	if err := ioutil.WriteFile(filePath, []byte(newCode), 0644); err != nil {
		log.Fatal(err)
	}
}
