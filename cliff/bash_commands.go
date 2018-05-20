package cliff

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

type bashCommands struct {
	Heading, Setup, Execute string
}

func (c *Command) addRunWithBashCommands(cmd *cobra.Command) {
	run := extractBashCommandsFromRun(c.Run)
	if len(run) == 0 {
		return
	}
	cmd.Run = func(cc *cobra.Command, args []string) {
		f := cmd.Flag("verbose")
		verbose := f != nil && f.Value.String() == "true"
		var content string
		for _, r := range run {
			heading := transformBashContent(r.Heading, args, cmd)
			bashSetup := transformBashContent(r.Setup, args, cmd)
			bashCommand := transformBashContent(r.Execute, args, cmd)

			content += fmt.Sprintf("%s\n", bashSetup)
			if verbose {
				if heading != "" {
					content += fmt.Sprintf("printf \"\n\033[0;32m%s...\033[m\n\"\n", heading)
				}
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
	content = strings.Replace(content, "args...", strings.Join(args, " "), 1)
	for i, arg := range args {
		content = strings.Replace(content, fmt.Sprintf("args[%v]", i), arg, -1)
	}
	// Replace the content flag placeholders with the values of the flags
	matches := regexp.MustCompile("flags\\[\"(.+?)\"\\]").FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		f := cmd.Flag(match[1])
		if f == nil {
			log.Fatalf("Error: Invalid flag \"%s\" used in command", match[1])
		}
		content = strings.Replace(content, match[0], f.Value.String(), 1)
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

func extractBashCommandsFromRun(run interface{}) []bashCommands {
	var bcs []bashCommands
	// Create and return a bashCommand if run is a string
	if exec, ok := run.(string); ok {
		return append(bcs, bashCommands{Execute: exec})
	}
	// Build and return a slice otherwise
	parts, ok := run.([]interface{})
	if !ok {
		return bcs
	}
	for _, part := range parts {
		p, ok := part.(map[interface{}]interface{})
		if !ok {
			return bcs
		}
		var heading, setup, execute string
		for k, v := range p {
			val := v.(string)
			switch k.(string) {
			case "heading":
				heading = val
			case "setup":
				setup = val
			case "execute":
				execute = val
			}
		}
		bcs = append(bcs, bashCommands{
			Heading: heading,
			Setup:   setup,
			Execute: execute,
		})
	}
	return bcs
}
