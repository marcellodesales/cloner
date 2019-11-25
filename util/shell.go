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

	// Make pointers to the stdout and stderr for debugging
	var stdoutBuf, stderrBuf, nodebugOutBuf, nodebugErrBuf bytes.Buffer
	stdoutIn, _ := execCmd.StdoutPipe()
	stderrIn, _ := execCmd.StderrPipe()
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)

	err := execCmd.Start()
	if err != nil {
		return "", "", errors.New(fmt.Sprintf("Can't execute command: '%v'", err))
	}

	var errStdout, errStderr error
	if IsLogInDebug() {
		// if it is in debug, show the error in the stdout
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			if IsLogInDebug() {
				_, errStdout = io.Copy(stdout, stdoutIn)
			}
			wg.Done()
		}()
		_, errStderr = io.Copy(stderr, stderrIn)
		wg.Wait()

	} else {
		_, errStdout = io.Copy(&nodebugOutBuf, stdoutIn)
		_, errStderr = io.Copy(&nodebugErrBuf, stderrIn)
	}

	err = execCmd.Wait()
	if err != nil {
		if IsLogInDebug() {
			return "", string(stderrBuf.Bytes()), errStderr
		}
		return "", string(nodebugErrBuf.Bytes()), errors.New(string(nodebugErrBuf.Bytes()))
	}

	if IsLogInDebug() {
		return string(stdoutBuf.Bytes()), "", nil
	}
	// Stdout and StdErr
	return string(nodebugOutBuf.Bytes()), "", nil
}

func ShellExecute(command string) (string, error) {
	// Execute in the background and collect the stdout
	syncStdout, syncStdErr, err := ShellExecuteAsync(command)
	if err != nil {
		return syncStdErr, err
	}
	return syncStdout, nil
}