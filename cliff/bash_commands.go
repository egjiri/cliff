package cliff

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	ex "github.com/egjiri/go-kit/exec"
	"github.com/egjiri/go-kit/ui/print"
	"github.com/spf13/cobra"
)

type bashCommands struct {
	Heading, Setup, Execute string
}

func (c *CommandConfig) addRunWithBashCommands() {
	run := extractBashCommandsFromRun(c.Run)
	if len(run) == 0 {
		return
	}
	command := newCommand(c)
	cmd := c.cobraCmd
	cmd.Run = func(cc *cobra.Command, args []string) {
		verbose := command.HasFlag("verbose") && command.FlagBool("verbose")
		var content string
		for _, r := range run {
			heading := transformBashContent(r.Heading, args, command)
			bashSetup := transformBashContent(r.Setup, args, command)
			bashCommand := transformBashContent(r.Execute, args, command)

			content += fmt.Sprintf("%s\n", bashSetup)
			if verbose {
				if heading != "" {
					print.Heading(heading)
				}
				print.Command(bashCommand)
			}
			content += fmt.Sprintf("%s\n", bashCommand)
		}
		ex.ExecuteBash(content)
		if verbose {
			print.Finished()
		}
	}
}

func transformBashContent(content string, args []string, c *Command) string {
	// Replace the content args placeholders with the values of the args
	content = strings.Replace(content, "args...", strings.Join(args, " "), 1)
	for i, arg := range args {
		content = strings.Replace(content, fmt.Sprintf("args[%v]", i), arg, -1)
	}
	// Replace the content flag placeholders with the values of the flags
	matches := regexp.MustCompile("flags\\[\"(.+?)\"\\]").FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		f := c.flag(match[1])
		if f == nil {
			log.Fatalf("Error: Invalid flag \"%s\" used in command", match[1])
		}
		content = strings.Replace(content, match[0], f.String(), 1)
	}
	return content
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
