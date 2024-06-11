package main

import (
	"os/exec"
)

// A "continous-integration" step
type step struct {
	name string
	// command used to run the step
	cmd  string
	args []string
	// msg is printed out, not used for errors
	msg  string
	proj string
}

func NewStep(name, cmd string, args []string, msg, proj string) step {
	return step{
		name: name,
		cmd:  cmd,
		args: args,
		msg:  msg,
		proj: proj,
	}
}

func (s step) execute() (string, error) {
	cmd := exec.Command(s.cmd, s.args...)
	cmd.Dir = s.proj
	if err := cmd.Run(); err != nil {
		return "", &stepErr{
			step:  s.name,
			msg:   "failed to execute",
			cause: err}
	}
	return s.msg, nil
}
