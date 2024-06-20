package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	// "time"
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
	const capacity int = 1
	manager := make(chan os.Signal, capacity)
	faults := make(chan error)
	done := make(chan struct{})
	signal.Notify(manager, syscall.SIGINT, syscall.SIGTERM)
	//	const numStep int = 4 // TODO refactor pipe building
	//	pipe := make([]executer, numStep)
	//	pipe[0] = NewStep("go build", "go", []string{"build", ".", "errors"},
	//		"go build: SUCCESS", project)
	//	pipe[1] = NewStep("go test", "go", []string{"test", "-v"}, "go test: SUCCESS", project)
	//	pipe[2] = newObservantStep("go formating", "gofmt", []string{"-l", "."}, "gofmt: SUCCESS", project)
	//	var sleep time.Duration = 5 * time.Second // arbitrary decision
	//	pipe[3] = NewTimeoutStep("git push", "git", []string{"push", "origin", "master"}, "git push: SUCCESS", project, sleep)
	c, err := loadCfg()
	if err != nil {
		fmt.Printf("Config err: %v", err)
		os.Exit(-1)
	}
	pipes := getPipelines(c)
	pre_commit, ok := pipes["pre-commit"]
	if !ok {
		fmt.Println("no such pipe")
		os.Exit(22)
	}
	pipe := makePipe(pre_commit.Steps, project)
	go func() {
		for _, s := range pipe {
			msg, err := s.execute()
			if err != nil {
				faults <- err
				return
			}
			_, err = fmt.Fprintln(out, msg)
			if err != nil {
				faults <- err
				return
			}
		}
		close(done)
	}()
	for {
		select {
		case got := <-manager:
			signal.Stop(manager)
			return fmt.Errorf("\nStoping pipeline due to: %s", got)
		case err := <-faults:
			//fmt.Printf("poped %s from err chanel\n", err)
			return err
		case <-done:
			return nil
		}
	}
}
