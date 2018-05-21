package cliff

import (
	"github.com/spf13/cobra"
)

// CommandConfig is the main type that stores all the command configuration
// such as argumens and flags as well as the command name, long and short
// descriptions as well as subcommands
type CommandConfig struct {
	Name, Short, Long string
	Args              interface{}
	Flags             []*flag
	Run               interface{}
	Children          []*CommandConfig `yaml:"commands"`
	parent            *CommandConfig
	cobraCmd          *cobra.Command
}

func (c *CommandConfig) configureCobraCommand() *CommandConfig {
	cmd := &cobra.Command{
		Use:   c.Name,
		Short: c.Short,
		Long:  c.Long,
	}
	c.cobraCmd = cmd
	updateTemplates(cmd)
	c.addRunWithBashCommands()
	c.addArgs()
	c.addFlags()
	c.addchildren()
	c.addToCommands()
	return c
}

func (c *CommandConfig) key() string {
	if c.parent == nil || c.parent == rootCmd {
		return c.Name
	}
	return c.parent.key() + "." + c.Name
}

func (c *CommandConfig) addArgs() {
	cmd := c.cobraCmd
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

func (c *CommandConfig) addFlags() {
	cmd := c.cobraCmd
	for _, flag := range c.Flags {
		flag.setFlag(cmd)
		flag.markRequiredFlags(cmd)
	}
	addHelpFlag(cmd) // Add the help flag to each command
}

func (c *CommandConfig) addchildren() {
	for _, command := range c.Children {
		command.parent = c
		command.configureCobraCommand()
		c.cobraCmd.AddCommand(command.cobraCmd)
	}
}

func (c *CommandConfig) addToCommands() {
	if c != rootCmd {
		commands[c.key()] = c
	}
}
