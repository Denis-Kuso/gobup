package actions

import (
	"fmt"
	"os/exec"
)

type Executer interface {
	Execute() (string, error)
}

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
		return "", &StepErr{
			step:  s.name,
			msg:   "failed to execute",
			cause: err}
	}
	return s.msg, nil
}

type StepErr struct {
	// which CI step caused the error
	step  string
	msg   string
	cause error
}

func (s *StepErr) Error() string {
	return fmt.Sprintf("step: %q: %s: cause: %v", s.step, s.msg, s.cause)
}

func (s *StepErr) Is(target error) bool {
	t, ok := target.(*StepErr)
	if !ok {
		return false
	}
	return t.step == s.step
}

func (s *StepErr) Unwrap() error {
	return s.cause
}
