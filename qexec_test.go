package qexec

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunSuccess(t *testing.T) {
	cmd := "ls testdata"

	output, err := Run(cmd)

	assert.Nil(t, err)
	assert.Equal(t, "truc\n", output)
}

func TestRunFailed(t *testing.T) {
	cmd := "ls testdata/notexist"

	output, err := Run(cmd)

	assert.NotNil(t, err)
	assert.Equal(t, "ls: testdata/notexist: No such file or directory\n", output)
}

func TestRunMixed(t *testing.T) {
	cmd := "ls testdata testdata/notexist"

	output, err := Run(cmd)

	assert.NotNil(t, err)
	assert.Equal(t, "ls: testdata/notexist: No such file or directory\ntestdata:\ntruc\n", output)
}

func TestRunWithVar(t *testing.T) {
	cmd := "echo $PWD"
	exp, err := Run("pwd")
	if err != nil {
		t.Fatalf("Unable to prepare test : %s", err)
	}

	output, err := Run(cmd)

	assert.Nil(t, err)
	assert.Equal(t, exp, output)
}

func TestExitStatusSuccess(t *testing.T) {
	cmd := "ls testdata"
	_, err := Run(cmd)

	status, err := ExitStatus(err)

	assert.Equal(t, 0, status)
	assert.Nil(t, err)
}

func TestExitStatusFailed(t *testing.T) {
	cmd := "ls testdata/notexists"
	_, err := Run(cmd)

	status, err := ExitStatus(err)

	assert.NotEqual(t, 0, status)
	assert.Nil(t, err)
}

func TestExitStatusError(t *testing.T) {
	err := errors.New("I'm not returned by cmd.Run")

	status, err := ExitStatus(err)

	assert.Equal(t, 0, status)
	assert.NotNil(t, err)
}

func TestQexecRun(t *testing.T) {
	q := New()
	q.Vars["MY_VAR"] = "value"

	output, err := q.Run("sh -c 'echo $MY_VAR'")

	assert.Nil(t, err)
	assert.Equal(t, "value\n", output)
}
