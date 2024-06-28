package actions

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

type timeoutStep struct {
	step
	timeout     time.Duration
	stdoutAsErr bool
}

func NewTimeoutStep(name, exe string, args []string, message, proj string,
	timeout time.Duration, stdoutAsErr bool) timeoutStep {
	const defaultTimeout time.Duration = 30
	t := timeoutStep{}
	t.step = NewStep(name, exe, args, message, proj)
	if timeout <= 0 {
		t.timeout = defaultTimeout * time.Second
	} else {
		t.timeout = timeout * time.Second
	}
	t.stdoutAsErr = stdoutAsErr
	return t
}

func (t timeoutStep) Execute() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, t.cmd, t.args...)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	cmd.Dir = t.proj
	if err := cmd.Run(); err != nil {
		if err == context.DeadlineExceeded {
			m := fmt.Sprintf("failed due to timeout, set to: %v", t.timeout)
			return "", &StepErr{
				step:  t.name,
				msg:   m,
				cause: context.DeadlineExceeded}
		}
		return "", &StepErr{
			step:  t.name,
			msg:   fmt.Sprintf("failed executing: out:\n %s: err: \n%s", out.String(), stderr.String()),
			cause: err}
	}
	// tools (e.g gofmt) which return 0 on success, but some "error" msg is returned
	// to stdout
	if t.stdoutAsErr && (out.Len() > 0) {
		return "", &StepErr{
			step:  t.name,
			msg:   fmt.Sprintf("%s", out.String()),
			cause: nil, // no underlying cause, apart from our interpretation of output as an err
		}
	}
	return t.msg, nil
}
