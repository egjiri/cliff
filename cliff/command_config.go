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

func (c *CommandConfig) buildCommand() *CommandConfig {
	cmd := &cobra.Command{
		Use:   c.Name,
		Short: c.Short,
		Long:  c.Long,
	}
	c.cobraCmd = cmd
	c.addRunWithBashCommands(cmd)
	c.addArgs(cmd)
	c.addFlags(cmd)
	updateTemplates(cmd)
	c.addchildren()
	commands[c.Name] = c
	return c
}

func (c *CommandConfig) addchildren() {
	for _, command := range c.Children {
		command.buildCommand()
		command.parent = c
		c.cobraCmd.AddCommand(command.cobraCmd)
	}
}

func (c *CommandConfig) addFlags(cmd *cobra.Command) {
	for _, flag := range c.Flags {
		flag.setFlag(cmd)
		flag.markRequiredFlags(cmd)
	}
	addHelpFlag(cmd) // Add the help flag to each command
}

func (c *CommandConfig) addArgs(cmd *cobra.Command) {
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
