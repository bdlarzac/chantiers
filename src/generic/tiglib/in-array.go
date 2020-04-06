package tiglib

func InArray(elt string, array []string) bool {
	for _, test := range array {
		if elt == test {
			return true
		}
	}
	return false
}
