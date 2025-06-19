package cache

import (
	"container/list"
	"encoding/json"
	"sync"
	"time"
)

type LRUCache struct {
	capacity        int
	cache           map[string]*list.Element
	list            *list.List
	persistenceSecs int
	mu              *sync.RWMutex
}

type pair struct {
	key   string
	value json.RawMessage
	start time.Time
}

func NewLRUCache(capacity, persistenceSecs int) *LRUCache {
	return &LRUCache{
		capacity:        capacity,
		cache:           make(map[string]*list.Element),
		list:            list.New(),
		persistenceSecs: persistenceSecs,
		mu:              &sync.RWMutex{},
	}
}

func (c *LRUCache) Get(key string) json.RawMessage {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.cache[key]; ok {
		curPair := elem.Value.(*pair)

		if time.Now().Unix()-curPair.start.Unix() > int64(c.persistenceSecs) {
			c.list.Remove(elem)
			delete(c.cache, key)

			return nil
		}

		c.list.MoveToFront(elem)
		return curPair.value

	}
	return nil
}

func (c *LRUCache) Set(key string, value json.RawMessage) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.cache[key]; ok {
		c.list.MoveToFront(elem)
		elem.Value.(*pair).value = value
		return
	}

	if c.cache == nil {
		c.cache = make(map[string]*list.Element)
	}

	c.cache[key] = c.list.PushFront(&pair{key: key, value: value, start: time.Now()})

	if c.list.Len() > c.capacity {
		elem := c.list.Back()
		if elem != nil {
			c.list.Remove(elem)
			delete(c.cache, elem.Value.(*pair).key)
		}
	}

	return
}

func (c *LRUCache) Del(key string) {
	elem, ok := c.cache[key]
	if !ok {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.list.Remove(elem)
	delete(c.cache, elem.Value.(*pair).key)

	return
}
