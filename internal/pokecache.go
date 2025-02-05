package pokecache

import(
	"time"
	"sync"
)

type Cache struct  {
	cachemap map[string] cacheEntry
	mu sync.RWMutex
	interval time.Duration
}

type cacheEntry struct  {
	createdAt time.Time
	val []byte
}

func NewCache(new_interval time.Duration) *Cache {

	cache := Cache{
		cachemap: make(map[string] cacheEntry),
		mu: sync.RWMutex{},
		interval: new_interval, 
	}
	go cache.reapLoop()
	return &cache

}

func (c *Cache) Add(newKey string, newVal []byte){
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cachemap[newKey] = cacheEntry{
		createdAt: time.Now(),
		val: newVal,
	}
}

func (c *Cache) Get(key string) ([]byte, bool){
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.cachemap[key]
	if !ok {
		return []byte{}, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for {
		select {
		case <- ticker.C:
			c.mu.Lock()
			for key, entry := range c.cachemap{
				if time.Since(entry.createdAt)>c.interval{
					delete(c.cachemap,key)
				}
			}
			c.mu.Unlock()
		}
	}
}