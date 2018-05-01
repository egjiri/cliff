package cliff

import "github.com/spf13/cobra"

func (c Command) addArgs(cmd *cobra.Command) {
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
