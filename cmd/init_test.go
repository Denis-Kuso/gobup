package cmd

import "testing"

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
