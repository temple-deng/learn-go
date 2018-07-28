# Go In Action

<!-- TOC -->

- [Go In Action](#go-in-action)
- [第 2 章 快速开始一个 Go 程序](#第-2-章-快速开始一个-go-程序)
- [第 3 章 打包和工具链](#第-3-章-打包和工具链)
- [第 4 章 数组、切片和映射](#第-4-章-数组切片和映射)
- [第 5 章 Go 语言的类型系统](#第-5-章-go-语言的类型系统)
  - [第 6 章 并发](#第-6-章-并发)
  - [6.1 并发与并行](#61-并发与并行)
  - [6.2 goroutine](#62-goroutine)
  - [6.4 锁住共享资源](#64-锁住共享资源)
    - [6.4.1 原子函数](#641-原子函数)
    - [6.4.2 互斥锁](#642-互斥锁)
- [第 7 章 并发模式](#第-7-章-并发模式)
  - [7.1 runner](#71-runner)
- [第 8 章 标准库](#第-8-章-标准库)
  - [8.1 记录日志](#81-记录日志)
  - [输入和输出](#输入和输出)

<!-- /TOC -->

# 第 2 章 快速开始一个 Go 程序

在 Go 语言里，标识符要么从包里公开，要么不从包里公开。以小写字母开头的标识符是不公开的，不能被
其他包中的代码直接访问。但是，其他包可以间接访问不公开的标识符。例如，一个函数可以返回一个未公开
类型的值，那么这个函数的任何调用者，哪怕调用者不是在这个包里声明的，都可以访问这个值。    

命名接口的时候，也需要遵守 Go 语言的命名惯例。如果接口类型只包含一个方法，那么这个类型的名字以
er 结尾。如果接口类型内部声明了多个方法，其名字需要与其行为关联。     

# 第 3 章 打包和工具链

每个包可以包含任意多个 init 函数，这些函数都会在程序开始执行的时候被调用。以数据库驱动为例，
database 下的驱动在启动时执行 init 函数会将自身注册到 sql 包里，因为 sql 包在编译时并不
知道这些驱动的存在，等启动之后 sql 才能调用这些驱动。     

```go
package postgres

import (
  "database/sql"
)

func init() {
  sql.Register("postgres", new(PostgresDriver))
}
```   

`go vet` 命令会帮开发人员检测代码的常见错误。其会捕获以下类型的错误：   

+ Printf 类函数调用时，类型匹配错误的参数
+ 定义常用的方法时，方法签名的错误
+ 错误的结构标签
+ 没有指定字段名的结构字面量     

`godoc -http=:6060` 会启动 Web 服务器，包含所有的 Go 标准库和 GOPATH 下的源代码的文档。   

为了在 godoc 生成文档里包含自己的代码文档，开发人员需要用下面的规则来写代码和注释：   

+ 用户需要在标识符之前，把自己想要的文档作为注释加入到代码中。这个规则对包、函数、类型和全局
变量都适用
+ 如果想给包写一段文字量比较大的文档，可以在工程里包含一个叫做 doc.go 的文件，使用同样的包名，
并把包的介绍使用注释加在包名声明之前。    

在 godep 和 vendor 这种社区工具已经使用第三方导入路径重写这种特性解决了依赖问题。其思想是把
所有依赖包复制到工程代码库的目录里，然后使用工程内部的依赖包所在目录来重写所有的导入路径。   

# 第 4 章 数组、切片和映射

数组所占用的内存是连续的。由于内存连续，CPU 能把正在使用的数据缓存更久的时间。而且内存连续很容易
计算索引，可以快速迭代数组里的所有元素。数组的类型信息可以提供每次访问访问一个元素时需要在内存中
移动的距离。既然数组的每个元素类型相同，又是连续分配，就可以以固定速度索引数组中的任意数据，速度非常快。    

在 Go 语言里，数组是一个值。这意味着数组可以用在赋值操作中。变量名代表整个数组。因此，同样类型的
数组可以赋值给另一个数组。   

当使用切片字面量时，可以设置初始长度和容量。要做的就是在初始化时给出所需长度和容量作为索引。   

```go
slice := []string{99: ""}
```    

函数 `append` 会智能地处理底层数组的容量增长。在切片容量小于 1000 个元素时，总是会成倍地增加
容量。一旦元素个数超过 1000，容量的增长因子会设为 1.25，也就是会每次增加 25% 的容量。     

在创建切片时（使用切割操作创建），还可以使用我们之前没有提及的第三个索引选项。第三个索引可以用来
控制新切片的容量。其目的不是要增加容量，而是要限制容量。    

```go
source := []{"Apple", "Orange", "Watermelon", "Grape", "Banana"}
slice := source[2:3:4]
// 现在 slice 的长度为 1，容量为2
// 2-3 只包含索引为 2 的元素，所以长度为 1,2-4 包含两个元素，所以容量为 2
```    

在 64 位架构的机器上，一个切片需要 24 字节的内存：指针字段需要 8 字节（地址空间是 64 位的），
长度和容量分别需要 8 字节。由于与切片关联的数据包含在底层数组里，不属于切片本身，所以将切片复制
到任意函数的时候，对底层数组大小都不会有影响。复制时只会复制切片本身，不会涉及底层数组。   

映射的实现使用散列表，映射的散列表包含一组桶。在存储、删除或者查找键值对的时候，所有操作都要先选择
一个同。把操作映射时指定的键传给映射的散列函数，就能选中对应的桶。这个散列函数的目的是生成一个
索引，这个索引最终将键值对分布到所有可用的桶里。     

生成散列值的大致过程是类似这样的：首先使用散列函数将键名转换为一个数值（散列值）。这个数值落在映射
已有桶的序号范围内表示一个可以用于存储的同的序号。之后，这个数值就被用于存储或者查找指定的键值对。
对于 Go 语言的映射来说，生成的散列键的一部分，具有来说是低位被用来选择桶。    

桶内部有两个数据结构，第一个数据结构是一个数组，存储散列键的高八位值，这个数组用来区分每个键值对
要存储在第二个数据结构的哪里。第二个数据结构是一个字节数组，用于存储键值对。该字节数组先一次存储了
这个桶里所有的键，之后依次存储了这个桶里所有的值。     

# 第 5 章 Go 语言的类型系统

结构体字面量有两种形式嘛，一种是按序声明各个字段，没有键名。第二种就是指定每个字段名和字段值嘛，
在这种情况下，字段名和值用冒号分隔，每一行以逗号结尾，也就是在这种情况下，才会出现复合字面量最后
一个字段/索引后面还要加逗号的情况。    

当声明一个引用类型的变量时，创建的变量被称作标头(header)值。从技术细节上来说，字符串也是一种引用
类型。每个引用类型创建的标头值是包含一个底层数据结构的指针。每个引用类型还包含一组独特的字段，用于
管理底层数据结构。因为标头值是为复制而设计的，所以永远不需要共享一个引用类型的值。标头列包含一个指针，
因此通过复制来传递一个引用类型的值的副本，本质上就是在共享底层数据结构。      

如果使用指针接收者来实现一个接口，那么只有指向那个类型的指针才能够实现对应的接口。如果使用值接收
者来实现一个接口，那么那个类型的值和指针都能实现对应的接口。    

那其实简单点说，一个 *user 类型实现了所有 user 方法，也就实现了所有 user 实现的接口。   

之所以指针的方法集只有指针实现的接口才行，是因为我们无法确保总是能获取到一个值的地址，所以值的方法集
只包括了使用者接收者实现的方法。因为这时候我们其实是无法获取到指针的，那也就无法调用到定义在指针
上的方法。     

嵌入类型是将已有的类型直接声明在新的结构类型里。被嵌入的类型被称为新的外部类型的内部类型。   

通过嵌入类型，与内部类型相关的标识符会提升到外部类型上。这些被提升的标识符就像直接声明在
外部类型里的标识符一样，也是外部类型的一部分。这些外部类型就组合了内部类型包含的所有属性，
并且可以添加新的字段和方法。这就是扩展或者修改已有类型的方法。    

```go
package main
import "fmt"

type user struct {
  name string
  email string
}

func (u *user) notify () {
  fmt.Printf("Send user email to %s<%s>\n", u.name, u.email)
}

type admin struct {
  user
  level string
}

func main() {
  ad := admin{
    user: user{
      name: "john smith",
      email: "john@yahoo.com",
    },
    level: "super",
  }

  ad.user.notify()

  // 内部类型的方法可以提升到外部类型
  ad.notify()
}
```    

由于内部类型的提升，内部类型实现的接口也会自动提升到外部类型。    

如果外部类型实现了某些内部类型上的方法，则内部类型的实现就不会被提升。不过内部类型的值
一直存在，因此还可以通过直接访问内部类型的值，来调用没有被提升的内部类型实现的方法。    

```go
// entities/entities.go
package entities

type User struct {
  Name string
  email string
}
```    

```go
package main

import (
  "fmt"
  "../entities"
)

func main() {
  u := entities.User{
    Name: "Bill",
    email: "bill@email.com"
  }


  // 结构字面量中结构 entities.User 的字段 'email' 未知

  fmt.Printf("User: %v\n", u)
}
```     

```go
// entities/entities.go
package entities

type user struct {
  Name string
  Email string
}

type Admin struct {
  user
  Right int
}
```   


```go
package main

import (
  "fmt"
  "../entities"
)

func main() {
  a := entities.Admin{
    Right: 10,
  }

  a.Name = "Bill"
  a.Email = "bill@email.com"

  fmt.Printf("User: %v\n", a)
}
```    

由于内部类型 user 是未公开的，这段代码无法直接通过结构字面量的方式初始化该内部类型。
不过，即便内部类型是未公开的，内部类型里声明的字段依旧是公开的。既然内部类型的标识符提升
到了外部类型，这些公开的字段也可以通过外部类型的字段的值访问。    

## 第 6 章 并发

Go 语言里的并发指的是能让某个函数独立于其他函数运行的能力。当一个函数创建为 goroutine 时，
Go 会将其视为一个独立的工作单元。这个单元会被调度到可用的逻辑处理器上执行。Go 语言运行时
的调度器是一个复杂的软件，能管理被创建的所有 goroutine 并为其分配执行时间。这个调度器
运行在操作系统之上，将操作系统的线程与语言运行时的逻辑处理器绑定，并在逻辑处理器上运行 
goroutine。    

Go 语言的并发同步模型来自一个叫做 **通信顺序进程**(Communicating Sequential Processes, CSP)
的范型。CSP 是一种消息传递模型，通过在 goroutine 之间传递数据来传递消息，而不是对数据进行
加锁来实现同步访问。    

## 6.1 并发与并行

通常来说，操作系统会在物理处理器上调度线程来运行，而 Go 语言的运行时会在逻辑处理器上调度
goroutine 来运行。每个逻辑处理器都分别绑定到单个操作系统线程。    

其实简单来说 goroutine 就是一种内核级线程与用户级线程的混合实现，goroutine 运行在用户
级线程上，而每个用户级线程又绑定到一个内核级线程上。    

![scheduler-manage-goroutine](https://raw.githubusercontent.com/temple-deng/markdown-images/master/go/scheduler-manage-goroutine.png)    

有时，正在运行的 goroutine 需要执行一个阻塞的系统调用，如打开一个文件。当这类调用发生时，
线程和 goroutine 会从逻辑处理器上分离，该线程会继续阻塞，等待系统调用的返回。与此同时，
这个逻辑处理器就失去了用来运行的线程。所以，调度器会创建一个新线程，并将其绑定到该逻辑
处理器上。之后，调度器会从本地运行队列里选择另一个 goroutine 来运行。一旦被阻塞的系统
调用执行完成被返回，对应的 goroutine 会放回本地运行队列，而之前的线程会保存好，以便
之后可以继续使用。     

调度器对可以创建的逻辑处理器的数量没有限制，但语言运行时默认限制每个程最多创建 10000 个
线程。这个限制值可以通过调用 runtime/debug 包的 SetMaxThreads 方法来更改。    

并发并不是并行。并行是让不同代码片段同时在不同的物理处理器上执行。并行的关键是同时做很多
事情，而并发是指同时管理很多事情，这些事情可能只做了一半就被暂停了去做别的事情了。    

![concurrency-parallelism](https://raw.githubusercontent.com/temple-deng/markdown-images/master/go/concurrency-parallelism.png)    

## 6.2 goroutine

```go
package main


import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	runtime.GOMAXPROCS(1)

	var wg sync.WaitGroup
	wg.Add(2)

	fmt.Println("Start Goroutines")

	go func() {
		defer wg.Done()

		for count := 0; count < 3; count++ {
			for char := 'a'; char < 'a' + 26; char++ {
				fmt.Printf("%c ", char)
			}
		}
	}()

	go func() {
		defer wg.Done()

		for count := 0; count < 3; count++ {
			for char := 'A'; char < 'A' + 26; char++ {
				fmt.Printf("%c ", char)
			}
		}
	}()

	fmt.Println("Waiting To Finish")
	wg.Wait()
	fmt.Println("\nTerminating Program")

}
```    

调用 runtime 包的 `GOMAXPROCS` 函数。这个函数允许程序更改调度器可以使用的逻辑处理器的
数量。如果不想在代码里做这个调用，也可以通过修改和这个函数名字一样的环境变量来更改逻辑
处理器的数量。    

`WaitGroup` 是一个计数信号量，可以用来记录并维护运行的 goroutine。如果 `WaitGroup` 的
值大于0，`Wait` 方法就会阻塞。    

基于调度器的内部算法，一个正运行的 goroutine 在工作结束前，可以被停止并重新调度。调度器
这样做的目的是防止某个 goroutine 长时间占用逻辑处理器。当 goroutine 占用时间过长时，
调度器会停止当前正运行的 goroutine，并给其他可运行的 goroutine 运行的机会。    

## 6.4 锁住共享资源

Go 语言提供了传统的同步 goroutine 的机制，就是对共享资源加锁。如果需要顺序访问一个整型
变量或者一段代码，atomic 和 sync 包里的函数提供了很好的解决方案。    

### 6.4.1 原子函数

```go
import (
	"fmt"
	"sync"
	"sync/atomic"
	"runtime"
)

var (
	counter int64
	wg sync.WaitGroup
)

func main() {
	wg.Add(2)

	go incCounter(1)
	go incCounter(2)

	wg.Wait()
	fmt.Println("Final Counter: ", counter)
}

func incCounter(id int) {
	defer wg.Done()

	for count := 0; count < 2; count++ {
		atomic.AddInt64(&counter, 1)

		runtime.Gosched()
	}
}
```    

`AddInt64` 函数会同步整型值的加法，方法是强制同一时刻只能有一个 goroutine 运行并完成
这个加法操作。另外两个有用的原子函数是 `LoadInt64` 和 `StoreInt64`。这两个函数提供了
一种安全地读和写一个整型值的方法。    

但是这里这些怎么保证的同步说的不是清楚，有点不透明。    

### 6.4.2 互斥锁

另一种同步访问共享资源的方式是使用互斥锁。    

略。这个在另一篇文档里有介绍。    

# 第 7 章 并发模式

## 7.1 runner

runner 包用于展示如何使用通道来监视程序的执行时间，如果程序运行时间过长，也可以使用
runner 包来终止程序。    

```go
package runner

import (
	"errors"
	"os"
	"os/signal"
	"time"
)

type Runner struct {
	interrupt chan os.Signal
	complete chan error
	timeout <- chan time.Time
	tasks []func(int)
}

var ErrTimeout = errors.New("received timeout")
var ErrInterrupt = errors.New("received interrupt")

func New(d time.Duration) *Runner {
	return &Runner{
		interrupt: make(chan os.Signal, 1),
		complete: make(chan error),
		timeout: time.After(d),
	}
}

func (r *Runner) Add(tasks ...func(int)) {
	r.tasks = append(r.tasks, tasks...)
}

func (r *Runner) Start() error {
	signal.Notify(r.interrupt, os.Interrupt)

	go func() {
		r.complete <- r.run()
	}()

	select {
	case err := <- r.complete:
		return err
	case <-r.timeout:
		return ErrTimeout
	}
}

func (r *Runner) run() error {
	for id, task := range r.tasks {
		if r.gotInterrupt() {
			return ErrInterrupt
		}

		task(id)
	}

	return nil
}

func (r *Runner) gotInterrupt() bool {
	select {
	case <- r.interrupt:
		signal.Stop(r.interrupt)
		return true
	default:
		return false
	}
}
```    

interrupt 通道收发 `os.Signal` 接口类型的值，用来从操作系统中接收中断时间，os.Signal
接口的声明代码如下：   

```go
type Signal interface {
	String() string
	Signal()
}
```    

通道 interrupt 被初始化为缓冲区容量为 1 的通道。这可以保证通道至少能接收一个来自语言运行时
的 os.Signal 值，确保语言运行时发送这个事件的时候不会被阻塞。    

`time.After` 函数返回一个 `time.Time` 类型的通道。语言运行时会在指定的 duration 时间
到期之后，向这个通道发送一个 time.Time 的值。    

需要注意的是 `gotInterrupt` 中的 `select` 语句，一般来说，`select` 语句在任何要接收的
数据时会阻塞，不过有了 `default` 分支就不会阻塞。`default` 分支会将接收 interrupt 通道
的阻塞调用转变为非阻塞的。如果 interrupt 通道有中断信号需要接收，就会接收并处理这个中断。
如果没有需要接收的信号，就会执行 `default` 分支。    

# 第 8 章 标准库

作为 Go 发布包的一部分，标准库的源代码是经过预编译的。这些预编译后文件，称作归档文件。可以
在 $GOROOT/pkg 文件夹中找到已经安装的各目标平台和操作系统的归档文件。扩展名为 .a 的文件，
就是归档文件。    

## 8.1 记录日志

```go
package main

import (
	"log"
)

func init() {
	log.SetPrefix("Trace: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
}

func main() {
	// Println 写到标准日志记录器
	log.Println("message")

	// Fatalln 在调用 Println() 之后会接着调用 os.Exit(1)
	log.Fatalln("fatal message")

	// Panicln 在调用 Println() 之后会接着调用 panic()
	log.Panicln("panic message")
}
```    

这个程序会在标准输出中产生如下的日志：    

```
Trace: 2018/07/28 11:24:41.560076 E:/go/learn-go/learn-go/goinaction/7/log-main/                                                             main.go:12: message
Trace: 2018/07/28 11:24:41.621238 E:/go/learn-go/learn-go/goinaction/7/log-main/                                                             main.go:15: fatal message
exit status 1
```    

有几个和 log 包关联的标志，这些标志用来控制可以写到每个日志项的其他信息：   

```go
// golang.org/src/log/log.go
const (
	Ldate = 1 << iota

	Ltime

	// 该设置会覆盖 Ltime 标志
	Lmicroseconds

	// 完整的文件名和行号
	Llongfile

	// 最终的文件名和行号
	Lshortfile

	LstdFlags = Ldate | Ltime
)
```    

```go
const (
	Ldate = 1 << iota   // 1 << 0 = 00000001 = 1
	Ltime						    // 1 << 1 = 00000010 = 2
	Lmicroseconds				// 1 << 2 = 00000100 = 4
	Llongfile						// 1 << 3 = 00001000 = 8
	Lshortfile					// 1 << 4 = 00010000 = 16
)
```    

log 包有一个很方便的地方就是，这些日志记录器是多 goroutine 安全的。这意味着在多个 goroutine
可以同时调用来自一个日志记录器的这些函数，而不会有彼此间的写冲突。    

要想创建一个定制的日志记录器，需要创建一个 `Logger` 类型值。可以给每个日志记录器配置一个
单独的目的地，并独立设置其前缀和标志。    

```go
package main

import (
	"log"
	"os"
	"io"
	"io/ioutil"
)

var (
	Trace 	*log.Logger
	Info 		*log.Logger
	Warning *log.Logger
	Error 	*log.Logger 
)

func init() {
	file, err := os.OpenFile("errors.txt",
		os.O_CREATE | os.O.WRONLY | os.O_APPEND, 0666)

	if err != nil {
		log.Fatalln("Failed to open errorlog file: ", err)
	}

	Trace = log.New(ioutil.Discard,
		"TRACE: ",
		log.Ldate | log.Ltime | log.Lshortfile)

	Info = log.New(os.Stdout,
		"INFO: ",
		log.Ldate | log.Ltime | log.Lshortfile)

	Warning = log.New(os.Stdout,
		"Warning: ",
		log.Ldate | log.Ltime | log.Lshortfile)

	Error = log.New(io.MultiWriter(file, os.Stderr),
		"Error: ",
		log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	Trace.Println("I have something standard to say")
	Info.Println("Special Information")
	Warning.Println("There is something you need to know about")
	Error.Println("Someing has failed")
}
```   

为了创建每个日志记录器，我们使用 log 包的 New 函数，它创建并正确初始化一个 Logger 类型
的值。函数 `New` 会返回新创建的值的地址。    

`io.MultiWriter` 函数调用会返回一个 io.Writer 接口类型值，这个值包含之前打开的文件 file，
以及 stderr。MultiWriter 函数是一个变参函数，可以接收个实现了 io.Writer 接口的值。这个
函数会返回一个 io.Writer 值，这个值会把所有传入个 io.Writer 的值绑在一起。当对这个返回
值进行写入时，会向所有绑在一起的 io.Writer 值做写入，这让类似 log.New 这样的函数可以同时
向多个 Writer做输出。    

## 输入和输出

```go
type Writer interface {
	Write (p []byte) (n int, err error)
}
```   

这个接口声明了唯一一个方法 Write，这个方法接收一个 byte 切片，并返回两个值。第一个值
是写入的字节数，第二个值是 error 错误值。     

Write 从 p 里向底层的数据流写入 len(p) 字节的数据。这个方法返回从 p 里写出的字节数
(0 &lt;= n &lt;= len(p))，以及任何可能导致写入提前结束的错误。Write 在返回 n &lt; len(p)
的时候，必须返回某个非 nil 值的 error。Write 绝不能改写切片里的数据。    

```go
type Reader interface {
	Read (p []byte) (n int, err error)
}
```    

1. Read 最多读入 len(p) 字节，保存到 p。这个方法返回读入的字节数(0 &lt;= n &lt;= len(p))
和任何读取时发生的错误。即便 Read 返回的 n &lt; len(p)，方法也可能使用所有的 p 空间
存储临时数据。如果数据可以读取，但是字节长度不足 len(p)，习惯上 Read 会立刻返回可用的
数据，而不等待更多的数据。
2. 当成功读取 n &gt; 0 字节后，如果遇到错误或者文件读取完成，Read 方法会返回读入的字节数。
方法可能会在本次调用返回一个非 nil 的错误，或者在下一次调用时返回错误。这种情况的一个例子是，
在输入的流结束时，Read 会返回非零的读取字节数，可能会返回 err = EOF，也可能返回 err == nil。
无论如何下一次调用 Read 应该返回 0，EOF。



