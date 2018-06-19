package main

import "fmt"

func main() {
	var s = []int{1,2,3}
	fmt.Printf("%v\t%[1]T\n",s)

	var s1 []int
	fmt.Println(s1 == nil)

	var runes []rune
	for _, r := range "Hello, world" {
		runes = append(runes, r)
	}
	
	fmt.Printf("%q\n", runes)
}