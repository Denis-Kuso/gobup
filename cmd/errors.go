package cmd

import "errors"

// ErrValidation is an error value associated with the correct state of
// provided arguments and/or usage of the gobup cmd.
var ErrValidation = errors.New("validation failure")
