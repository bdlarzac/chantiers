package tiglib

import (
	"time"
)

// IsBefore() is the same as time.IsBefore(), but takes IsZero() into account.
func IsBefore(t1, t2 time.Time) bool {
	if t1.IsZero() {
		return false
	}
	if t2.IsZero() {
		return true
	}
	return t1.Before(t2)
}

