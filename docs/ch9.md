# 第九章 基于共享变量的并发

## 9.1 竞争条件

数据竞争会在两个以上的 goroutine 并发访问相同的变量且至少其中一个为写操作时发生。根据上述定义，
有三种方式可以避免数据竞争：   

+ 第一种方法是不要去写变量。
+ 第二种方法是避免从多个 goroutine 访问变量。
+ 第三种方法是允许多个 goroutine 去访问变量，但是同一时刻最多只有一个 goroutine 在访问。

## 9.2 sync.Mutex 互斥锁

我们可以用一个容量只有 1 的 channel 来保证最多只有一个 goroutine 在同一时刻访问一个共享变量。
一个只能为 1 和 0 的信号量叫做二元信号量。    

```go
var (
  sema = make(chan struct{}, 1)
  balance int
)

func Deposit(amount int) {
  sema <- struct{}{}
  balance = balance + amount
  <- sema
}

func Balance() int {
  sema <- struct{}{}
  b := balance
  <- sema
  return b
}
```     

这种互斥很使用，而且被 sync 包里的 Mutex 类型直接支持。它的 Lock 方法能够获取到 token，并且
Unlock 方法会释放这个 token：    

```go
import "sync"

var (
  mu sync.Mutex
  balance int
)

func Deposit(amount int) {
  mu.Lock()
  balance = balance + amount
  mu.Unlock()
}

func Balance() int {
  mu.Lock()
  b := balance
  mu.Unlock()
  return b
}
```     

上面的程序例证了一种通用的并发模式。一系列的导出函数封装了一个或多个变量，那么访问这些变量唯一的
方式就是通过这些函数来做。每个函数在一开始就获取互斥锁并在最后释放锁，从而保证共享变量不会被并发
方法。这种函数、互斥锁和变量的编排叫做监控 monitor。     

## 9.3 sync.RWMutex

在上面的例子中，Balance 函数只需要读取变量的状态，所以我们同时让多个 Balance 调用并发执行事实
上是安全的，只要在运行的时候没有存款或取款的操作就行。这种情况下我们需要一种特殊类型的锁，其允许
多个只读操作并行执行，但写操作会完全互斥。这种锁叫做”多读单写“锁，Go 语言中这样的锁是 `sync.RWMutex`:   

```go
var mu sync.RWMutex
var balance int
func Balalce () int {
  mu.RLock()
  defer mu.RUnlock()
  return balance
}
```     

## 9.4 内存同步

我们可能会纠结为什么 Balance 方法需要用到互斥条件，毕竟和存款不一样，它只由一个简单的操作组成，
所以不会碰到其它 goroutine 在其执行”中“执行其它逻辑的风险。    

原因呢是，在现代计算机中处理器上都有缓存，对内存的写入可能会在每一个处理器中缓冲，并在必要时一起
flush 到主存。这种情况下这些数据可能会以与当初 goroutine 写入顺序不同的顺序被提交到主存。像
channel 通信或者互斥量操作这样的原语会使处理器将其聚集地写入 flush 并commit。    

## 9.5 sync.Once 初始化

如果初始化成本比较大的话，那么将初始化延迟到需要的时候再去做就是一个比较好的选择。    

略。   

## 9.6 竞争条件检测

只要在 `go build, go run` 或者 `go test` 命令后加 --race 标志，就会使编译器创建一个你的
应用的修改版或者一个附带了能够记录所有运行期对共享变量访问工具的 test，并且会记录下每一个读或者
写共享变量的 goroutine 的身份信息。      

竞争检查器会检查这些事件，会寻找在哪一个 goroutine 中出现了这样的 case，例如其读或写了一个
共享变量，这个共享变量是被另一个 goroutine 在没有进行干预同步操作便直接写入的。这种情况也就
表明了是对一个共享变量的并发访问，即数据竞争。这个工具会打印一份报告，内容包含变量身份，读取和
写入 goroutine 中活跃的函数的调用栈。    

## 9.8 goroutines 和线程

### 9.8.1 动态栈

每一个 OS 线程都有一个固定大小的内存块（一般是2M)）来做栈，这个栈会用来存储当前正在被调用会挂起
的函数的内部变量。这个固定大小的栈同时很大又很小。     

相反，一个 goroutine 会以一个很小的栈开始其生命周期，一般只需要 2KB。一个 goroutine 的栈，
和操作系统线程一样，会保存其活跃或挂起的函数调用的本地变量，但是和 OS 线程不太一样的是一个 goroutine
的栈大小不固定，可以动态地伸缩，最大1GB。     

### 9.8.2 goroutine 调度

Go 的运行时包含了其自己的调度器，这个调度器使用了一些技术手段，比如m:n 调度，因为其会在 n 个
操作系统线程上调度 m 个goroutine。和系统的线程调度不同，Go 调度器并不是用一个硬件定时器而
是被 Go 本身进行调度的。例如当一个 goroutine 调用了 time.Sleep 或者被 channel 调用或者
mutex 操作阻塞时，调用器会使其进入休眠并开始执行另一个 goroutine 知道时机到了再去唤醒第一个
goroutine。    

### 9.8.3 GOMAXPROCS

go 的调度器使用了一个叫做 GOMAXPROCS 的变量来决定会有多少个操作系统的线程同时执行 Go 的代码。
其默认值是 CPU 的核心数。    


