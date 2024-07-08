// Package actions defines the Step type and its error along with the Execute
// method. Together, these enable enviroment setup for the step to execute in.
package actions

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

// Step represents information neccessary to prepare an enviroment for an
// executable command (e.g. gofmt, go build) to run
type Step struct {
	Name string
	// command used to run the step
	cmd         string
	args        []string
	proj        string
	timeout     time.Duration
	stdoutAsErr bool
}

// StepErr represent an error type caused by a Step, containing information
// about its cause.
type StepErr struct {
	// which CI step caused the error
	step  string
	msg   string
	cause error
}

func (s *StepErr) Error() string {
	if s.cause == nil {
		return fmt.Sprintf("step: %q: %s", s.step, s.msg)
	}
	return fmt.Sprintf("step: %q\n %scause: %s", s.step, s.msg, s.cause)
}

// Is returns true if any of target's type in its error chain equals to StepErr
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

// NewStep creates a new Step with provided arguments. Step is used primarily,
// to set up the environment for a command to execute. If timeout t is provided
// as t <= 0, default value of 30s is used.
func NewStep(stepName, exe string, args []string, proj string, timeout time.Duration, stdoutAsErr bool) Step {
	// arbitrary choice of 30s, most commands used complete well below 30s
	const defaultTimeout time.Duration = 30
	t := Step{
		Name:        stepName,
		cmd:         exe,
		args:        args,
		proj:        proj,
		timeout:     timeout,
		stdoutAsErr: stdoutAsErr,
	}

	if timeout <= 0 {
		t.timeout = defaultTimeout * time.Second
	} else {
		t.timeout = timeout * time.Second
	}
	return t
}

// Execute executes the cmd associated with step s. On success it return nil.
// Returns an error:
// - if default or provided timeout exceeded
// - if s.stdoutAsErr is true and command wrote to stdout
// - the cmd executes unsuccessfuly
func (s Step) Execute() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, s.cmd, s.args...)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	cmd.Dir = s.proj
	if err := cmd.Run(); err != nil {
		if err == context.DeadlineExceeded {
			m := fmt.Sprintf("timeout: %d exceeded", s.timeout)
			return &StepErr{
				step:  s.Name,
				msg:   m,
				cause: context.DeadlineExceeded}
		}
		return &StepErr{
			step:  s.Name,
			msg:   fmt.Sprintf("%s %s", out.String(), stderr.String()),
			cause: err}
	}
	// tools (e.g gofmt) which return 0 on success, but some "error" msg is returned
	// to stdout
	if s.stdoutAsErr && (out.Len() > 0) {
		return &StepErr{
			step:  s.Name,
			msg:   fmt.Sprintf("%s", out.String()),
			cause: nil, // no underlying cause, apart from our interpretation of output as an err
		}
	}
	return nil
}
