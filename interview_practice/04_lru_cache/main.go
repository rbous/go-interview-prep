package main

import "fmt"

// LRU Cache
//
// Implement a Least Recently Used cache with a fixed capacity.
// When the cache is full and a new key is inserted, the LEAST recently
// used key should be evicted. Both Get and Put count as "using" a key.
//
// Expected output:
//   get(1) = 1
//   get(2) = 2
//   After adding key 3 (capacity=2, should evict least recent):
//   get(1) = -1
//   get(2) = 2
//   get(3) = 3

type entry struct {
	key, val int
	prev     *entry
	next     *entry
}

type LRUCache struct {
	capacity int
	cache    map[int]*entry
	head     *entry // Most recent
	tail     *entry // Least recent
}

func NewLRUCache(capacity int) *LRUCache {
	head := &entry{}
	tail := &entry{}
	head.next = tail
	tail.prev = head
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[int]*entry),
		head:     head,
		tail:     tail,
	}
}

func (c *LRUCache) remove(e *entry) {
	e.prev.next = e.next
	e.next.prev = e.prev
}

func (c *LRUCache) addToFront(e *entry) {
	e.next = c.head.next
	e.prev = c.head
	c.head.next.prev = e
	c.head.next = e
}

func (c *LRUCache) Get(key int) int {
	if e, ok := c.cache[key]; ok {
		return e.val
	}
	return -1
}

func (c *LRUCache) Put(key, value int) {
	if e, ok := c.cache[key]; ok {
		e.val = value
		c.remove(e)
		c.addToFront(e)
		return
	}

	e := &entry{key: key, val: value}
	c.cache[key] = e
	c.addToFront(e)

	if len(c.cache) > c.capacity {
		// Evict least recently used (tail)
		lru := c.tail.prev
		c.remove(lru)
		delete(c.cache, lru.key)
	}
}

func main() {
	cache := NewLRUCache(2)

	cache.Put(1, 1)
	cache.Put(2, 2)
	fmt.Printf("get(1) = %d\n", cache.Get(1)) // Access key 1 (makes it recently used)
	fmt.Printf("get(2) = %d\n", cache.Get(2)) // Access key 2

	cache.Put(3, 3) // Should evict key 1 (least recently used)

	fmt.Println("After adding key 3 (capacity=2, should evict least recent):")
	fmt.Printf("get(1) = %d\n", cache.Get(1)) // Should be -1 (evicted)
	fmt.Printf("get(2) = %d\n", cache.Get(2)) // Should be 2
	fmt.Printf("get(3) = %d\n", cache.Get(3)) // Should be 3
}
