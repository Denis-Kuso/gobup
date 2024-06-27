package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"github.com/Denis-Kuso/gobup/internal/actions"
)

func TestRun(t *testing.T) {
	type testcase struct {
		name     string
		project  string
		expErr   error
		expOut   string
		setupGit bool
	}
	tCases := []testcase{
		{
			name:     "succesful build",
			project:  "./testdata/testingTool",
			expErr:   nil,
			expOut:   "go build: SUCCESS\ngo test: SUCCESS\ngofmt: SUCCESS\ngit push: SUCCESS\n",
			setupGit: true,
		},
		{
			name:     "failed build",
			project:  "./testdata/testingToolErr",
			expErr:   &actions.StepErr{step: "go build"},
			expOut:   "",
			setupGit: false,
		},
		{
			name:     "failed formating",
			project:  "./testdata/testingToolFmtErr",
			expErr:   &actions.StepErr{step: "go formating"},
			expOut:   "",
			setupGit: false,
		},
	}

	for _, tc := range tCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupGit {
				cleanup := setupGit(t, tc.project)
				defer cleanup()
			}
			var out bytes.Buffer
			gotErr := run(tc.project, &out)
			if tc.expErr != nil {
				if gotErr == nil {
					t.Errorf("Expected error: %q. Got nil instead", tc.expErr)
					return
				}
				if !errors.Is(gotErr, tc.expErr) {
					t.Errorf("Error types differ. Expected: %q, got :%q instead", tc.expErr, gotErr)
				}
				return
			}
			if gotErr != nil {
				t.Errorf("Expected no error, got: %q", gotErr)
				return
			}
			gotOut := out.String()
			if tc.expOut != gotOut {
				t.Errorf("Expected output: %s, recevied: %s", tc.expOut, gotOut)
			}
		})
	}
}

// Helper function that
// creates a temporary directory.
// creates a Bare git repository on this temporary directory.
// initializes a git repository on the target project directory.
// add the bare git repository as a remote repository in the empty Git repository in the target project directory.
// stages a file to commit.
// commits the changes to the Git repository
// returns a function that deletes the temp dir created
func setupGit(t *testing.T, proj string) func() {
	t.Helper()
	// fail if absent
	gitExec, err := exec.LookPath("git")
	if err != nil {
		t.Fatalf("git not found: %s", err)
	}
	tempDir, err := os.MkdirTemp("", "gobup")
	if err != nil {
		t.Fatalf("could not create directory: %s", tempDir)
	}
	projPath, err := filepath.Abs(proj)
	if err != nil {
		t.Fatalf("cannot create pathname: %s", proj)
	}
	remoteRepo := fmt.Sprintf("file://%s", tempDir)
	gitCmds := []struct {
		dir  string
		args []string
		env  []string
	}{
		{tempDir, []string{"init", "--bare"}, nil},
		{projPath, []string{"init"}, nil},
		{projPath, []string{"remote", "add", "origin", remoteRepo}, nil},
		{projPath, []string{"add", "."}, nil},
		{projPath, []string{"commit", "-m", "testRepo"},
			[]string{"GIT_COMMITTER_NAME=tester",
				"GIT_COMMITTER_EMAIL=tester@example.com",
				"GIT_AUTHOR_NAME=tester",
				"GIT_AUTHOR_EMAIL=tester@example.com"},
		},
	}
	for _, gc := range gitCmds {
		gitCmd := exec.Command(gitExec, gc.args...)
		gitCmd.Dir = gc.dir
		if gc.env != nil {
			gitCmd.Env = append(os.Environ(), gc.env...)
		}
		if err := gitCmd.Run(); err != nil {
			t.Fatalf("failed executing gitcmd: %s, err:%q", gc.args, err)
		}
	}
	return func() {
		os.RemoveAll(tempDir)
		os.RemoveAll(filepath.Join(projPath, ".git"))
	}
}
