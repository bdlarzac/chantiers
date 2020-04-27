package tiglib

import "testing"

func TestArraysEqualString(t *testing.T) {
	a1 := []string{"a", "b", "c"}
	a2 := []string{"a", "b", "c"}
	a3 := []string{"a", "c", "b"}
	a4 := []string{"a", "b", "c", "d"}
	a5 := []string{}
	a6 := []string{}
	if !ArraysEqualString(a1, a2) {
		t.Error(`{"a", "b", "c"} and {"a", "b", "c"} should be equal`)
	}
	if ArraysEqualString(a1, a3) {
		t.Error(`{"a", "b", "c"} and {"a", "c", "b"} should not be equal`)
	}
	if ArraysEqualString(a1, a4) {
		t.Error(`{"a", "b", "c"} and {"a", "b", "c", "d"} should not be equal`)
	}
	if ArraysEqualString(a1, a5) {
		t.Error(`{"a", "b", "c"} and {} should not be equal`)
	}
	if !ArraysEqualString(a5, a6) {
		t.Error(`{} and {} should be equal`)
	}
}

func TestArraysEqualInt(t *testing.T) {
	a1 := []int{1, 2, 3}
	a2 := []int{1, 2, 3}
	a3 := []int{1, 3, 2}
	a4 := []int{1, 2, 3, 4}
	a5 := []int{}
	a6 := []int{}
	if !ArraysEqualInt(a1, a2) {
		t.Error(`{1, 2, 3} and {1, 2, 3} should be equal`)
	}
	if ArraysEqualInt(a1, a3) {
		t.Error(`{1, 2, 3} and {1, 3, 2} should not be equal`)
	}
	if ArraysEqualInt(a1, a4) {
		t.Error(`{1, 2, 3} and {1, 2, 3, 4} should not be equal`)
	}
	if ArraysEqualInt(a1, a5) {
		t.Error(`{1, 2, 3} and {} should not be equal`)
	}
	if !ArraysEqualInt(a5, a6) {
		t.Error(`{} and {} should be equal`)
	}
}
