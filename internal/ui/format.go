package ui

import (
	"strings"
)

func MarkupAsProgress(s string, w int, pct float64) string {
	if w == 0 {
		return s
	}

	if pct < 0 {
		pct = 0
	} else if pct > 1 {
		pct = 1
	}

	pad := w - len(s)
	if pad > 0 {
		s = s + strings.Repeat(" ", pad)
	}
	pre := int(float64(w) * pct)
	if pre < 0 {
		return s
	} else if pre >= len(s) {
		return s
	}
	return "%S" + EscapePercent(s[:pre]) + "%N" + EscapePercent(s[pre:])
}

func EscapePercent(s string) string {
	return strings.Replace(s, "%", "%%", -1)
}
