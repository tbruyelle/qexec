package qexec

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/kballard/go-shellquote"
)

// Run executes the command in parameter after having correctly quoted it.
// The command stdout is returned.
//
// It handles a common error when the path to the executable contains one or more
// environment variable, which usually produces an error `no such file
// or directory`. This is because `os/exec` checks the existence of the
// executable and it doesn't interpret the environment variables.
// Here if the executable contains any $ char, then the whole command is
// wrapped by `sh -c "<command>"`.
func Run(cmds ...string) (string, error) {
	if strings.Contains(cmds[0], "$") {
		// If the path to the executable contains env variables,
		// then the command must be wrapped by `sh -c "<command>"`
		wrap := []string{"sh", "-c", `"`}
		wrap = append(wrap, cmds...)
		wrap = append(wrap, `"`)
		cmds = wrap
	}
	name, args, err := quote(cmds)
	if err != nil {
		return "", err
	}
	return run(name, args)
}

// ExitStatus tries to extract the exit status from the error.
// This won't work on every platforms.
func ExitStatus(err error) int {
	if err == nil {
		return 0
	}
	if exiterr, ok := err.(*exec.ExitError); ok {
		// The program has exited with an exit code != 0

		// This works on both Unix and Windows. Although package
		// syscall is generally platform dependent, WaitStatus is
		// defined for both Unix and Windows and in both cases has
		// an ExitStatus() method with the same signature.
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus()
		}
	}
	return 0
}

func quote(cmds []string) (string, []string, error) {
	toRun := strings.Join(cmds, " ")
	input, err := shellquote.Split(toRun)
	if err != nil {
		return "", nil, err
	}
	return input[0], input[1:], nil
}

func run(name string, args []string) (string, error) {
	cmd := exec.Command(name, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stdin = os.Stdin
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}
