package tiglib

import "testing"

func TestInArray(t *testing.T) {
	arr := []string{"a", "b", "brt"}
	if !InArray("b", arr) {
		t.Error("b should belong to {a, b, brt}")
	}
	if InArray("r", arr) {
		t.Error("r should not belong to {a, b, brt}")
	}
}
