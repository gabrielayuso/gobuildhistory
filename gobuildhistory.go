package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

type CmdError struct {
	cmd string
	msg string
}

func (e *CmdError) Error() string {
	return fmt.Sprintf("'%s' command failed. | %s", e.cmd, e.msg)
}

func execCommand(command string) (out string, err *CmdError) {
	cmdParts := strings.Split(command, " ")
	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)

	cmdOut, _ := cmd.StdoutPipe()
	cmdErr, _ := cmd.StderrPipe()
	cmd.Start()
	cmdOutBytes, _ := ioutil.ReadAll(cmdOut)
	cmdErrBytes, _ := ioutil.ReadAll(cmdErr)
	cmd.Wait()

	if len(cmdErrBytes) > 0 {
		return "", &CmdError{command, string(cmdErrBytes)}
	}

	return string(cmdOutBytes), nil
}

func getCommits() ([]string, error) {
	out, err := execCommand("git log --pretty=format:%h")

	if err != nil {
		return nil, err
	}

	return strings.Split(out, "\n"), nil
}

func checkoutCommit(commit string) error {
	if _, err := execCommand("git checkout " + commit); err != nil {
		return err
	}

	return nil
}

func goBuild() error {
	if _, err := execCommand("go build"); err != nil {
		return err
	}

	return nil
}

func main() {
	commits, err := getCommits()

	if err != nil {
		panic(err)
	}

	for i := len(commits) - 1; i >= 0; i-- {
		commit := commits[i]
		checkoutCommit(commit)
		msg := "Successfully built project at commit: " + commit

		if err := goBuild(); err != nil {
			msg = "!! Failed to build project at commit: " + commit + " Reason: " + err.Error()
		}

		fmt.Println(msg)
	}
}
