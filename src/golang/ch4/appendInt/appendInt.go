package main

import "fmt"

func main() {
	var a = []int{1,2,3}
	var b = 5

	fmt.Println(appendInt(a, b))	
}

func appendInt(x []int, y int) []int {
	var z []int
	zlen := len(x) + 1
	if zlen <= cap(x) {
		z = x[:zlen]
	} else {
		zcap := zlen

		if zcap < zlen * 2 {
			zcap = zlen * 2
		}

		z = make([]int, zlen, zcap)
		copy(z, x)
	}
	z[len(x)] = y
	return z
}