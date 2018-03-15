package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

type command struct {
	Name, Short, Long string
	Args              interface{}
	Flags             []flag
	Commands          []command
	Run               []bashCommands
}

type flag struct {
	Long, Short, Type, Description string
	Default                        interface{}
	Global, Required               bool
}

type bashCommands struct {
	Heading, Setup, Execute string
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
	c.addRunWithBashCommands(cmd)
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
