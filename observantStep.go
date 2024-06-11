package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

// A "continous-integration" step that observes the output of the step
type observantStep struct {
	step
}

func newObservantStep(name, exe string, args []string, message, proj string) observantStep {
	s := observantStep{}
	s.step = NewStep(name, exe, args, message, proj)
	return s
}

func (o observantStep) execute() (string, error) {
	cmd := exec.Command(o.cmd, o.args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Dir = o.proj
	if err := cmd.Run(); err != nil {
		return "", &stepErr{
			step:  o.name,
			msg:   "failed to execute",
			cause: err,
		}
	}
	if out.Len() > 0 {
		return "", &stepErr{
			step: o.name,
			msg:  fmt.Sprintf("poorly formated files: %s", out.String()),
			cause: nil, // gofmt is based on opinion/conventions
		}
	}
	return o.msg, nil
}
