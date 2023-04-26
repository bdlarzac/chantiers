/**
    Reverse a slice in-place
    From https://stackoverflow.com/questions/13190836/is-there-a-way-to-iterate-over-a-slice-in-reverse-in-go
    (solution of Ivan Voras)

    @copyright  Thierry Graff
    @license    GPL - conforms to file LICENCE located in root directory of current repository.
**/
package tiglib

func ArrayReverse[S ~[]E, E any](s S) {
    for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
        s[i], s[j] = s[j], s[i]
    }
}