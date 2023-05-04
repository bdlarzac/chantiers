/**
    Functions to find the index of an element in a slice.
    
    @copyright  Thierry Graff
    @license    GPL - conforms to file LICENCE located in root directory of current repository.
**/
package tiglib

// ArraySearch Returns the index of the first occurence of findMe in arr
// Returns -1 if findMe is not in arr
func ArraySearch[T comparable](arr []T, findMe T) int {
	for i, s := range arr {
		if s == findMe {
			return i
		}
	}
	return -1
}
