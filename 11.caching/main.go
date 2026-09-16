package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Cache struct {
	data   map[string]CacheItem
	mutex  sync.RWMutex
	hits   int64
	misses int64
}

type CacheItem struct {
	value     string
	expiresAt time.Time
}

func (c *Cache) Cleanup() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for key, value := range c.data {
		if time.Now().After(value.expiresAt) {
			delete(c.data, key)
		}
	}
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
		atomic.AddInt64(&c.misses, 1)
		c.mutex.RUnlock()
		return "", false

	}

	if time.Now().Before(value.expiresAt) {
		atomic.AddInt64(&c.hits, 1)
		c.mutex.RUnlock()
		return value.value, true
	}

	c.mutex.RUnlock()

	c.mutex.Lock()
	defer c.mutex.Unlock()

	value, ok = c.data[key]

	if !ok {
		atomic.AddInt64(&c.misses, 1)
		return "", false
	}

	if time.Now().Before(value.expiresAt) {
		atomic.AddInt64(&c.hits, 1)
		return value.value, true
	}
	atomic.AddInt64(&c.misses, 1)
	delete(c.data, key)
	return "", false
}

func (c *Cache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.data, key)
}

type CacheStore interface {
	Set(key string, value string, timeout time.Time)
	Get(key string) (string, bool)
	Delete(key string)
	Cleanup()
}

func main() {
	var store CacheStore = &Cache{
		data: map[string]CacheItem{},
	}

	var wg sync.WaitGroup
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				fmt.Println("Done!")
				return
			case <-ticker.C:
				store.Cleanup()
			}
		}
	}()
	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			key := fmt.Sprintf("user:%d", i)
			value := fmt.Sprintf("rahul-%d", i)

			store.Set(key, value, time.Now().Add(5*time.Second))

			value, ok := store.Get(key)

			fmt.Println(key, value, ok)
		}(i)
	}

	wg.Wait()
	close(done)
}
