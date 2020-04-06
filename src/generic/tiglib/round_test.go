package tiglib

import (
	"strconv"
	"testing"
)

func TestRound(t *testing.T) {
	var tests = []struct {
		input float64
		round int
		want  float64
	}{
		{1.12345678, 2, 1.12},
		{1.129, 2, 1.13},
		{1.12345678, 3, 1.123},
	}
	for _, test := range tests {
		if got := Round(test.input, test.round); got != test.want {
			t.Error("Expect " + strconv.FormatFloat(test.want, 'b', -1, 64) + ", got " + strconv.FormatFloat(got, 'b', -1, 64))
		}
	}
}
