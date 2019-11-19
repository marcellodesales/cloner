package util

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
	"strings"
)

// Docker Executor struct to show
type DockerExecutor struct {
	// The docker command to execute, any command processed using Sprintf
	Command string
	// The arguments to be passed to the format string
	Arguments []interface{}
	// The actual command to execute processed during execution
	RawCommand string

	// TODO: Request return type to commands
	// * json or yaml so that we can process them properly instead of interpreting stdout
}

/**
 * Factory method that creates a new instance of the Docker Executor
 */
func NewDockerExecutor(command string, args []string) *DockerExecutor {
	// https://stackoverflow.com/questions/12990338/cannot-convert-string-to-interface/12990540#12990540
	// Convert the array of strings into array of interfaces
	dockerArgs := formatCommandArguments(args)

	// Add the env var to the docker run command when it's in debug mode
	if IsLogInDebug() && funk.Contains(command, "--rm") {
		command = strings.ReplaceAll(command, "--rm", "--rm -e PROTOCOOL_DEBUG=1")
	}

	// Create the instance of the executor
	dockerExecutor := DockerExecutor{
		Command:   command,
		Arguments: dockerArgs,
		RawCommand: fmt.Sprintf(command, dockerArgs...), // https://stackoverflow.com/questions/7145905/fmt-sprintf-passing-an-array-of-arguments/7153343#7153343
	}

	// Return the reference
	return &dockerExecutor
}

/**
 * Format the arguments
 */
func formatCommandArguments(args []string) []interface{} {
	commandArgs := make([]interface{}, len(args))
	for i, v := range args {
		commandArgs[i] = v
	}
	return commandArgs
}

/**
 * Executes a given Docker Executor reference returning the (standout, error)
 */
func (dockerExecutor *DockerExecutor) Execute() (string, error) {
	log.Debugf("Executing Docker Command: %s", dockerExecutor.RawCommand)

	stdout, err := ShellExecute(dockerExecutor.RawCommand)
	if err != nil {
		log.Debug("Failed to execute the docker command...")
		log.Errorf("Couldn't execute docker command: %v", err)
		return "", err

	} else {
		log.Debug("Result from the docker command execution was successful")
		return stdout, nil
	}
}