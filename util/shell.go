package util

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

/**
 * Executes the shell and wait until it finishes executing
 * The result or error is provided
 * https://stackoverflow.com/questions/1877045/how-do-you-get-the-output-of-a-system-command-in-go/54586179#54586179
 */
func ShellExecuteSync(command string) (string, error) {
	commandArgs := strings.Split(command, " ")
	baseCmd := commandArgs[0]
	cmdArgs := commandArgs[1:]

	var execCmd *exec.Cmd
	if strings.Contains(command, "|") {
		// https://stackoverflow.com/questions/10781516/how-to-pipe-several-commands-in-go/30329351#30329351
		execCmd = exec.Command("bash","-c", baseCmd + " " + strings.Join(cmdArgs, " "))

	} else {
		execCmd = exec.Command(baseCmd, cmdArgs...)
	}
	out, err := execCmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

/**
 * Executes the shell and show the live results, returning the
 * output and error.
 * https://hackmongo.com/post/reading-os-exec-cmd-output-without-race-conditions/
 */
func ShellExecuteAsync(command string) (string, string, error) {
	commandArgs := strings.Split(command, " ")
	baseCmd := commandArgs[0]
	cmdArgs := commandArgs[1:]

	var execCmd *exec.Cmd
	if strings.Contains(command, "|") {
		// https://stackoverflow.com/questions/10781516/how-to-pipe-several-commands-in-go/30329351#30329351
		execCmd = exec.Command("bash","-c", baseCmd + " " + strings.Join(cmdArgs, " "))

	} else {
		execCmd = exec.Command(baseCmd, cmdArgs...)
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	stdoutIn, _ := execCmd.StdoutPipe()
	stderrIn, _ := execCmd.StderrPipe()

	var errStdout, errStderr error
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)
	err := execCmd.Start()
	if err != nil {
		return "", "", errors.New(fmt.Sprintf("Can't execute command: '%v'", err))
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
		wg.Done()
	}()
	_, errStderr = io.Copy(stderr, stderrIn)
	wg.Wait()

	err = execCmd.Wait()
	if err != nil {
		return "", "", errors.New(fmt.Sprintf("Can't wait for execution: '%v'", err))
	}
	if errStdout != nil || errStderr != nil {
		return "", "", errors.New(fmt.Sprintf("Can't capture stdout or stderr"))
	}

	// Stdout and StdErr
	return string(stdoutBuf.Bytes()), string(stderrBuf.Bytes()), nil
}

func ShellExecute(command string) (string, error) {
	var stdout string
	if IsLogInDebug() {
		// Execute in the background and collect the stdout
		syncStdout, _, err := ShellExecuteAsync(command)
		if err != nil {
			return "", err
		}
		stdout = syncStdout

	} else {
		// Executes async, show the docker execution logs and collect the log lines
		syncStdout, err := ShellExecuteSync(command)
		if err != nil {
			return "", err
		}
		stdout = syncStdout
	}
	return stdout, nil
}