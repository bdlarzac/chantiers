package tiglib

import "testing"

func TestArraySearchString(t *testing.T) {
	a := []string{"a", "b", "c", "d"}
	if ArraySearchString(a, "c") != 2 {
		t.Error(`ArraySearchString([]string{"a", "b", "c", "d"}, "c") should be = 2`)
	}
}
