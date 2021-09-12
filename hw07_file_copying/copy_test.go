package main

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	err := os.Mkdir("testFiles", os.ModePerm)
	checkSystemError(t, err)

	defer func() {
		err = os.RemoveAll("testFiles")
		if err != nil {
			t.FailNow()
		}
	}()

	t.Run("unsupported file", func(t *testing.T) {
		err = Copy("./testFiles/src", "./testFiles/dst.txt", 0, 0)
		require.ErrorIs(t, ErrUnsupportedFile, err)

		err = Copy("./testFiles/src.txt", "./testFiles/dst", 0, 0)
		require.ErrorIs(t, ErrUnsupportedFile, err)
	})

	t.Run("not existed src file", func(t *testing.T) {
		err = Copy("./testFiles/src.txt", "./testFiles/dst.txt", 0, 0)
		require.Equal(t, "open ./testFiles/src.txt: no such file or directory", err.Error())
	})

	t.Run("wrong toPath", func(t *testing.T) {
		srcFile, err := os.CreateTemp("testFiles", "src*.txt")
		checkSystemError(t, err)

		stat, err := srcFile.Stat()
		checkSystemError(t, err)

		err = Copy("./testFiles/"+stat.Name(), "./testFiles/something/dst.txt", 0, 0)
		require.Equal(t, "open ./testFiles/something/dst.txt: no such file or directory", err.Error())
	})

	t.Run("offsets and limits", func(t *testing.T) {
		srcFilePath := "./testFiles/src.txt"
		srcFile, err := os.Create(srcFilePath)
		checkSystemError(t, err)

		text := "something really important\nand long"

		_, err = srcFile.WriteString(text)
		checkSystemError(t, err)

		tests := []struct {
			offset int64
			limit  int64
			result string
			error  error
		}{
			{offset: 0, limit: 0, result: text},
			{offset: 0, limit: 10, result: text[0:10]},
			{offset: 1, limit: 10, result: text[1:11]},
			{offset: 0, limit: 1000, result: text},
			{offset: -10, limit: 1000, error: ErrNegativeOffset},
			{offset: -10, limit: -10, error: ErrNegativeOffset},
			{offset: 10, limit: -10, error: ErrNegativeLimit},
			{offset: 100, limit: 10, error: ErrOffsetExceedsFileSize},
		}

		for _, data := range tests {
			dstFile, err := os.CreateTemp("./testFiles", "dst*.txt")
			checkSystemError(t, err)

			stat, err := dstFile.Stat()
			checkSystemError(t, err)

			err = Copy(srcFilePath, "./testFiles/"+stat.Name(), data.offset, data.limit)

			if data.error != nil {
				require.ErrorIs(t, data.error, err)
				continue
			}

			require.Equal(t, nil, err)

			buf, err := io.ReadAll(dstFile)
			checkSystemError(t, err)

			err = dstFile.Close()
			checkSystemError(t, err)

			require.Equal(t, data.result, string(buf))
		}

		err = srcFile.Close()
		checkSystemError(t, err)
	})
}

func checkSystemError(t *testing.T, err error) {
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
}
