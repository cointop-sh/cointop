package cointop

import (
	"bytes"
	"encoding/gob"
	"strings"

	"github.com/miguelmota/cointop/cointop/common/open"
)

// OpenLink opens the url in a browser
func (ct *Cointop) OpenLink() error {
	ct.debuglog("openLink()")
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
