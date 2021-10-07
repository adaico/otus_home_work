package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	err := os.Mkdir("test", os.ModePerm)

	defer func() {
		err = os.RemoveAll("test")

		if err != nil {
			require.FailNow(t, err.Error())
		}
	}()

	if err != nil {
		require.FailNow(t, err.Error())
	}

	t.Run("test data", func(t *testing.T) {
		env, err := ReadDir("testdata/env")
		if err != nil {
			require.FailNow(t, err.Error())
		}

		result := Environment{
			"BAR":   {Value: "bar", NeedRemove: false},
			"EMPTY": {Value: "", NeedRemove: false},
			"FOO":   {Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO": {Value: "\"hello\"", NeedRemove: false},
			"UNSET": {Value: "", NeedRemove: true},
		}

		require.Equal(t, result, env)
	})

	t.Run("= in file name", func(t *testing.T) {
		_, err = os.Create("test/NAME_WITH_=")

		defer func() {
			err = os.Remove("test/NAME_WITH_=")

			if err != nil {
				require.FailNow(t, err.Error())
			}
		}()

		if err != nil {
			require.FailNow(t, err.Error())
		}

		if err != nil {
			require.FailNow(t, err.Error())
		}

		_, err = ReadDir("test")

		require.ErrorIs(t, ErrWrongName, err)
	})

	t.Run("only tab and whitespace in file, dir among files", func(t *testing.T) {
		file, err := os.Create("test/TAB_AND_WHITESPACE")
		if err != nil {
			require.FailNow(t, err.Error())
		}

		err = os.Mkdir("test/SOME_DIR", os.ModeDir)

		if err != nil {
			require.FailNow(t, err.Error())
		}

		_, err = file.WriteString("\t ")

		if err != nil {
			require.FailNow(t, err.Error())
		}

		env, err := ReadDir("test")
		if err != nil {
			require.FailNow(t, err.Error())
		}

		result := Environment{
			"TAB_AND_WHITESPACE": {Value: "", NeedRemove: false},
		}

		require.Equal(t, result, env)
	})
}
