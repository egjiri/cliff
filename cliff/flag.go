package cliff

import (
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Flag struct {
	Long, Short, Type, Description string
	Default                        interface{}
	Global, Required               bool
	cobraFlag                      *pflag.Flag
}

// String returns the string value of the flag
func (f *Flag) String() string {
	if f == nil || f.cobraFlag == nil {
		if f == nil {
			log.Fatal("No flag found!")
		}
		log.Fatalf("No flag \"%s\" found!\n", f.Long)
	}
	return f.cobraFlag.Value.String()
}

// Bool returns the bool value of the flag
func (f *Flag) Bool() bool {
	value, err := strconv.ParseBool(f.String())
	if err != nil {
		log.Fatal(err)
	}
	return value
}

// Int returns the int value of the flag
func (f *Flag) Int() int {
	value, err := strconv.Atoi(f.String())
	if err != nil {
		log.Fatal(err)
	}
	return value
}

func (f *Flag) setFlag(cmd *cobra.Command) {
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
	f.cobraFlag = cmd.Flag(f.Long)
}

func (f *Flag) stringValue() string {
	var defaultValue string
	if value := (*f).Default; value != nil {
		defaultValue = value.(string)
	}
	return defaultValue
}

func (f *Flag) boolValue() bool {
	var defaultValue bool
	if value := (*f).Default; value != nil {
		defaultValue = value.(bool)
	}
	return defaultValue
}

func (f *Flag) intValue() int {
	var defaultValue int
	if value := (*f).Default; value != nil {
		defaultValue = value.(int)
	}
	return defaultValue
}

func (f *Flag) cmdFlags(cmd *cobra.Command) *pflag.FlagSet {
	if (*f).Global {
		return cmd.PersistentFlags()
	}
	return cmd.Flags()
}

func (f *Flag) markRequiredFlags(cmd *cobra.Command) {
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
	(&Flag{
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
	(&Flag{
		Long:        "verbose",
		Type:        "boolean",
		Description: "Verbosity of the logs",
		Default:     false,
		Global:      true,
	}).setFlag(rootCmd.cobraCmd)
}
