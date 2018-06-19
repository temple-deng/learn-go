package main

import "fmt"

func main() {
	var arr [3]int = [...]int{1,2,3}
	var brr [3]int = [...]int{4,5,6}
	var m = make(map[[3]int]string)
	m[arr] = "123"
	m[brr] = "456"
	fmt.Println(m)

	var s = "hello world"
	fmt.Println(s[5:10])
}