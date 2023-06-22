package align

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/acarl005/stripansi"
)

// Left align left
func Left(t string, n int) string {
	s := stripansi.Strip(t)
	slen := utf8.RuneCountInString(s)
	if slen > n {
		return s[:n]
	}

	return fmt.Sprintf("%s%s", t, strings.Repeat(" ", n-slen))
}

// Right align right
func Right(t string, n int) string {
	s := stripansi.Strip(t)
	slen := utf8.RuneCountInString(s)
	if slen > n {
		return s[:n]
	}

	return fmt.Sprintf("%s%s", strings.Repeat(" ", n-slen), t)
}

// Center align center
func Center(t string, n int) string {
	s := stripansi.Strip(t)
	slen := utf8.RuneCountInString(s)
	if slen > n {
		return s[:n]
	}

	pad := (n - slen) / 2
	lpad := pad
	rpad := n - slen - lpad

	return fmt.Sprintf("%s%s%s", strings.Repeat(" ", lpad), t, strings.Repeat(" ", rpad))
}
