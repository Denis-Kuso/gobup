package main

import (
	"errors"
	"io"
	"strings"
	"testing"
)

func TestLoadCfg(t *testing.T) {
	// able to read file but not unmarshall according to structure of cfg
	// unable to read (but compiles)
	// read a valid cfg "file"
	type testcase struct {
		name   string
		expErr error
		input  io.Reader
	}
	validCfg := strings.NewReader(`# pipeline name
pre-commit:
  run: true
  # if possible, all commands run as warnings
  fail_fast: false
  # sequence of commands to run in this pipeline
  cmds:
    - build:
        cmdName: go build
        # args are ordered
        args:
          - one
          - two
          - three
    - test:
        cmdName: go test
        args:
          - "."
    - format:
        cmdName: gofmt
        args: ["-l", "."]`)
	validReader := strings.NewReader("Perfectly valid reader, but not valid cfg")
	noData := strings.NewReader("")
	tCases := []testcase{
		{
			name:   "valid cfg",
			expErr: nil,
			input:  validCfg,
		},
		{
			name:   "invalid cfg format",
			expErr: ErrConfig,
			input:  validReader,
		},
		{
			name:   "no input",
			expErr: ErrConfig,
			input:  noData,
		},
	}
	for _, tc := range tCases {
		t.Run(tc.name, func(t *testing.T) {
			_, gotErr := LoadCfg(tc.input)
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
		})
	}
}
