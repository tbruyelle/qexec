package qexec

import (
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
