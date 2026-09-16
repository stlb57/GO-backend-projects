package main

import (
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	data  map[string]CacheItem
	mutex sync.RWMutex
}

type CacheItem struct {
	value     string
	expiresAt time.Time
}

func (c *Cache) Set(key string, value string, timeout time.Time) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[key] = CacheItem{
		value:     value,
		expiresAt: timeout,
	}
}

func (c *Cache) Get(key string) (string, bool) {
	c.mutex.RLock()

	value, ok := c.data[key]

	if !ok {
		c.mutex.RUnlock()
		return "", false
	}

	if time.Now().Before(value.expiresAt) {
		c.mutex.RUnlock()
		return value.value, true
	}
	value, ok = c.data[key]
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if !ok {
		return "", false
	}

	if time.Now().Before(value.expiresAt) {
		return value.value, true
	}

	delete(c.data, key)
	return "", false

}

func (c *Cache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.data, key)
}

func main() {
	var cache_inst Cache = Cache{
		data: map[string]CacheItem{},
	}

	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			key := fmt.Sprintf("user:%d", i)
			value := fmt.Sprintf("rahul-%d", i)

			cache_inst.Set(key, value, time.Now().Add(5*time.Second))

			value, ok := cache_inst.Get(key)

			fmt.Println(key, value, ok)
		}(i)
	}

	wg.Wait()
}
