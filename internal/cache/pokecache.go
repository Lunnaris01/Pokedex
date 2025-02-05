package pokecache

import(
	"time"
	"sync"
)

type Cache struct  {
	Cachemap map[string] cacheEntry
	Mu sync.RWMutex
	Interval time.Duration
}

type cacheEntry struct  {
	CreatedAt time.Time
	Val []byte
}

func NewCache(new_interval time.Duration) *Cache {

	cache := Cache{
		Cachemap: make(map[string] cacheEntry),
		Mu: sync.RWMutex{},
		Interval: new_interval, 
	}
	go cache.reapLoop()
	return &cache

}

func (c *Cache) Add(newKey string, newVal []byte){
	c.Mu.Lock()
	defer c.Mu.Unlock()
	c.Cachemap[newKey] = cacheEntry{
		CreatedAt: time.Now(),
		Val: newVal,
	}
}

func (c *Cache) Get(key string) ([]byte, bool){
	c.Mu.RLock()
	defer c.Mu.RUnlock()
	entry, ok := c.Cachemap[key]
	if !ok {
		return []byte{}, false
	}
	return entry.Val, true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.Interval)
	defer ticker.Stop()
	for {
		select {
		case <- ticker.C:
			c.Mu.Lock()
			for key, entry := range c.Cachemap{
				if time.Since(entry.CreatedAt)>c.Interval{
					delete(c.Cachemap,key)
				}
			}
			c.Mu.Unlock()
		}
	}
}