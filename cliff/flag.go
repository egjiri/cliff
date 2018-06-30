package cliff

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type flag struct {
	Long, Short, Type, Description string
	Default                        interface{}
	Global, Required               bool
	Enum                           []string
	cobraFlag                      *pflag.Flag
}

// String returns the string value of the flag
func (f *flag) String() string {
	if f == nil || f.cobraFlag == nil {
		if f == nil {
			log.Fatal("No flag found!")
		}
		log.Fatalf("No flag \"%s\" found!\n", f.Long)
	}
	if err := f.validate(); err != nil {
		log.Fatal(err)
	}
	return f.cobraFlag.Value.String()
}

// Bool returns the bool value of the flag
func (f *flag) Bool() bool {
	value, err := strconv.ParseBool(f.String())
	if err != nil {
		log.Fatal(err)
	}
	return value
}

// Int returns the int value of the flag
func (f *flag) Int() int {
	value, err := strconv.Atoi(f.String())
	if err != nil {
		log.Fatal(err)
	}
	return value
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
	f.cobraFlag = cmd.Flag(f.Long)
}

func (f *flag) stringValue() string {
	var defaultValue string
	if value := f.Default; value != nil {
		defaultValue = value.(string)
	}
	return defaultValue
}

func (f *flag) boolValue() bool {
	var defaultValue bool
	if value := f.Default; value != nil {
		defaultValue = value.(bool)
	}
	return defaultValue
}

func (f *flag) intValue() int {
	var defaultValue int
	if value := f.Default; value != nil {
		defaultValue = value.(int)
	}
	return defaultValue
}

func (f *flag) cmdFlags(cmd *cobra.Command) *pflag.FlagSet {
	if f.Global {
		return cmd.PersistentFlags()
	}
	return cmd.Flags()
}

func (f *flag) markRequiredFlags(cmd *cobra.Command) {
	if f.Required {
		name := f.Long
		if f.Global {
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
	for _, f := range rootCmd.Flags {
		if f.Long == "verbose" {
			return
		}
	}
	(&flag{
		Long:        "verbose",
		Type:        "boolean",
		Description: "Verbosity of the logs",
		Default:     false,
		Global:      true,
	}).setFlag(rootCmd.cobraCmd)
}

// validate checks if the flag type is a string and if there are any enums set and
// ensures that the flag value matches one of the enums otherwise it returns an error
func (f *flag) validate() error {
	if f.Type != "string" || len(f.Enum) == 0 {
		return nil
	}
	val := f.cobraFlag.Value.String()
	for _, validVal := range f.Enum {
		if val == validVal {
			return nil
		}
	}
	options := fmt.Sprintf("[%v]", strings.Join(f.Enum, ", "))
	errMsg := fmt.Sprintf("Invalid \"%s\" flag option: %s\nAvailable options are: %s\n", f.Long, val, options)
	return errors.New(errMsg)
	// ===== Sample Output =====
	// Invalid "environment" flag option: testing
	// Available options are: [development, staging, production]
}
