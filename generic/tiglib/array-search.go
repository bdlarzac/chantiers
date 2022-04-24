/**
    Functions the index of an element in a slice
**/
package tiglib

// Returns the index of the first occurence of findMe in arr
// Returns -1 if findMe is not in arr
func ArraySearchString(arr []string, findMe string) int {
	for i, s := range arr {
		if s == findMe {
			return i
		}
	}
	return -1
}
