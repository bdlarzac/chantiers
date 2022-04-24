package tiglib

import (
	"fmt"
	"time"
)

// Renvoie une date YYYY-MM-DD
func DateIso(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())
}

// Renvoie une date JJ/MM/AAAA
func DateFr(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return fmt.Sprintf("%02d/%02d/%d", t.Day(), t.Month(), t.Year())
}

// Renvoie une date du style "12 septembre 2019"
func DateFrText(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	month := ""
	switch t.Month() {
	case 1:
		month = "janvier"
	case 2:
		month = "février"
	case 3:
		month = "mars"
	case 4:
		month = "avril"
	case 5:
		month = "mai"
	case 6:
		month = "juin"
	case 7:
		month = "juillet"
	case 8:
		month = "août"
	case 9:
		month = "septembre"
	case 10:
		month = "octobre"
	case 11:
		month = "novembre"
	case 12:
		month = "décembre"
	}
	return fmt.Sprintf("%d %s %d", t.Day(), month, t.Year())
}
