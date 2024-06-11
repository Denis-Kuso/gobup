package main

import (
	"bytes"
	"context"
	"os/exec"
	"time"
	"fmt"
)

type timeoutStep struct {
	step
	timeout time.Duration
}

func NewTimeoutStep(name, exe string, args []string, message, proj string, 
	timeout time.Duration) timeoutStep {
	const defaultTimeout time.Duration = 30 
	t := timeoutStep{}
	t.step = NewStep(name, exe, args, message, proj)
	if timeout == 0 {
		t.timeout = defaultTimeout * time.Second
	}else {
		t.timeout = timeout;
	}
	return t
}

func (t timeoutStep) execute() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, t.cmd, t.args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Dir = t.proj
	if err := cmd.Run(); err != nil {
		if err == context.DeadlineExceeded {
			m := fmt.Sprintf("failed due to timeout, set to: %v", t.timeout)
			return "", &stepErr{
			step: t.name, 
			msg: m,
			cause: context.DeadlineExceeded, }
		}
		return "", &stepErr{
			step: t.name, 
			msg: "failed executing",
			cause: err, }
	}
	return t.msg, nil
}
