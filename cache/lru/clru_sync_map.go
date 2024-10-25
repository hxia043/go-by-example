package lru

import (
	"fmt"
	"sync"
)

type LRUCache2 struct {
	head, tail *Node
	caches     sync.Map
	capacity   int
	mu         sync.RWMutex
}

type Node2 struct {
	key, value int
	Pre, Next  *Node
}

func (c *LRUCache2) len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	len := 0
	for cur := c.head; cur != nil; cur = cur.Next {
		len++
	}
	return len
}

func (c *LRUCache2) Print() {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for cur := c.head; cur != nil; cur = cur.Next {
		fmt.Printf("%d:%d ", cur.key, cur.value)
	}
	fmt.Println()
}

func (c *LRUCache2) Get(key int) bool {
	c.mu.RLock()
	node, ok := c.caches.Load(key)
	c.mu.RUnlock()

	if ok {
		c.moveToTail(node.(*Node))
		return true
	}

	return false
}

// moveToTail is to move the existed node to tail
func (c *LRUCache2) moveToTail(node *Node) {
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

func (c *LRUCache2) removeHead() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.caches.Delete(c.head.key)

	if c.head == nil {
		fmt.Println("Failed to move the head node because the cache is empty")
		return
	}

	if c.head == c.tail {
		c.head, c.tail = nil, nil
		return
	}

	n := c.head
	c.head = c.head.Next
	n.Next = nil
	c.head.Pre = nil
}

func (c *LRUCache2) addToTail(node *Node) {
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

func (c *LRUCache2) Put(key int, value int) {
	node, ok := c.caches.Load(key)
	if ok {
		node.(*Node).value = value
		c.moveToTail(node.(*Node))
		return
	}

	newNode := &Node{key: key, value: value}
	c.addToTail(newNode)
	c.caches.Store(key, newNode)
}

func NewCLRUCache2(capacity int) *LRUCache2 {
	return &LRUCache2{
		caches:   sync.Map{},
		capacity: capacity,
	}
}
