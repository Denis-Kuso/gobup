package cmd

import (
	"bytes"
	"testing"
)

func TestCreateTemplate(t *testing.T) {
	type testCase struct {
		name   string
		in     string
		expErr bool
		// TODO add a cleanup function?
	}
	tCases := []testCase{
		{
			// currently assumes there is no valid config file
			name:   "validProject - relative path",
			expErr: false,
			in:     ".",
		},
		{
			name:   "invalid filepath",
			expErr: true,
			in:     "&znj",
		},
		{
			name:   "file exists",
			expErr: true,
			in:     ".",
		},
	}
	for _, tc := range tCases {
		t.Run(tc.name, func(t *testing.T) {
			gotErr := createTemplate(tc.in)
			if gotErr != nil && !tc.expErr {
				t.Errorf("expected no error, got %q\n", gotErr)
				return
			}
			if gotErr == nil && tc.expErr {
				t.Errorf("expected err, got nil\n")
				return
			}
		})
	}
}

func TestWriteCmd(t *testing.T) {

	type testCase struct {
		name   string
		in     string
		expErr bool
	}
	tCases := []testCase{
		{
			name:   "happy",
			in:     "hello world",
			expErr: false,
		},
		{
			name:   "empty",
			in:     "",
			expErr: true,
		},
		{
			name:   "actual command",
			in:     "gobup run -a pre-commit",
			expErr: false,
		},
	}
	for _, tc := range tCases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			gotErr := writeCmd([]byte(tc.in), &buf)
			if gotErr != nil && !tc.expErr {
				t.Errorf("expected no error, got %q\n", gotErr)
				return
			}
			if gotErr == nil && tc.expErr {
				t.Errorf("expected err, got nil\n")
				return
			}
		})
	}
}
