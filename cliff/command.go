package cliff

import (
	"log"

	"github.com/spf13/cobra"
)

// Command is the main type that stores all the command configuration
// such as argumens and flags as well as the command name, long and short
// descriptions as well as subcommands
type Command struct {
	Name, Short, Long string
	Args              interface{}
	Flags             []*Flag
	Run               interface{}
	Children          []*Command `yaml:"commands"`
	parent            *Command
	cobraCmd          *cobra.Command
	args              []string
}

// Flag returns a *Flag with the passed name presenet in the Command
func (c *Command) Flag(name string) *Flag {
	for _, f := range c.Flags {
		if f.Long == name {
			return f
		}
	}
	if c.parent != nil {
		return c.parent.Flag(name)
	}
	return &Flag{Long: name} // Return an empty Flag with just the name for debugging
}

// Arg returns the value of the argument as the specified index
func (c *Command) Arg(index int) string {
	if len(c.args) <= index {
		log.Fatalf("Arg with index \"%v\" is out of range!\n", index)
	}
	return c.args[index]
}

func (c *Command) buildCommand() *Command {
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

func (c *Command) addchildren() {
	for _, command := range c.Children {
		command.buildCommand()
		command.parent = c
		c.cobraCmd.AddCommand(command.cobraCmd)
	}
}

func (c *Command) addFlags(cmd *cobra.Command) {
	for _, flag := range c.Flags {
		flag.setFlag(cmd)
		flag.markRequiredFlags(cmd)
	}
	addHelpFlag(cmd) // Add the help flag to each command
}

func (c *Command) addArgs(cmd *cobra.Command) {
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
