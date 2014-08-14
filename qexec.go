package qexec

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

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
	res, err := run(name, args)
	if err != nil {
		return "", err
	}
	return res, nil
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
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}
