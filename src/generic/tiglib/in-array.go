/*
*

	Functions equivalent to php function in_array().

	@copyright  Thierry Graff
	@license    GPL - conforms to file LICENCE located in root directory of current repository.

*
*/
package tiglib

// InArray returns true if elt belongs to slice arr
// Equivalent of php function in_array for a slice - generic version
func InArray[T comparable](elt T, arr []T) bool {
	for _, test := range arr {
		if elt == test {
			return true
		}
	}
	return false
}
