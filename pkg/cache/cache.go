package cache

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
)

// Cache is cache struct
type Cache struct {
	cache *gocache.Cache
}

// NoExpiration is constant
const NoExpiration = gocache.NoExpiration

// New returns new cache instance
func New(defaultExpiration, cleanupInterval time.Duration) *Cache {
	return &Cache{
		cache: gocache.New(defaultExpiration, cleanupInterval),
	}
}

// Set sets cache item
func (c *Cache) Set(k string, x interface{}, d time.Duration) {
	c.cache.Set(k, x, d)
}

// Get gets cache item
func (c *Cache) Get(k string) (interface{}, bool) {
	return c.cache.Get(k)
}

// Delete deletes cache item
func (c *Cache) Delete(k string) {
	c.cache.Delete(k)
}
