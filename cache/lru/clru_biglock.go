package lru

import (
	"fmt"
	"sync"
)

type LRUCache1 struct {
	capacity int
	cache    map[int]*Node
	head     *Node
	tail     *Node
	mu       sync.Mutex
}

func NewCLRUCache1(capacity int) *LRUCache1 {
	return &LRUCache1{
		capacity: capacity,
		cache:    make(map[int]*Node),
	}
}

func (c *LRUCache1) Get(key int) int {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, ok := c.cache[key]; ok {
		c.moveToHead(node)
		return node.value
	}
	return -1
}

func (c *LRUCache1) Put(key int, value int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, ok := c.cache[key]; ok {
		node.value = value
		c.moveToHead(node)
	} else {
		newNode := &Node{key: key, value: value}
		c.cache[key] = newNode
		c.addToHead(newNode)

		if len(c.cache) > c.capacity {
			removedNode := c.removeTail()
			delete(c.cache, removedNode.key)
		}
	}
}

func (c *LRUCache1) moveToHead(node *Node) {
	c.removeNode(node)
	c.addToHead(node)
}

func (c *LRUCache1) addToHead(node *Node) {
	node.Pre = nil
	node.Next = c.head

	if c.head != nil {
		c.head.Pre = node
	}
	c.head = node

	if c.tail == nil {
		c.tail = node
	}
}

func (c *LRUCache1) removeNode(node *Node) {
	if node.Pre != nil {
		node.Pre.Next = node.Next
	} else {
		c.head = node.Next
	}

	if node.Next != nil {
		node.Next.Pre = node.Pre
	} else {
		c.tail = node.Pre
	}
}

func (c *LRUCache1) removeTail() *Node {
	node := c.tail
	c.removeNode(node)
	return node
}

func (c *LRUCache1) Print() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for node := c.head; node != nil; node = node.Next {
		fmt.Printf("%d:%d ", node.key, node.value)
	}
	fmt.Println()
}
