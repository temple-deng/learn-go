package main

import (
	"fmt"
	"os"
)

func main() {
	var s, sep string
	for i, arg := range os.Args[1:] {
		s += sep + arg
		fmt.Println(i)
		sep = " "
	}

	fmt.Println(s)
}