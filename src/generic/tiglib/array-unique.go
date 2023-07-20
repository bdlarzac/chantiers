/*
Function equivalent to php function array_unique().

Note: depends on in-array.go

@copyright  Thierry Graff
@license    GPL - conforms to file LICENCE located in root directory of current repository.
*/
package tiglib

// ArrayUnique is an equivalent of php function array_unique - generic version
func ArrayUnique[T comparable](e []T) []T {
	r := []T{}
	for _, s := range e {
		if !InArray(s, r[:]) {
			r = append(r, s)
		}
	}
	return r
}
