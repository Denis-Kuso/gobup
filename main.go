package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"

)

type executer interface {
	execute() (string, error)
}

func main() {
	proj := flag.String("p", "", "Project directory")
	flag.Parse()
	if err := run(*proj, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(project string, out io.Writer) error {
	if project == "" {
		return fmt.Errorf("Project directory required: %w", ErrValidation)
	}
	const numStep int = 4// TODO refactor pipe building 
	pipe := make([]executer, numStep)
	pipe[0] = NewStep("go build", "go", []string{"build", ".", "errors"},
		"go build: SUCCESS", project)
	pipe[1] = NewStep("go test", "go", []string{"test", "-v"}, "go test: SUCCESS", project)
	pipe[2] = newObservantStep("go formating", "gofmt", []string{"-l", "."}, "gofmt: SUCCESS", project)
	var sleep time.Duration = 5 * time.Second // arbitrary decision
	pipe[3] = NewTimeoutStep("git push", "git", []string{"push"}, "git push: SUCCESS", project, sleep)

	for _, s := range pipe {
		msg, err := s.execute()
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(out, msg)
		if err != nil {
			return err
		}
	}
	return nil
}
