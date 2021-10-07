package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("work test", func(t *testing.T) {
		env := Environment{
			"VALUE":   {Value: "VALUE", NeedRemove: false},
			"EMPTY":   {Value: "", NeedRemove: false},
			"REMOVED": {Value: "", NeedRemove: true},
		}

		result := RunCmd([]string{"/bin/bash", "testdata/echo.sh", "arg1=1", "arg2=2"}, env)
		require.Equal(t, 0, result)
	})

	t.Run("command error", func(t *testing.T) {
		result := RunCmd([]string{"/bin/bash", "wrong", "arg1=1", "arg2=2"}, Environment{})

		require.NotEqual(t, 0, result)
	})
}
