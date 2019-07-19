package main

import "fmt"

type Flags uint

const (
	FlagUp Flags = 1 << iota
	FlagBroadcast
	FlagLoopback
	FlagPointToPoint
	FlagMulticast
)

func main() {
	const a, b = 3, 5

	const (
		c = 1
		d
		e = 2
		f
	)

	fmt.Println(a+b)
	fmt.Println(c,d,e,f)

	type Weekday int

	const (
		Sunday Weekday = iota
		Monday
		Tuesday
		Wednesday
		Thursday
		Friday
		Saturday
	)

	fmt.Printf("Type: %T\t%2[1]d\n", Sunday)
	fmt.Printf("Type: %T\t%2[1]d\n", Monday)
	fmt.Printf("Type: %T\t%2[1]d\n", Tuesday)
	fmt.Printf("Type: %T\t%2[1]d\n", Wednesday)
	fmt.Printf("Type: %T\t%2[1]d\n", Thursday)
	fmt.Printf("Type: %T\t%2[1]d\n", Friday)
	fmt.Printf("Type: %T\t%2[1]d\n", Saturday)
}