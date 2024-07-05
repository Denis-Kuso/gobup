package actions

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

type Step struct {
	Name string
	// command used to run the step
	cmd         string
	args        []string
	proj        string
	timeout     time.Duration
	stdoutAsErr bool
}

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

func NewStep(stepName, exe string, args []string, proj string, timeout time.Duration, stdoutAsErr bool) Step {
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
