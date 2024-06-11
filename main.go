package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
)

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
	const cmdName string = "go"
	args := []string{"build", ".", "fmt"}
	cmd := exec.Command(cmdName, args...)
	// assuming project a valid directory
	cmd.Dir = project
	if err := cmd.Run(); err != nil {
		return &stepErr{step: "go build", msg: "go build failed", cause: err}
	}
	_, err := fmt.Fprintln(out, "go build: SUCCESS")
	return err
}
