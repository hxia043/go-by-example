package lru

import (
	"fmt"
)

type LRUCache struct {
	head, tail *Node
	caches     map[int]*Node
	capacity   int
}

func (c *LRUCache) len() int {
	len := 0
	for cur := c.head; cur != nil; cur = cur.Next {
		len++
	}
	return len
}

func (c *LRUCache) Print() {
	for cur := c.head; cur != nil; cur = cur.Next {
		fmt.Printf("%d:%d ", cur.key, cur.value)
	}
	fmt.Println()
}

func (c *LRUCache) Get(key int) bool {
	node, ok := c.caches[key]
	if ok {
		c.moveToTail(node)
		return true
	}

	return false
}

// modeToTail is to move the existed node to tail
func (c *LRUCache) moveToTail(node *Node) {
	if c.head == nil {
		fmt.Println("Failed to move the node to tail because the cache is empty")
		return
	}

	if c.tail == node {
		return
	}

	if c.head == node {
		c.head = node.Next
		c.head.Pre = nil
		node.Next = nil
	} else {
		node.Pre.Next = node.Next
		node.Next.Pre = node.Pre
		node.Pre, node.Next = nil, nil
	}

	node.Pre = c.tail
	c.tail.Next = node
	c.tail = node
}

func (c *LRUCache) removeHead() {
	delete(c.caches, c.head.key)

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

func (c *LRUCache) addToTail(node *Node) {
	if c.head == nil {
		c.head, c.tail = node, node
		return
	}

	node.Pre = c.tail
	c.tail.Next = node
	node.Next = nil
	c.tail = node

	if c.len() > c.capacity {
		c.removeHead()
	}
}

func (c *LRUCache) Put(key int, value int) {
	node, ok := c.caches[key]
	if ok {
		node.value = value
		c.moveToTail(node)
		return
	}

	newNode := &Node{key: key, value: value}
	c.addToTail(newNode)
	c.caches[key] = newNode
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		caches:   map[int]*Node{},
		capacity: capacity,
	}
}
