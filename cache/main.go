package main

import (
	"cache/lru"
	"fmt"
)

func main() {
	c := lru.NewLRUCache(6)
	c.Put(1, 1)
	fmt.Println(c.Get(1))
	c.Print()

	c.Put(2, 2)
	c.Put(3, 3)
	c.Put(4, 4)
	c.Print()

	c.Put(5, 5)
	c.Print()

	c.Put(6, 6)
	c.Print()

	c.Put(7, 7)
	c.Print()

	c.Get(2)
	c.Print()

	c.Put(1, 1)
	c.Print()

	fmt.Println(c.Get(3))
}
