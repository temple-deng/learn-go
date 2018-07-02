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

### 8.4.2 串联的 Channels

Channels 也可以用于将多个 goroutine 链接在一起，一个 Channel 的输出作为下一个Channel
的输入。这种串联的 channels 就是所谓的管道(pipeline)。     

```go
func main() {
	naturals := make(chan int)
	squares := make(chan int)

	go func() {
		for x := 0;; x++ {
			naturals <- x
		}
	}()

	go func() {
		for {
			x := <- naturals
			squares <- x * x
		}
	}()

	for {
		fmt.Println(<- squares)
	}
}
```      

如果发送者知道，没有更多的值需要发送到 channel 的话，那么让接收者也能及时知道没有多余的值
可接收将是有用的，因为接收者可以停止不必要的接收等待。这可以通过内置的 close 函数来关闭
channel 实现：   

`close(naturals)`     

当一个 channel 被关闭后，再向该 channel 发送数据将导致 panic 异常。当一个被关闭的 channel
中已经发送的数据都被成功接收后，后续的接收操作将不再阻塞，它们会立即返回一个零值。     

没有办法直接测试一个 channel 是否被关闭，但是接收操作有一个变体形式：它多接收一个结果，
多接收的第二个结果是一个布尔值 ok，true 表示成功从 channels 接收到值，false 表示
channels 已经被关闭并且里面没有值可接收。     

```go
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
```     

上面的语法太过笨拙，因此 Go 语言允许 `range` 循环可直接在 channel 上迭代，当 channel 被关闭
且没有值可接收时跳出循环。     

```go
func main() {
	naturals := make(chan int)
	squares := make(chan int)

	go func() {
		for x:= 0; x < 100; x++ {
			naturals <- x
		}
		close(naturals)
	}()

	go func() {
		for x := range naturals {
			squares <- x * x
		}
		close(squares)
	}()

	for x := range squares {
		fmt.Println(x)
	}
}
```    

试图重复关闭一个 channel 将导致 panic 异常，关闭一个 nil 值的 channel 也将导致 panic
异常。关闭一个 channel 还会触发一个广播机制。    

### 8.4.3 单方向的 channel    

当一个 channel 作为一个函数参数时，它一般总是被专门用于只发送或者只接收。为了表明这种意图并防止
被滥用，Go 语言的类型系统提供了单方向的 channel 类型，分别用于只发送或只接收的 channel。类型
`chan <- int` 表示一个只发送 int 的channel，只能发送不能接收，相反，类型 `<- chan int` 
表示一个只接收 int 的channel，只能接收不能发送。     

因为关闭操作只用于断言不再向 channel 发送新的数据，所以只有在发送者所在的 goroutine 才能
调用 `close` 函数，对一个只接收的 channel 调用 `close` 会发生编译错误。   

任何双向 channel 向单向 channel 变量的赋值操作都将导致该隐式转换。    

### 8.4.4 带缓存的 channels

带缓存的 channel 内部持有一个元素队列。队列的最大容量是在调用 `make` 函数创建 channel 时通过
第二个参数指定的。    

带缓存的 channel 的发送操作就是向内部缓存队列的尾部插入元素，接收操作则是从队列的头部删除元素。
如果内部缓存队列是满的，那么发送操作将阻塞直到因另一个 goroutine 执行接收操作而释放了新的队列
空间。相反，如果 channel 是空的，接收操作将阻塞直到有另一个 goroutine 执行发送操作而向队列
插入元素。     

使用 `cap` 函数可以得到 channel 内部缓存的容量。`len` 函数返回内部缓存队列中有效元素的个数。    

`cap(ch)`     

## 8.7 基于 select 的多路复用

下面的程序会进行火箭发射的倒计时。`time.Tick` 函数返回一个 channel，程序会周期性地像一个
节拍器一样向这个 channel 发送事件。每一个事件的值是一个时间戳。   

```go
func main() {
	fmt.Println("Commencing countdown")
	tick := time.Tick(1 * time.Second)
	for countdown := 10; countdown > 0; countdown-- {
		fmt.Println(countdown)
		fmt.Println(<- tick)
	}
}
```     

现在我们让这个程序支持在倒计时中，用户按下 return 键时直接中断发射流程。     

```go
abort := make(chan struct{})
go func() {
	os.Stdin.Read(make([]byte, 1))
	abort <- struct{}{}
}()
```      

现在每一次计数循环的迭代都要等待两个 channel 中的其中一个返回事件了。我们无法做到从每一个
channel 中接收信息，如果我们这么做的话，如果第一个 channel 中没有事件发过来那么程序就会立刻
被阻塞，这样我们就无法收到第二个 channel 中发过来的事件。这时候就需要多路复用了。    

```go
select {
	case <- ch1:
		// ...
	case x := <- ch2:
		// ...
	case ch3 <- y:
		// ...
	default:
		// ...
}
```    

每一个 case 代表一个通信操作并且会包含一些语句组成的一个语句块。select 会等待 case 中有能够
执行的 case 时去执行。当条件满足时，select 才会去通信并执行 case 之后的语句；这时候其他通信
是不会执行。一个没有任何 case 的 select 语句写作 `select{}`，会永远地等待下去。   

```go
ch := make(chan int, 1)
for i := 0; i < 10; i++ {
	select {
		case x := <-ch
			fmt.Println(x)     // 0 2 4 6 8
		case ch <- i
	}
}
```     

如果多个 case 同时就绪时，`select` 会随机地选择一个执行，这样来保证每一个 channel 都有平等
的被 select 的机会。     

有时候我们希望能够从 channel 中发送或者接收值，并避免因为发送或者接收导致的阻塞，尤其是当 channel
没有准备好写或者读时。select 语句就可以实现这样的功能。select 会有一个default 来设置当其他
的操作都不能够被马上处理时程序需要执行哪些逻辑。     

