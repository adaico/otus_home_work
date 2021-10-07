package main

import (
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	command := cmd[0]
	params := cmd[1:]
	execCmd := exec.Command(command, params...)

	execCmd.Env = getEnv(env)
	execCmd.Stdout = os.Stdout

	if err := execCmd.Run(); err != nil {
		r := regexp.MustCompile(`^exit status (\d+)$`)
		stringCode := r.FindStringSubmatch(err.Error())[1]
		code, err := strconv.Atoi(stringCode)
		if err != nil {
			log.Fatal(err)
		}

		return code
	}

	return
}

func getEnv(env Environment) []string {
	finalEnvMap := make(map[string]string)
	for _, stringValue := range os.Environ() {
		arrayValue := strings.Split(stringValue, "=")
		name := arrayValue[0]
		value := arrayValue[1]

		finalEnvMap[name] = value
	}

	for name, value := range env {
		if value.NeedRemove {
			delete(finalEnvMap, name)

			continue
		}

		finalEnvMap[name] = value.Value
	}

	finalEnvSlice := make([]string, 0, len(finalEnvMap))
	for name, value := range finalEnvMap {
		finalEnvSlice = append(finalEnvSlice, strings.Join([]string{name, value}, "="))
	}

	return finalEnvSlice
}
