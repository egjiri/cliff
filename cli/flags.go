package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func (c Command) addFlags(cmd *cobra.Command) {
	for _, flag := range c.Flags {
		f := &flag
		f.setFlag(cmd)
		f.markRequiredFlags(cmd)
	}
	addHelpFlag(cmd) // Add the help flag to each command
}

func (f *flag) setFlag(cmd *cobra.Command) {
	ff := *f
	cmdFlags := ff.cmdFlags(cmd)
	if ff.Short != "" {
		switch ff.Type {
		case "string":
			cmdFlags.StringP(ff.Long, ff.Short, f.stringValue(), ff.Description)
		case "boolean":
			cmdFlags.BoolP(ff.Long, ff.Short, f.boolValue(), ff.Description)
		case "integer":
			cmdFlags.IntP(ff.Long, ff.Short, f.intValue(), ff.Description)
		}
	} else {
		switch ff.Type {
		case "string":
			cmdFlags.String(ff.Long, f.stringValue(), ff.Description)
		case "boolean":
			cmdFlags.Bool(ff.Long, f.boolValue(), ff.Description)
		case "integer":
			cmdFlags.Int(ff.Long, f.intValue(), ff.Description)
		}
	}
}

func (f *flag) stringValue() string {
	var defaultValue string
	if value := (*f).Default; value != nil {
		defaultValue = value.(string)
	}
	return defaultValue
}

func (f *flag) boolValue() bool {
	var defaultValue bool
	if value := (*f).Default; value != nil {
		defaultValue = value.(bool)
	}
	return defaultValue
}

func (f *flag) intValue() int {
	var defaultValue int
	if value := (*f).Default; value != nil {
		defaultValue = value.(int)
	}
	return defaultValue
}

func (f *flag) cmdFlags(cmd *cobra.Command) *pflag.FlagSet {
	if (*f).Global {
		return cmd.PersistentFlags()
	}
	return cmd.Flags()
}

func (f *flag) markRequiredFlags(cmd *cobra.Command) {
	flag := *f
	if flag.Required {
		name := flag.Long
		if flag.Global {
			cmd.MarkPersistentFlagRequired(name)
		} else {
			cmd.MarkFlagRequired(name)
		}
	}
}

func addHelpFlag(cmd *cobra.Command) {
	if cmd.Flag("help") != nil {
		return
	}
	(&flag{
		Long:        "help",
		Short:       "h",
		Type:        "boolean",
		Description: fmt.Sprintf("Help for %s", cmd.Name()),
		Default:     false,
	}).setFlag(cmd)
}

func addVerboseFlagToRootCmd() {
	if rootCmd.Flag("verbose") != nil {
		return
	}
	(&flag{
		Long:        "verbose",
		Type:        "boolean",
		Description: "Verbosity of the logs",
		Default:     false,
		Global:      true,
	}).setFlag(rootCmd)
}
