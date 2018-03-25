package cli

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	yaml "gopkg.in/yaml.v2"
)

type Command struct {
	Name, Short, Long string
	Args              interface{}
	Flags             []flag
	Commands          []Command
	Run               interface{}
	cobraCmd          *cobra.Command
}

func (cmd Command) Flag(name string) *pflag.Flag {
	return cmd.cobraCmd.Flag(name)
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
	Run  func(cmd Command, args []string)
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

// ConfigureFromFile reads the contented of a passed file path and then calls Configure with it
func ConfigureFromFile(path string) error {
	yamlConfigContent, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	Configure(yamlConfigContent)
	return nil
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
func AddRunToCommand(name string, runFunc func(cmd Command, arg []string)) {
	*runs = append(*runs, run{name, runFunc})
}

func init() {
	log.SetFlags(0)
}

func setupRootCmd() {
	*rootCmd = *rootCommandFromConfigFile().buildCommand()
	addVerboseFlagToRootCmd()
}

func attachRunToCommands() {
	for _, r := range *runs {
		if cmd, ok := (*commands)[r.Name]; ok {
			cmd.Run = func(cmd *cobra.Command, args []string) {
				r.Run(Command{cobraCmd: cmd}, args)
			}
		}
	}
}

func rootCommandFromConfigFile() *Command {
	var rootCommand Command
	if err := yaml.Unmarshal(*config, &rootCommand); err != nil {
		log.Fatal(err)
	}
	return &rootCommand
}

func (c Command) buildCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   c.Name,
		Short: c.Short,
		Long:  c.Long,
	}
	c.addRunWithBashCommands(cmd)
	c.addArgs(cmd)
	c.addFlags(cmd)
	c.addCommands(cmd)
	updateTemplates(cmd)
	(*commands)[c.Name] = cmd
	return cmd
}

func (c Command) addCommands(parentCmd *cobra.Command) {
	for _, command := range c.Commands {
		parentCmd.AddCommand(command.buildCommand())
	}
}
