package cache

import(
	"time"

	memoryCache "github.com/patrickmn/go-cache"
)

type cache struct {
	memory *memoryCache.Cache
}

type MemoryCache interface {
	Set(key string, value any, d time.Duration)
	Get(key string) (any, bool)
	Remove(key string)
}

func NewMemoryCache() MemoryCache {
	c := memoryCache.New(memoryCache.NoExpiration, time.Minute)

	return &cache{memory: c}
}

func (c *cache) Set(key string, value any, ttl time.Duration) {
	c.memory.Set(key, value, ttl)
}

func (c *cache) Get(key string) (any, bool) {
	return c.memory.Get(key)
}

func (c *cache) Remove(key string) {
	c.memory.Delete(key)
}