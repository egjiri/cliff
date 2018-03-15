package cli

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	ex "github.com/egjiri/go-utils/exec"
	"github.com/spf13/cobra"
)

func (c command) addRunWithBashCommands(cmd *cobra.Command) {
	if len(c.Run) == 0 {
		return
	}
	cmd.Run = func(cc *cobra.Command, args []string) {
		verbose := cmd.Flag("verbose").Value.String() == "true"
		var content string
		for _, run := range c.Run {
			heading := transformBashContent(run.Heading, args, cmd)
			bashSetup := transformBashContent(run.Setup, args, cmd)
			bashCommand := transformBashContent(run.Execute, args, cmd)

			content += fmt.Sprintf("%s\n", bashSetup)
			if verbose && heading != "" {
				content += fmt.Sprintf("printf \"\n\033[0;32m%s...\033[m\n\"\n", heading)
				content += fmt.Sprintf("printf \"\033[0;34m==> \033[m\033[0;1m%s\033[m\n\"\n", bashCommand)
			}
			content += fmt.Sprintf("%s\n", bashCommand)
		}
		if verbose {
			content += fmt.Sprintf("printf \"\n\033[0;32m%s\033[m\n\"\n", "Finished!")
		}
		executeBash(content)
	}
}

func transformBashContent(content string, args []string, cmd *cobra.Command) string {
	// Replace the content args placeholders with the values of the args
	for i, arg := range args {
		content = strings.Replace(content, fmt.Sprintf("args[%v]", i), arg, 1)
	}
	// Replace the content flag placeholders with the values of the flags
	matches := regexp.MustCompile("flags\\[\"(.+?)\"\\]").FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		content = strings.Replace(content, match[0], cmd.Flag(match[1]).Value.String(), 1)
	}
	return content
}

func executeBash(content string) {
	tmpfile, err := ioutil.TempFile("", "cli")
	defer os.Remove(tmpfile.Name()) // clean up
	if err != nil {
		log.Fatal(err)
	}
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}
	ex.Execute("/bin/bash", tmpfile.Name())
}
