package actions

import (
	"errors"
	"fmt"
	"testing"

	"time"
)

// testCases
// [x]happy path - all should go well
// [x] some other command not assumed to be usually part of pipeline
// [x]command does not exists - throw err
// how about getting a root/reverse shell? or preventing one?
// [x] command takes to long - throw err
// [x] command errs to stdout (goFmt) err should be thrown
// [x] command throws err - should be "caught"
// [x] timeout called with negative (nonsensical val) is that a cfg error though?
// folder/project does not exists --> throw err
func TestRun(t *testing.T) {
	type testcase struct {
		name     string
		project  string
		expErr   error
		stepName string
		cmd      string
		args     []string
		special  bool
		timeout  time.Duration
		//expOut   string
	}
	tCases := []testcase{
		{
			name:     "happy_path",
			project:  "../../testdata/testingTool",
			expErr:   nil,
			stepName: "building",
			cmd:      "go",
			args:     []string{"build", "."},
		},
		{
			name:     "non-existing command (not installed, typo,...)",
			project:  "../../testdata/testingTool",
			expErr:   &StepErr{step: "building"},
			stepName: "building",
			cmd:      "golangtoolcommand",
			args:     []string{"build", "."},
		},
		{
			name:     "cmd takes too long",
			timeout:  2,
			project:  "../../testdata/testingTool",
			expErr:   &StepErr{step: "simulated building"},
			stepName: "simulated building",
			cmd:      "sleep",
			args:     []string{"4"},
		},
		{
			name:     "cmd called with 'negative' time", // instantiation uses default value for timeout
			timeout:  -3,
			project:  "../../testdata/testingTool",
			expErr:   nil, //&StepErr{step: "simulated building"},
			stepName: "simulated building",
			cmd:      "sleep",
			args:     []string{"2"},
		},
		{
			name:     "project does not exist",
			project:  "../testdata/testingTool",
			expErr:   &StepErr{step: "building"},
			stepName: "building",
			cmd:      "go",
			args:     []string{"build", "."},
		},
		{
			name:     "cmd's stdout treated as err",
			project:  "../../testdata/testingToolFmtErr",
			expErr:   &StepErr{step: "formating"},
			stepName: "formating",
			cmd:      "gofmt",
			args:     []string{"-l", "."},
			special:  true,
		},
		{
			name:     "actions fails",
			project:  "../../testdata/testingToolErr",
			expErr:   &StepErr{step: "building"},
			stepName: "building",
			cmd:      "go",
			args:     []string{"build", "."},
		},
		{
			name:     "simulate using some valid tool",
			project:  "../../testdata/testingTool",
			expErr:   nil,
			stepName: "listing",
			cmd:      "ls",
			args:     []string{"-la"},
			special:  true,
		},
	}
	for _, tc := range tCases {
		t.Run(tc.name, func(t *testing.T) {
			//	var out bytes.Buffer
			s := NewTimeoutStep(tc.stepName, tc.cmd, tc.args, "SUCCESS", tc.project, tc.timeout, tc.special)
			msg, gotErr := s.Execute()
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
			fmt.Printf("got msg: %s\n", msg)
		})
	}
}
