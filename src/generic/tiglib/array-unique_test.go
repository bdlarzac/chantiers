/**
    Note: tests of array-unique.go depend on array-equal.go
**/

package tiglib

import "testing"

func TestArrayUniqueString(t *testing.T) {
	a1 := []string{"a", "a", "b"}
	a2 := []string{"a", "b"}
	a3 := []string{}
	uniq1 := ArrayUniqueString(a1)
	uniq2 := ArrayUniqueString(a2)
	uniq3 := ArrayUniqueString(a3)
	if !ArraysEqualString(uniq1, a2) {
		t.Error(`ArrayUniqueString({"a", "a", "b"}) should be = to {"a", "b"}`)
	}
	if !ArraysEqualString(uniq2, a2) {
		t.Error(`ArrayUniqueString({"a", "b"}) should be = to {"a", "b"}`)
	}
	if !ArraysEqualString(uniq3, a3) {
		t.Error(`ArrayUniqueString({}) should be = to {}`)
	}
}

func TestArrayUniqueInt(t *testing.T) {
	a1 := []int{1, 1, 2}
	a2 := []int{1, 2}
	a3 := []int{}
	uniq1 := ArrayUniqueInt(a1)
	uniq2 := ArrayUniqueInt(a2)
	uniq3 := ArrayUniqueInt(a3)
	if !ArraysEqualInt(uniq1, a2) {
		t.Error(`ArrayUniqueInt({1, 1, 2}) should be = to {1, 2}`)
	}
	if !ArraysEqualInt(uniq2, a2) {
		t.Error(`ArrayUniqueInt({1, 2}) should be = to {1, 2}`)
	}
	if !ArraysEqualInt(uniq3, a3) {
		t.Error(`ArrayUniqueInt({}) should be = to {}`)
	}
}
