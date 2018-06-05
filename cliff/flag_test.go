package cliff

import (
	"testing"

	et "github.com/egjiri/go-kit/testing"
)

func Test_validate(t *testing.T) {
	f := &flag{
		Long: "environment",
		Type: "string",
		Enum: []string{"development", "staging", "production"},
	}
	et.Assert(t,
		et.TestErr{
			Name: "value present in enum",
			Actual: func(f *flag) error {
				f.cobraFlag = cobraFlag(f.Long, "development")
				return f.validate()
			}(f),
			Expected: false,
		},
		et.TestErr{
			Name: "value absent in enum",
			Actual: func(f *flag) error {
				f.cobraFlag = cobraFlag(f.Long, "testing")
				return f.validate()
			}(f),
			Expected: true,
		},
		et.TestErr{
			Name: "no enums",
			Actual: func(f *flag) error {
				f.Enum = []string{}
				return f.validate()
			}(f),
			Expected: false,
		},
	)
}
