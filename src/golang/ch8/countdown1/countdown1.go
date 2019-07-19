package main

import (
	"os"
	"fmt"
	"time"
)

func main() {
	fmt.Println("Commencing countdown")
	abort := make(chan struct{})
	go func() {
		os.Stdin.Read(make([]byte, 1))
		abort <- struct{}{}
	}()

	for countdown := 10; countdown > 0; countdown-- {
		select {
		case <- time.After(10 * time.Second):
			fmt.Println("Launch")
			return
		case <- abort:
			fmt.Println("Launch aborted")
			return
		}
	}
}