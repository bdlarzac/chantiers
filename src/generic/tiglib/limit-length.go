/******************************************************************************
    Given a string, computes an array of strings whose length is maximum <limit>.
    Existing newlines (\n) are preserved.
    Only lines (substrings without \n) containing more than <limit> chars are split
    Ex, limit = 20
    0---------0---------0
    This is the input string.
    A short line.
    The rest of the string
    0---------0---------0
    This is the input
    string.
    A new line.
    The rest of the
    string

    @copyright  BDL, Bois du Larzac
    @license    GPL
    @history    2020-04-07 14:39:39+02:00, Thierry Graff : Creation
********************************************************************************/
package tiglib

import (
	"strings"
	//"fmt"
)

func LimitLength(str string, limit int) []string {
	str = strings.TrimSpace(str)
	var res []string
	var curRes string
	splits := strings.Split(str, "\n")
	for _, elt := range splits {
		if len(elt) <= limit {
			res = append(res, strings.TrimSpace(elt))
			continue
		}
		words := strings.Split(elt, " ")
		curRes = ""
		for _, word := range words {
			if len(curRes)+len(word)+1 <= limit {
				curRes = curRes + " " + word
			} else {
				res = append(res, strings.TrimSpace(curRes))
				curRes = word
			}
		}
		if curRes != "" {
			res = append(res, strings.TrimSpace(curRes))
		}
	}
	return res
}
