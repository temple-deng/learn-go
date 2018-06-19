package main

import "fmt"

func main() {
	var s []string
	var m map[string]int

	// fmt.Println(s[0])
	fmt.Println(m["0"])

	for i, v := range s {
		fmt.Println(i, v)
	}

	for i, v := range m {
		fmt.Println(i, v)
	}
}