package pokecache

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	mu       sync.RWMutex
	data     map[string]cacheEntry
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) Cache {
	return Cache{sync.RWMutex{}, make(map[string]cacheEntry, 0), interval}
}

func (C *Cache) Add(key string, val []byte) {
	C.reapLoop()
	fmt.Println(fmt.Sprintf("Adding key '%v' to cache", key))
	C.mu.Lock()
	if _, ok := C.data[key]; !ok {
		C.data[key] = cacheEntry{time.Now(), val}
	}
	C.mu.Unlock()
}

func (C *Cache) Get(key string) ([]byte, bool) {
	C.reapLoop()
	if _, ok := C.data[key]; ok {
		fmt.Println(fmt.Sprintf("Hit cache key'%v'", key))
		return C.data[key].val, true
	}
	return nil, false
}

func (C *Cache) reapLoop() {
	C.mu.Lock()
	for k, v := range C.data {
		if v.createdAt.Before(time.Now().Add(-C.interval)) {
			fmt.Println(fmt.Sprintf("Cleanup of cache key'%v'", k))
			delete(C.data, k)
		}
	}
	C.mu.Unlock()
}

func cacheManager(ticker *time.Ticker, C *Cache) {
	for true {
		<-ticker.C
		C.reapLoop()
	}
}
