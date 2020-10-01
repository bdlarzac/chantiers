/**
    Functions equivalent to php function array_unique()

    Note: depend on in-array.go
**/
package tiglib

// Equivalent of php function array_unique for a slice of strings
func ArrayUniqueString(e []string) []string {
	r := []string{}
	for _, s := range e {
		if !InArrayString(s, r[:]) {
			r = append(r, s)
		}
	}
	return r
}

// Equivalent of php function array_unique for a slice of ints
func ArrayUniqueInt(e []int) []int {
	r := []int{}
	for _, s := range e {
		if !InArrayInt(s, r[:]) {
			r = append(r, s)
		}
	}
	return r
}
