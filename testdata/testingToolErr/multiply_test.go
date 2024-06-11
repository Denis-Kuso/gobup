package multiply

import (
	"testing"
)

func TestAdd(t *testing.T) {
	a := 3
	b := 3
	exp := 9
	got := mult(a, b)
	if exp != got {
		t.Errorf("Expected %d, got %d.", exp, got)
	}
}
