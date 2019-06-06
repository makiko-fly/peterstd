package cache

import (
	"container/list"
	"sync"
)

var _ Cache = (*lruCache)(nil)

type lruEntry struct {
	key   interface{}
	value interface{}
	hits  int64
}

type LRUCache = *lruCache

type lruCache struct {
	cap     int
	fast    map[interface{}]*list.Element
	ll      *list.List
	lock    sync.RWMutex
	onEvict func(k, v interface{})
}

func NewLRUCache() *lruCache {
	return &lruCache{
		cap:     1118,
		fast:    make(map[interface{}]*list.Element),
		ll:      list.New(),
		onEvict: nil,
	}
}

func (c *lruCache) Cap(cap int) *lruCache {
	c.cap = cap
	return c
}

func (c *lruCache) OnEvict(fn func(k, v interface{})) *lruCache {
	c.onEvict = fn
	return c
}

func (c *lruCache) evict() {
	// remove oldest items
	for c.ll.Len() > c.cap {
		elem := c.ll.Back()
		if c.onEvict != nil {
			en := elem.Value.(*lruEntry)
			c.onEvict(en.key, en.value)
		}
		c.removeElement(elem)
	}
}

func (c *lruCache) remove(key interface{}) {
	elem, ok := c.fast[key]
	if ok {
		c.ll.Remove(elem)
		delete(c.fast, key)
	}
}

func (c *lruCache) removeElement(elem *list.Element) {
	e := elem.Value.(*lruEntry)
	delete(c.fast, e.key)
	c.ll.Remove(elem)
}

func (c *lruCache) Set(key, value interface{}) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	defer c.evict()

	elem, ok := c.fast[key]
	if ok {
		elem.Value.(*lruEntry).value = value
		c.ll.MoveToFront(elem)
		return nil
	}
	c.fast[key] = c.ll.PushFront(&lruEntry{
		key:   key,
		value: value,
	})
	return nil
}

func (c *lruCache) Get(k interface{}) (interface{}, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	elem, ok := c.fast[k]
	if !ok {
		return nil, ErrorMissingKey
	}
	c.ll.MoveToFront(elem)
	e := elem.Value.(*lruEntry)
	return e.value, nil
}

func (c *lruCache) Remove(k interface{}) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.remove(k)
	return nil
}

func (c *lruCache) Len() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.ll.Len()
}
