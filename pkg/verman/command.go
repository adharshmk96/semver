package verman

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

func Clean(output string, err error) (string, error) {
	output = strings.ReplaceAll(strings.Split(output, "\n")[0], "'", "")
	if err != nil {
		err = errors.New(strings.TrimSuffix(err.Error(), "\n"))
	}
	return output, err
}

func RunCmd(name string, args ...string) (string, error) {
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}

	c := exec.Command(name, args...)
	c.Stdout = &stdout
	c.Stderr = &stderr

	err := c.Run()
	if err != nil {
		return "", errors.New(stderr.String())
	}
	return Clean(stdout.String(), nil)
}
