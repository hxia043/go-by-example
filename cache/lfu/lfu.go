package lfu

import "container/list"

type LFUCache struct {
	lists    map[int]*list.List
	nodes    map[int]*list.Element
	min      int
	capacity int
}

type node struct {
	key, value int
	frenquency int
}

func (c *LFUCache) updateNode(n *node) {
	c.lists[n.frenquency].Remove(c.nodes[n.key])

	if n.frenquency == c.min && c.lists[c.min].Len() == 0 {
		c.min++
	}

	n.frenquency++
	if _, ok := c.lists[n.frenquency]; !ok {
		c.lists[n.frenquency] = list.New()
	}

	ne := c.lists[n.frenquency].PushFront(n)
	c.nodes[n.key] = ne
}

func (c *LFUCache) Get(key int) bool {
	e, ok := c.nodes[key]
	if ok {
		n := e.Value.(*node)
		c.updateNode(n)
		return true
	}

	return false
}

func (c *LFUCache) Put(key int, value int) {
	e, ok := c.nodes[key]
	if ok {
		n := e.Value.(*node)
		c.updateNode(n)
		return
	}

	if len(c.nodes) == c.capacity {
		e := c.lists[c.min].Back()
		c.lists[c.min].Remove(e)
		delete(c.nodes, e.Value.(*node).key)
	}

	c.min = 1
	if _, ok = c.lists[1]; !ok {
		c.lists[1] = list.New()
	}

	n := &node{key: key, value: value, frenquency: 1}
	e = c.lists[1].PushFront(n)
	c.nodes[key] = e
}

func New(capacity int) *LFUCache {
	return &LFUCache{
		lists:    make(map[int]*list.List),
		nodes:    make(map[int]*list.Element),
		capacity: capacity,
	}
}
