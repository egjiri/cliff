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
	yaml "gopkg.in/yaml.v2"
)

type command struct {
	Name, Short, Long, Run string
	Args                   interface{}
	Flags                  []flag
	Commands               []command
}

type flag struct {
	Long, Short, Type, Description string
	Default                        interface{}
	Global, Required               bool
}

type run struct {
	Name string
	Run  func(cmd *cobra.Command, args []string)
}

var config = &[]byte{}
var rootCmd = &cobra.Command{}
var commands = &map[string]*cobra.Command{}
var runs = &[]run{}

// Configure sets the content of the yaml config file and sets up the commands
func Configure(yamlConfigContent []byte) {
	*config = yamlConfigContent
	setupRootCmd()
	attachRunToCommands()
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// AddRunToCommand provies a mechanism to attach a Run function to a command
func AddRunToCommand(name string, runFunc func(cmd *cobra.Command, args []string)) {
	*runs = append(*runs, run{name, runFunc})
}

func init() {
	log.SetFlags(0)
}

func setupRootCmd() {
	*rootCmd = *rootCommandFromConfigFile().buildCommand()
	updateTemplates()
}

func attachRunToCommands() {
	for _, r := range *runs {
		if cmd, ok := (*commands)[r.Name]; ok {
			cmd.Run = r.Run
		}
	}
}

// TODO Run this if there is no config file passed in
func readConfigFile() {
	fileName := "cli.yml"
	yamlContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	*config = yamlContent
}

func rootCommandFromConfigFile() *command {
	var rootCommand command
	if err := yaml.Unmarshal(*config, &rootCommand); err != nil {
		log.Fatal(err)
	}
	return &rootCommand
}

func (c command) buildCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   c.Name,
		Short: c.Short,
		Long:  c.Long,
	}
	c.addRun(cmd)
	c.addArgs(cmd)
	c.addFlags(cmd)
	c.addCommands(cmd)
	(*commands)[c.Name] = cmd
	return cmd
}

func (c command) addCommands(parentCmd *cobra.Command) {
	for _, command := range c.Commands {
		parentCmd.AddCommand(command.buildCommand())
	}
}

func (c command) addRun(cmd *cobra.Command) {
	if c.Run == "" {
		return
	}
	cmd.Run = func(cc *cobra.Command, args []string) {
		content := c.Run
		// Replace the content args placeholders with the values of the args
		for i, arg := range args {
			content = strings.Replace(content, fmt.Sprintf("args[%v]", i), arg, 1)
		}
		// Replace the content flag placeholders with the values of the flags
		matches := regexp.MustCompile("flags\\[\"(.+?)\"\\]").FindAllStringSubmatch(content, -1)
		for _, match := range matches {
			content = strings.Replace(content, match[0], cmd.Flag(match[1]).Value.String(), 1)
		}
		// Write the content to a bash file
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
}

func (c command) addArgs(cmd *cobra.Command) {
	if num, ok := c.Args.(int); ok {
		cmd.Args = cobra.ExactArgs(num)
	} else {
		if args, ok := c.Args.(map[interface{}]interface{}); ok {
			if len(args) == 1 {
				for k, v := range args {
					if value, ok := v.(int); ok {
						if k == "min" {
							cmd.Args = cobra.MinimumNArgs(value)
						} else if k == "max" {
							cmd.Args = cobra.MaximumNArgs(value)
						}
					}
				}
			} else if len(args) == 2 {
				var min, max int
				for k, v := range args {
					if value, ok := v.(int); ok {
						if k == "min" {
							min = value
						} else if k == "max" {
							max = value
						}
					}
				}
				cmd.Args = cobra.RangeArgs(min, max)
			}
		}
	}
}