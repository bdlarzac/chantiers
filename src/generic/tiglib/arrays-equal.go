/**
    Functions testing if two slices are equal (contain the same elements)
    Fore slices of bytes, use bytes.Equal()
**/
package tiglib

// Tests if two slices of strings contain the same elements
func ArraysEqualString(a1, a2 []string) bool {
	if len(a1) != len(a2) {
		return false
	}
	for i := range a1 {
		if a1[i] != a2[i] {
			return false
		}
	}
	return true
}

// Tests if two slices of integers contain the same elements
func ArraysEqualInt(a1, a2 []int) bool {
	if len(a1) != len(a2) {
		return false
	}
	for i := range a1 {
		if a1[i] != a2[i] {
			return false
		}
	}
	return true
}
