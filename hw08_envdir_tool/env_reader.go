package main

import (
	"bufio"
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var ErrWrongName = errors.New("wrong file name error: name can not contain \"=\"")

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	infoSlice, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := Environment{}
	for _, info := range infoSlice {
		if strings.Contains(info.Name(), "=") {
			return nil, ErrWrongName
		}

		if info.IsDir() {
			continue
		}

		path := filepath.Join(dir, info.Name())
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		scanner := bufio.NewScanner(file)
		scanner.Scan()

		buf := bytes.ReplaceAll(scanner.Bytes(), []byte{0x00}, []byte("\n"))
		env[info.Name()] = EnvValue{
			Value:      strings.TrimRight(string(buf), " \t"),
			NeedRemove: len(buf) == 0,
		}

		err = file.Close()

		if err != nil {
			return nil, err
		}
	}

	return env, nil
}
