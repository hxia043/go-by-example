package lru

import (
	"fmt"
	"hash/fnv"
	"strconv"
	"sync"
)

type CLRUCache struct {
	mu         sync.RWMutex
	capacity   int
	buckets    []*bucket
	bucketMask uint32
	head, tail *Node
}

type bucket struct {
	mu   sync.RWMutex
	keys map[int]*Node
}

func (b *bucket) get(key int) (*Node, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	node, ok := b.keys[key]
	if ok {
		return node, true
	}
	return nil, false
}

func (b *bucket) put(key int, node *Node) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.keys[key] = node
}

func (b *bucket) delete(key int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	delete(b.keys, key)
}

func hashIdx(key int, mask uint32) uint32 {
	h := fnv.New32a()
	h.Write([]byte(strconv.Itoa(key)))
	return h.Sum32() & mask
}

func (c *CLRUCache) moveToTail(node *Node) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.head == nil {
		fmt.Println("Failed to move the node to tail because the cache is empty")
		return
	}

	if c.tail == node {
		return
	}

	if c.head == node {
		c.head = node.Next
		if c.head != nil {
			c.head.Pre = nil
		}
		node.Next = nil
	} else {
		if node.Pre != nil {
			node.Pre.Next = node.Next
		}
		if node.Next != nil {
			node.Next.Pre = node.Pre
		}
		node.Pre, node.Next = nil, nil
	}

	node.Pre = c.tail
	if c.tail != nil {
		c.tail.Next = node
	}

	c.tail = node
}

func (c *CLRUCache) Get(key int) bool {
	i := hashIdx(key, c.bucketMask)
	node, ok := c.buckets[i].get(key)
	if ok {
		c.moveToTail(node)
		return true
	}

	return false
}

func (c *CLRUCache) len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	len := 0
	for cur := c.head; cur != nil; cur = cur.Next {
		len++
	}
	return len
}

func (c *CLRUCache) removeHead() {
	i := hashIdx(c.head.key, c.bucketMask)
	c.buckets[i].delete(c.head.key)

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.head == nil {
		fmt.Println("Failed to move the head node because the cache is empty")
		return
	}

	if c.head == c.tail {
		c.head, c.tail = nil, nil
		return
	}

	c.head = c.head.Next
	c.head.Pre.Next = nil
	c.head.Pre = nil
}

func (c *CLRUCache) addToTail(node *Node) {
	c.mu.Lock()

	if c.head == nil {
		c.head, c.tail = node, node
		c.mu.Unlock()
		return
	}

	node.Pre = c.tail
	c.tail.Next = node
	node.Next = nil
	c.tail = node

	c.mu.Unlock()

	if c.len() > c.capacity {
		c.removeHead()
	}
}

func (c *CLRUCache) Put(key int, value int) {
	i := hashIdx(key, c.bucketMask)
	node, ok := c.buckets[i].get(key)
	if ok {
		node.value = value
		c.moveToTail(node)
		return
	}

	newNode := &Node{key: key, value: value}
	c.addToTail(newNode)
	c.buckets[i].put(key, newNode)
}

func NewCLRUCache(capacity int) *CLRUCache {
	c := &CLRUCache{
		capacity:   capacity,
		buckets:    make([]*bucket, 1024),
		bucketMask: 1023,
	}

	for i := range c.buckets {
		c.buckets[i] = &bucket{keys: map[int]*Node{}}
	}

	return c
}
