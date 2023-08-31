package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	run := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	run.Stdin = os.Stdin
	run.Stdout = os.Stdout
	run.Stderr = os.Stderr
	for name, envVal := range env {
		if envVal.NeedRemove {
			err := os.Unsetenv(name)
			if err != nil {
				log.Fatal(fmt.Errorf("unable to unset env: %w", err))
			}
			continue
		}
		os.Setenv(name, envVal.Value)
	}

	if err := run.Start(); err != nil {
		log.Fatal(err)
	}
	if err := run.Wait(); err != nil {
		log.Fatal(err)
	}

	return run.ProcessState.ExitCode()
}
