package main

import (
	"fmt"
	"io"
	"os/exec"
	"os"
	"flag"
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
		return fmt.Errorf("no project name provided")
	}
	const cmdName string = "go"
	args := []string{"build", ".", "fmt"}
	cmd := exec.Command(cmdName, args...)
	// assuming project a valid directory
	cmd.Dir = project
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("'%s build' failed: %v", cmdName, err)
	}
	_, err := fmt.Fprintln(out, "Go Build: SUCCESS")
	return err
}
