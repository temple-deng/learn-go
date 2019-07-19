package main

import (
	"fmt"
	"../tempconv"
)

// type Celsius float64
// type Fahrenheit float64

// const (
// 	AbsoluteZeroC Celsius = -273.15
// 	FreezingC Celsius = 0
// 	BoilingC Celsius = 100
// )

// func CToF (c Celsius) Fahrenheit {
// 	return Fahrenheit(c * 9 / 5 + 32)
// }

// func FToC (f Fahrenheit) Celsius {
// 	return Celsius((f - 32) * 5 / 9)
// }

func main() {
	fmt.Println("绝对零度的华式度数是：", tempconv.CToF(tempconv.AbsoluteZeroC))
	fmt.Println("华式36度的摄氏表示式：", tempconv.FToC(36))
	c := tempconv.FToC(212.0)
	fmt.Println(c.String())
	fmt.Printf("%v\n", c)
	fmt.Printf("%s\n", c)
	fmt.Println(c)
	fmt.Printf("%g\n", c)
	fmt.Println(float64(c))
}

// func (c Celsius) String() string {
// 	return fmt.Sprintf("%g°C", c)
// }