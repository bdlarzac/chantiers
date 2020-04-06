package tiglib

import "testing"

func TestJoinInt(t *testing.T) {
	arr := []int{1, 2, 3}
	test := JoinInt(arr, " ")
	if test != "1 2 3" {
		t.Error("Expect \"1 2 3\", got \"" + test + "\"")
	}
	test = JoinInt(arr, "Z")
	if test != "1Z2Z3" {
		t.Error("Expect \"1Z2Z3\", got \"" + test + "\"")
	}
}
