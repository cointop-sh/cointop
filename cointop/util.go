package cointop

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/cointop-sh/cointop/pkg/open"
	log "github.com/sirupsen/logrus"
)

// OpenLink opens the url in a browser
func (ct *Cointop) OpenLink() error {
	log.Debug("OpenLink()")
	open.URL(ct.RowLink())
	return nil
}

// GetBytes returns the interface in bytes form
func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Slugify returns a slugified string
func Slugify(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	return s
}

// TruncateString returns a truncated string
func TruncateString(value string, maxLen int) string {
	dots := "..."
	if len(value) > maxLen {
		value = fmt.Sprintf("%s%s", value[0:maxLen-3], dots)
	}
	return value
}

// ClearSyncMap clears a sync.Map
func (ct *Cointop) ClearSyncMap(syncMap *sync.Map) {
	syncMap.Range(func(key interface{}, value interface{}) bool {
		syncMap.Delete(key)
		return true
	})
}

// NormalizeFloatString normalizes a float as a string
func normalizeFloatString(input string, allowNegative bool) string {
	re := regexp.MustCompile(`(\d+\.\d+|\.\d+|\d+)`)
	if allowNegative {
		re = regexp.MustCompile(`-?(\d+\.\d+|\.\d+|\d+)`)
	}
	result := re.FindStringSubmatch(input)
	if len(result) > 0 {
		return result[0]
	}

	return ""
}
