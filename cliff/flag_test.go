package cliff

import (
	"testing"

	et "github.com/egjiri/go-kit/testing"
)

func Test_validate(t *testing.T) {
	type args struct {
		value string
		enums []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"value present in enums", args{"development", []string{"development", "production"}}, false},
		{"value absent in enums", args{"staging", []string{"development", "production"}}, true},
		{"no enums everything valid", args{"development", []string{}}, false},
	}
	var testables []et.Testable
	for _, test := range tests {
		testables = append(testables, et.TestErr{
			Name: test.name,
			Actual: (&flag{
				Long:      test.args.value,
				Type:      "string",
				Enum:      test.args.enums,
				cobraFlag: cobraFlag("environment", test.args.value),
			}).validate(),
			Expected: test.wantErr,
		})
	}
	et.Assert(t, testables...)
}
