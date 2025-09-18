package gocache_ttl

import (
	"container/list"
	"sync"
	"time"
)

type item struct {
	value      any
	expiration int64
	elem       *list.Element
}

type Cache struct {
	mu       sync.RWMutex
	items    map[string]item
	order    *list.List
	maxSize  int
	interval time.Duration
	stop     chan struct{}
}

func NewCache(interval time.Duration, maxSize int) *Cache {
	c := &Cache{
		items:    make(map[string]item),
		order:    list.New(),
		maxSize:  maxSize,
		interval: interval,
		stop:     make(chan struct{}),
	}

	if interval > 0 {
		go c.cleanup()
	}

	return c
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var exp int64

	if ttl > 0 {
		exp = time.Now().Add(ttl).UnixNano()
	}

	if it, ok := c.items[key]; ok {
		it.value = value
		it.expiration = exp
		c.items[key] = it
		return
	}

	if c.maxSize > 0 && len(c.items) >= c.maxSize {
		front := c.order.Front()
		if front != nil {
			oldestKey := front.Value.(string)
			c.order.Remove(front)
			delete(c.items, oldestKey)
		}
	}

	elem := c.order.PushBack(key)
	c.items[key] = item{
		value:      value,
		expiration: exp,
		elem:       elem,
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	it, ok := c.items[key]
	if !ok {
		return nil, false
	}
	if it.expiration > 0 && time.Now().UnixNano() > it.expiration {
		return nil, false
	}
	return it.value, true
}

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if it, ok := c.items[key]; ok {
		c.order.Remove(it.elem)
		delete(c.items, key)
	}
}

func (c *Cache) cleanup() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-c.stop:
			return
		case <-ticker.C:
			now := time.Now().UnixNano()
			c.mu.Lock()
			for k, it := range c.items {
				if it.expiration > 0 && now > it.expiration {
					c.order.Remove(it.elem)
					delete(c.items, k)
				}
			}
			c.mu.Unlock()
		}
	}
}

func (c *Cache) Exists(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	it, ok := c.items[key]
	if !ok {
		return false
	}

	if it.expiration > 0 && time.Now().UnixNano() > it.expiration {
		return false
	}
	return true
}

func (c *Cache) Keys() []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	var keys []string
	now := time.Now().UnixNano()

	for k, it := range c.items {
		if it.expiration > 0 && now > it.expiration {
			c.order.Remove(it.elem)
			delete(c.items, k)
			continue
		}
		keys = append(keys, k)
	}

	return keys
}

func (c *Cache) Close() {
	close(c.stop)
}
