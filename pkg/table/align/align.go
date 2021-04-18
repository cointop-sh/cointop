package align

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// AlignLeft align left
func AlignLeft(s string, n int) string {
	slen := utf8.RuneCountInString(s)
	if slen > n {
		return s[:n]
	}

	return fmt.Sprintf("%s%s", s, strings.Repeat(" ", n-slen))
}

// AlignRight align right
func AlignRight(s string, n int) string {
	slen := utf8.RuneCountInString(s)
	if slen > n {
		return s[:n]
	}

	return fmt.Sprintf("%s%s", strings.Repeat(" ", n-slen), s)
}

// AlignCenter align center
func AlignCenter(s string, n int) string {
	slen := utf8.RuneCountInString(s)
	if slen > n {
		return s[:n]
	}

	pad := (n - slen) / 2
	lpad := pad
	rpad := n - slen - lpad

	return fmt.Sprintf("%s%s%s", strings.Repeat(" ", lpad), s, strings.Repeat(" ", rpad))
}
