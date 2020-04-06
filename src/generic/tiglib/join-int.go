package tiglib

import (
	"strconv"
	"strings"
)

/**
    Same as strings.Join, for a slice of int
**/
func JoinInt(arr []int, sep string) string {
	arr2 := []string{}
	for _, elt := range arr {
		arr2 = append(arr2, strconv.Itoa(elt))
	}
	return strings.Join(arr2, sep)
}
