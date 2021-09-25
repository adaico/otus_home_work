package main

import (
	"log"
	"os"
)

func main() {
	envPath := os.Args[1]
	cmd := os.Args[2:]

	env, err := ReadDir(envPath)
	if err != nil {
		log.Fatal(err)
	}

	RunCmd(cmd, env)
}
