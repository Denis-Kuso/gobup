package main

import (
	"bytes"
	"errors"
	"testing"
)

func TestRun(t *testing.T) {
	type testcase struct {
		name    string
		project string
		expErr  error
		expOut  string
	}
	tCases := []testcase{
		{
			name:    "succesful build",
			project: "./testdata/testingTool",
			expErr:  nil,
			expOut:  "go build: SUCCESS\ngo test: SUCCESS\ngofmt: SUCCESS\n",
		},
		{
			name:    "failed build",
			project: "./testdata/testingToolErr",
			expErr:  &stepErr{step: "go build"},
			expOut:  "",
		},
		{
			name:    "failed formating",
			project: "./testdata/testingToolFmtErr",
			expErr:  &stepErr{step: "go formating"},
			expOut:  "",
		},
	}

	for _, tc := range tCases {
		t.Run(tc.name, func(t *testing.T) {
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
