package cliff

import "log"

// Command is the main type that gets exposed outside of the package and
// includes the command name and exposes the args and flags through methods
type Command struct {
	Name   string
	config *CommandConfig
	args   []string
}

// Arg returns the value of the argument as the specified index
func (c *Command) Arg(index int) string {
	if len(c.args) <= index {
		log.Fatalf("Arg with index \"%v\" is out of range!\n", index)
	}
	return c.args[index]
}

// FlagString returns the string value of a named flag
func (c *Command) FlagString(name string) string {
	return c.flag(name).String()
}

// FlagBool returns the bool value of a named flag
func (c *Command) FlagBool(name string) bool {
	return c.flag(name).Bool()
}

// FlagInt returns the int value of a named flag
func (c *Command) FlagInt(name string) int {
	return c.flag(name).Int()
}

// HasFlag returns a bool indicating the presence of a flag
func (c *Command) HasFlag(name string) bool {
	f := c.flag(name)
	return f != nil && f.cobraFlag != nil
}

// newCommand create a new *Command with a name and config
func newCommand(config *CommandConfig) *Command {
	return &Command{
		Name:   config.Name,
		config: config,
	}
}

// flag returns a *flag with the passed name presenet in the Command
func (c *Command) flag(name string) *flag {
	cc := c.config
	for _, f := range cc.Flags {
		if f.Long == name {
			return f
		}
	}
	if cc.parent != nil {
		return newCommand(cc.parent).flag(name)
	}
	return &flag{Long: name} // Return an empty Flag with just the name for debugging
}
