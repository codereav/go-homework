package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	type Assert struct {
		VarName  string
		Expected interface{}
	}
	testCases := []struct {
		Descr     string
		Cmd       []string
		EnvValues Environment
		Asserts   []Assert
	}{
		{
			Descr: "All env vars are available after command running",
			Cmd:   []string{"./testdata/echo.sh", `arg1`, `arg2`},
			EnvValues: Environment{
				"FOO":   EnvValue{Value: "foo"},
				"ADDED": EnvValue{Value: `"false"`},
				"ANY":   EnvValue{Value: `123`},
			},
			Asserts: []Assert{
				{VarName: "FOO", Expected: `foo`},
				{VarName: "ANY", Expected: `123`},
				{VarName: "ADDED", Expected: `"false"`},
			},
		},
		{
			Descr: "All env vars marked to delete are not available after command running",
			Cmd:   []string{"./testdata/echo.sh", `arg1`, `arg2`},
			EnvValues: Environment{
				"FOO":   EnvValue{Value: "foo", NeedRemove: true},
				"ADDED": EnvValue{Value: `"false"`, NeedRemove: true},
				"ANY":   EnvValue{Value: `123`, NeedRemove: true},
			},
			Asserts: []Assert{
				{VarName: "FOO", Expected: ""},
				{VarName: "ANY", Expected: ""},
				{VarName: "ADDED", Expected: ""},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Descr, func(t *testing.T) {
			RunCmd(tc.Cmd, tc.EnvValues)
			for _, assert := range tc.Asserts {
				val, ok := os.LookupEnv(assert.VarName)
				if assert.Expected == "" {
					require.Equal(t, false, ok)
				}
				require.Equal(t, assert.Expected, val)
			}
		})
	}
}
