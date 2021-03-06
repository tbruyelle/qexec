package qexec

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/kballard/go-shellquote"
)

// Qexec holds the execution context.
type Qexec struct {
	// Vars can be used to add some new environment variable to the execution context.
	Vars map[string]string
	Out  io.Writer
}

// New returns an initialized Qexec struct.
func New() *Qexec {
	q := &Qexec{}
	q.Vars = make(map[string]string)
	q.Out = os.Stdout
	return q
}

func (q *Qexec) Run(cmds ...string) error {
	var prefix []string
	for k, v := range q.Vars {
		prefix = append(prefix, fmt.Sprintf("%s=%s", k, v))
	}
	return run(q.Out, append(prefix, cmds...)...)
}

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
	var out bytes.Buffer
	err := run(&out, cmds...)
	return out.String(), err
}

func run(out io.Writer, cmds ...string) error {
	// Wrap the command with `sh -c '<command>'`
	wrap := []string{"sh", "-c", `"`}
	wrap = append(wrap, cmds...)
	cmds = append(wrap, `"`)
	name, args, err := quote(cmds)
	if err != nil {
		return err
	}
	cmd := exec.Command(name, args...)
	cmd.Stdout = out
	cmd.Stdin = os.Stdin
	cmd.Stderr = out
	return cmd.Run()
}

// ExitStatus tries to extract the exit status from the error.
// This won't work on every platforms.
//
// If a status has been extracted from the error, then the returned
// error is null. Else the error in parameter is propagated.
func ExitStatus(err error) (int, error) {
	if err == nil {
		return 0, nil
	}
	if exiterr, ok := err.(*exec.ExitError); ok {
		// The program has exited with an exit code != 0

		// This works on both Unix and Windows. Although package
		// syscall is generally platform dependent, WaitStatus is
		// defined for both Unix and Windows and in both cases has
		// an ExitStatus() method with the same signature.
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus(), nil
		}
	}
	return 0, err
}

func quote(cmds []string) (string, []string, error) {
	toRun := strings.Join(cmds, " ")
	input, err := shellquote.Split(toRun)
	if err != nil {
		return "", nil, err
	}
	return input[0], input[1:], nil
}
