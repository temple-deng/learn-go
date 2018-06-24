# 第八章 Goroutines 和 Channels

Go 语言中的并发程序可以用两种手段来实现。本章讲解 goroutine 和 channel，其支持
“顺序通信进程”(communicating sequential processes)或被简称为 CSP。CSP 是一种现代
的并发编程模型，在这种编程模型中值会在不同的运行实例 goroutine 中传递。第九章覆盖更为
传统的并发模型：多线程共享内存。     

## 8.1 Goroutines

在 Go 语言中，每一个并发的执行单元叫做一个 goroutine。目前为止，我们可以简单地把 goroutine
类比为一个线程。     

当一个程序启动时，其主函数即在一个单独的 goroutine 中运行，我们叫它 main goroutine。新的
goroutine 会用 `go` 语句来创建。在语法上，`go` 语句是在一个普通的函数或方法调用前加上
关键字 `go`。`go` 语句会使其语句中的函数在一个新创建的 goroutine 中运行。而 `go` 语句
本身会迅速地完成。     

```go
f()      // call f(); wait for it to return
go f()   // create a new goroutine that calls f(); don't wait
```     

```go
package main

import (
	"time"
	"fmt"
)

func main() {
	go spinner(100 * time.Millisecond)
	const n = 45
	fibN := fib(n)
	fmt.Printf("\rFibonacci(%d) = %d\n", n, fibN)
}

func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

func fib(x int) int {
	if x < 2 {
		return x
	}
	return fib(x-1) + fib(x-2)
}
```    

动画显示了几秒之后，`fib(45)` 的调用成功地返回，并且打印结果。然后主函数返回。主函数返回
时，所有的goroutine 都会被直接打断，程序退出。除了从主函数退出或者直接终止程序之外，没有
其它的编程方法能够让一个 goroutine 来打断另一个的执行，但是之后可以看到一种方式来实现
这个目的，通过 goroutine 之间的通信来让一个 goroutine 请求其它的 goroutine，让其它
的 goroutine 自行结束执行。     

## 8.4 Channels

channel 是一种通信机制，它可以让一个 goroutine 通过它给另一个 goroutine 发送信息。每个
channel 都有一个特殊的类型，也就是channel 可发送数据的类型。一个可以发送 int 类型数据的
channel 一般写为 chan int。      

使用内置的 `make` 函数，可以创建一个 channel:    

`ch := make(chan int)`     

一个 channel 有发送和接受两个主要操作，都是通信行为。一个发送语句将一个值从一个 goroutine
通过 channel 发送到另一个执行接收操作的 goroutine。发送和接收两个操作都是用 `<-` 运算符。
在发送语句中， `<-` 运算符分割 channel 和要发送的值。在接收语句中，`<-` 运算符写在
channel 对象之前。一个不使用接收结果的接收操作也是合法的。    

```go
ch <- x   // a send statement
x = <- ch  // a receive expression in an assignment statement
<- ch
```    

channel 还支持 `close` 操作，用于关闭 channel，关闭后的 channel 不能再进行发送操作，
否则会导致 panic。不过已经close 的 channel 还可以接收到之前已经成功发送的数据。如果 
channel 中已经没有数据的话，产生一个零值的数据。     

`close(ch)`      

以最简单方式调用 `make` 创建的是一个无缓冲的 channel，但是我们也可以指定第二个整型参数，
对应 channel 的容量。如果 channel 的容量大于零，那么该 channel 就是带缓冲的channel.    

### 8.4.1 无缓存的 channel

一个无缓存的channel 如果进行发送操作，将会导致发送者自身 goroutine 阻塞，直到另一个
goroutine 在相同的 channel 上执行接收操作，当发送的值通过 channel 成功传输之后，两个
goroutine 可以继续执行后面的语句。反之，如果接收操作先发生，那么接收者goroutine 也将
阻塞，直到另一个 goroutine 在相同的channel 上执行发送操作。     

基于无缓存的channel 的发送和接收操作将导致两个 goroutine 做一次同步操作。因为这个原因，
无缓存 channel 也被称为同步 channel。当通过一个无缓存 channel 发送数据时，接收者收到
数据发生在唤醒发送者 goroutine 之前。      

