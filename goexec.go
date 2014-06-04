package qexec

import (
	"bytes"
	"github.com/kballard/go-shellquote"
	"os"
	"os/exec"
	"strings"
)

// Run executes the command in parameter after having correctly quoted it.
// The command stdout is returned.
func Run(cmds ...string) (string, error) {
	toRun := strings.Join(cmds, " ")
	input, err := shellquote.Split(toRun)
	if err != nil {
		return "", err
	}
	name := input[0]
	arg := input[1:]
	cmd := exec.Command(name, arg...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}
