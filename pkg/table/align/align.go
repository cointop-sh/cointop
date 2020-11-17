package align

import (
	"fmt"
	"strings"
)

// AlignLeft align left
func AlignLeft(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}

	return fmt.Sprintf("%s%s", s, strings.Repeat(" ", n-len(s)))
}

// AlignRight align right
func AlignRight(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}

	return fmt.Sprintf("%s%s", strings.Repeat(" ", n-len(s)), s)
}

// AlignCenter align center
func AlignCenter(s string, n int) string {
	if len(s) > n {
		return s[:n]
	}

	pad := (n - len(s)) / 2
	lpad := pad
	rpad := n - len(s) - lpad

	return fmt.Sprintf("%s%s%s", strings.Repeat(" ", lpad), s, strings.Repeat(" ", rpad))
}
