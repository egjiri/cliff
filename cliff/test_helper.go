package cliff

import "github.com/spf13/pflag"

// Implement the pflag.Value interface
type flagValue struct {
	Value string
}

func (v flagValue) String() string   { return v.Value }
func (v flagValue) Set(string) error { return nil }
func (v flagValue) Type() string     { return "" }

func cobraFlag(name, value string) *pflag.Flag {
	return &pflag.Flag{
		Name:  name,
		Value: flagValue{value},
	}
}
