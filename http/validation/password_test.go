package validation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"goyave.dev/goyave/v5/validation"
)

func TestPasswordValidator(t *testing.T) {
	cases := []struct {
		value any
		want  bool
	}{
		{value: "weak", want: false},
		{value: "p4ssW0rd_", want: true},
		{value: "p4ssW0rd", want: false},
		{value: 123, want: false},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%v", c.value), func(t *testing.T) {
			v := Password()
			assert.Equal(t, "password", v.Name())
			pass := v.Validate(&validation.Context{Value: c.value})
			assert.Equal(t, c.want, pass)
		})
	}
}
