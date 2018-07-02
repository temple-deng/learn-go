package main

import (
	"fmt"
)

func main() {
	naturals := make(chan int)
	squares := make(chan int)

	go func() {
		for x := 1; x < 10; x++ {
			naturals <- x
		}
		close(naturals)
	}()

	go func() {
		for {
			x, ok := <- naturals
			if !ok {
				break
			}
			squares <- x * x
		}
		close(squares)
	}()

	for {
		x, ok := <- squares
		if !ok {
			return
		}
		fmt.Println(x)
	}
}

