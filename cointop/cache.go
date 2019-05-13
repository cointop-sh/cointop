package cointop

import (
	"fmt"
	"strings"
)

func (ct *Cointop) cacheKey(key string) string {
	return strings.ToLower(fmt.Sprintf("%s_%s", ct.apiChoice, key))
}
