package main

import (
	"fmt"
	"time"
)

func main() {
	var x = make(chan int, 10)
	var c = 0

	for i := 0; i < 10; i++ {
		go func(i int) {
			x <- i
		}(i)
	}

	for i := 0; i < 10; i++ {
		go func(c *int) {
			for n := range x {
				*c = *c + n
				fmt.Printf("%d ", *c)
			}
		}(&c)
	}

	time.Sleep(1 * time.Second)
	fmt.Println()
}
