/**
    Functions equivalent to php function in_array()

    Note: functions depending on in-array.go
**/
package tiglib

// Equivalent of php function in_array for a slice of strings
func InArrayString(elt string, array []string) bool {
	for _, test := range array {
		if elt == test {
			return true
		}
	}
	return false
}

// Equivalent of php function in_array for a slice of ints
func InArrayInt(elt int, array []int) bool {
	for _, test := range array {
		if elt == test {
			return true
		}
	}
	return false
}
