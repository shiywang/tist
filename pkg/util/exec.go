package util

import (
	"bytes"
	"os/exec"
)

type Execer interface {
	Exec(s string, args ...string) (string, error)
}

type BashExec struct {
	out bytes.Buffer
}

func (b BashExec) Exec(s string, args ...string) (string, error) {
	cmd := exec.Command(s, args...)
	cmd.Stdout = &b.out
	cmd.Stderr = &b.out
	err := cmd.Run()
	return b.out.String(), err
}
