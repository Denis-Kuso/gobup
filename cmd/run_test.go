package cmd

import (
	"strings"
	"testing"
	//"github.com/Denis-Kuso/gobup/internal/config"
)

// func preparePipes(cfg io.Reader, pipeline string) ([]config.Action, error) {
func TestPreparePipes(t *testing.T) {
	// both nil - expect error
	// one is nil expect err
	// pipeline not present in cfg
	type testCase struct {
		name        string
		expErr      bool
		inputReader string
		inPipe      string
	}
	tCases := []testCase{
		{
			name:        "invalid config provided",
			expErr:      true,
			inputReader: "frodoBaggins",
			inPipe:      "pre-commit",
		},
		{
			name:        "pipeline does not exist",
			expErr:      true,
			inputReader: returnSampleCfg(),
			inPipe:      "madeup-p1p3",
		},
		{
			name:        "nil cfg provided",
			expErr:      true,
			inputReader: "",
			inPipe:      "",
		},
		{
			name:        "expected/normal case",
			expErr:      false,
			inputReader: returnSampleCfg(),
			inPipe:      "pre-commit",
		},
	}
	for _, tc := range tCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := strings.NewReader(tc.inputReader)
			_, gotErr := preparePipes(cfg, tc.inPipe)
			if tc.expErr && (gotErr == nil) {
				t.Errorf("Expected err, recevied nil")
				return
			}
			if gotErr != nil && !tc.expErr {
				t.Errorf("Expected no error, got %v\n", gotErr)
				return
			}
		})
	}
}
func returnSampleCfg() string {
	return `pre-commit:
    run: true
    cmds:
        - build:
            cmdName: go
            args:
                - build
                - .
        - test:
            cmdName: go
            args:
                - test
                - .
        - format:
            cmdName: gofmt
            args:
                - -l
                - .
            stdoutAsErr: true`
}
