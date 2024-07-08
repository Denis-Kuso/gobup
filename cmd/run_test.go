package cmd

import (
	"bytes"
	"strings"
	"testing"
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
func TestRunPipelines(t *testing.T) {
	type input struct {
		project  string
		pipeline string
		cfg      string
	}
	type testCase struct {
		name   string
		expErr bool
		in     input
	}
	tCases := []testCase{
		{
			name:   "normal case",
			expErr: false,
			in: input{
				project: "../testdata/testingTool",
				cfg:     returnSampleCfg(),
			},
		},
		{
			name:   "pipeline does not exist",
			expErr: true,
			in: input{
				project:  "../testdata/testingTool",
				pipeline: "madeup",
				cfg:      returnSampleCfg(),
			},
		},
		{
			name:   "step in pipeline fails",
			expErr: true,
			in: input{
				project:  "../testdata/testingToolErr",
				pipeline: "pre-commit",
				cfg:      returnSampleCfg(),
			},
		},
		{
			name:   "write to stdout treated as err",
			expErr: true,
			in: input{
				project:  "../testdata/testingToolFmtErr",
				pipeline: "pre-commit",
				cfg:      returnSampleCfg(),
			},
		},
		{
			name:   "no valid cfg",
			expErr: true,
			in: input{
				project:  "../testdata/testingTool",
				pipeline: "pre-commit",
				cfg:      "bla",
			},
		},
		{
			name:   "no pipeline to run",
			expErr: true,
			in: input{
				project: "../testdata/testingTool",
				cfg:     unrunnableCfg,
			},
		},
	}
	for _, tc := range tCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := strings.NewReader(tc.in.cfg)
			var out bytes.Buffer
			gotErr := runPipelines(tc.in.project, tc.in.pipeline, cfg, &out)
			if gotErr == nil && tc.expErr {
				t.Errorf("Expected err, got no error")
			}
			if gotErr != nil && !tc.expErr {
				t.Errorf("Expected no error, got: %v\n", gotErr)
			}
		})
	}
}

const unrunnableCfg = `pre-commit:
    run: false
    cmds:
        - build:
            cmdName: go
            args:
                - build
                - .
`

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
