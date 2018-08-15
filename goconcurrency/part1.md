# 第 1 章 初识 Go 语言

## 1.2 安装和设置

在解压 go 的安装包后，其中包含以下的主要的文件夹：   

+ **api**。用于存放依照 Go 版本顺序的 API 增量列表文件。这些 API 增量列表文件用于 Go
语言 API 检查
+ **bin**。用于存放主要的标准命令文件
+ **blog**。用于存放官方博客中的所有文章
+ **doc**。用于存放标准库的 HTML 格式的文档
+ **lib**。用于存放一些特殊的库文件
+ **misc**。用于存放一些辅助类的说明和工具
+ **pkg**。用于存放安装 Go 标准库后的所有归档文件。其中有类似 linux_amd64 的文件夹，称
为平台相关目录。/pkg/tool/linux_amd64 中存放了很多 Go 的命令和工具
+ **src**。用于存放 Go 自身、Go 标准工具以及标准库的所有源码文件。
+ **test**    

# 第 2 章 语法概览

## 2.1 基本构成要素

Go 的语言符号又称为词法元素，共包括 5 类内容——标识符、关键字、字面量、分隔符和操作符，
它们可以组成各种表达式和语句，而后者都无需以分号结尾。   

### 2.1.1 标识符

标识符可以表示程序实体，前者即为后者的名称。在一般情况下，同一个代码块中不允许出现同名的
程序实体。    

### 2.1.2 关键字

关键字是指被编程语言保留的字符序列，编程人员不能把它们用作标识符。   

### 2.1.3 字面量

简单来说，字面量就是值的一种标记法。    

### 2.1.4 操作符

操作符，也称运算符，它是用于执行特定算术或逻辑操作的符号，操作的对象称为操作数。    

+ 逻辑操作符：`||, &&, !`
+ 比较操作符：`==, !=, <, <=, >, >=`
+ 算术运算符：`+, -, |, ^, *, /, %, <<, >>, &, &^`
+ 取址操作符：`&`
+ 接收操作符：`<-`    

最后需要注意的是，++和-- 是语句而不是表达式，因而它们不存在于任何操作符优先级层次之内。
例如，表达式 *p-- 等同于 (\*p)--。    

### 2.1.5 表达式

表达式是把操作符和函数作用于操作数的计算方法。Go 中的表达式有很多种。    


种类 | 用途 | 示例
---------|----------|---------
 选择表达式 | 选中一个值中的字段或方法 | context.Speaker
 索引表达式 | 选取数组、切片、字符串或字典值中的某个元素 | array[1]
 切片表达式 | 选取数组、数组指针、切片或字符串值中的某个范围的元素 | slice[0:2]
 类型断言 | 判断一个接口值的实际类型是否为某个类型，或一个非接口值的类型是否实现了某个接口类型 | v1.(I1)
 调用表达式 | 调用一个函数或一个值的方法 | v1.M1()

## 2.2 基本类型

只有基本类型及其别名类型才可以作为常量的类型。    

## 2.3 高级类型

数组类型的零值一定是一个不包含任何元素的空数组。   

# 第 3 章 并发编程综述

## 3.1 并发编程基础

并发程序内部会被划分为多个部分，每个部分都可以看做一个串行程序，在这些串行程序之间，可能
会存在交互的虎丘。比如，多个串行程序可能都要对一个共享的资源的进行访问。又比如，它们需要
相互传递一些数据。在这些情况下，我们就需要协调它们的执行，这就涉及同步。同步的作用是避免
在并发访问共享资源时可能发生的冲突，以及确保有条不紊地传递数据。    

根据同步的原则，程序如果想使用一个共享资源，就必须先请求该资源并获取到对它的访问权。当程序
不再需要某个资源的时候，它应该放弃对该资源的访问权。    

传递数据是并发程序内部的另一种交互方式，也称为并发程序内部的通信。实际上，协调这种内部
通信的方式不只“同步”这一种。我们也可以使用异步的方式对通信进行管理，这种方式使得数据可以
不加延迟地发送给数据接收方。即使数据接收方还没有为接收数据做好准备，也不会造成数据发送方
的等待。数据会被临时存放在一个称为通信缓存的数据结构中。    

## 3.2 多进程编程

在多进程程序中，如果多个进程之间需要协作完成任务，那么进程间通信的方式就是需要重点考虑的
事项之一。这种通信常被叫做 IPC。    

在 Linux 操作系统中可以使用的 IPC 方法有很多种。从处理机制的角度看，它们可以分为三大类：
**基于通信** 的 IPC 方法、**基于信号** 的 IPC 方法以及 **基于同步** 的 IPC 方法。其中，基于通信的 IPC
方法又分为以数据传送为手段的 IPC 方法和以共享内存为手段的 IPC 方法，前者包括了管道和
消息队列。管道可以用来传送字节流，消息队列可以用来传送结构化的消息对象。以共享内存为手段
的 IPC 方法主要以共享内存区为代表。基于信号的 IPC 方法就是我们常说的操作系统的信号机制，它
是唯一一种异步 IPC 方法。在基于同步的 IPC 方法中，最重要的就是信号量。    

Go 支持的 IPC 方法有管道、信号和 socket。    

### 3.2.1 管道

管道是一种半双工（或者说单向）的通信方式，只能用于父进程与子进程以及同祖先的子进程之间的
通信。例如，在使用 shell 命令的时候，常常会用到管道：   

`$ ps aux | grep go`    

管道的优点在于简单，而缺点则是只能单向通信以及对通信双方关系上的严格限制。    

对于管道，Go 是支持的。通过标准库包 os/exec 中的 API，我们可以执行操作系统命令并在此之上
建立管道。下面创建一个 `exec.Cmd` 类型的值：   

`cmd0 := exec.Command("echo", "-n", "My first command comes from golang.")`    

对应的 shell 命令：`echo -n "My first command comes from golang."`    

在 `exec.Cmd` 类型之上有一个名为 `Start` 的方法，可以使用它启动命令：   

```go
if err := cmd0.Start(); err != nil {
  fmt.Printf("Error: The command No.0 can not be startup: %s\n", err)
  return
}
```    

为了创建一个能够获取此命令的输出管道，需要在 `if` 语句之前加入如下语句：   

```go
stdout0, err := cmd0.StdoutPipe()
if err != nil {
  fmt.Printf("Error: Couldn't obtain the stdout pipe for command No.0: %s\n", err)
  return
}
```    

变量 cmd0 的 `StdoutPipe` 方法会返回一个输出管道，这里把代表这个输出管道的值赋给了变量
stdout0. stdout0 的类型是 `io.ReadCloser`。    

有了 stdout0，启动上述命令之后，就可以通过调用它的 `Read` 方法来获取命令的输出：   

```go
output0 := make([]byte, 30)
n, err := stdout0.Read(output0)
if err != nil {
  fmt.Printf("Error: Couldn't read data from the pipe: %s\n", err)
  return
}
fmt.Printf("%s\n", output0[:n])
```   

完整的程序：   

```go
package main

import (
	"fmt"
	"os/exec"
)

func main() {
	cmd0 := exec.Command("echo", "-n", "My first command comes from golang")

	stdout0, err := cmd0.StdoutPipe()
	if err != nil {
		fmt.Printf("Error: Couldn't obtain the stdout pipe for command No.0: %s\n", err)
		return
	}

	if err := cmd0.Start(); err != nil {
		fmt.Printf("Error: The command No.0 cannot be startup: %s\n", err)
		return
	}

	output0 := make([]byte, 40)
	n, err := stdout0.Read(output0)
	if err != nil {
		fmt.Printf("Error: Couldn't read data from the pipe: %s\n", err)
		return
	}

	fmt.Printf("%s\n", output0[:n])
}
```   

Go 管道可以把一个命令的输出作为另一个命令的输入，Go 代码也可以做到这一点：   

```go
cmd1 := exec.Command("ps", "aux")
cmd2 := exec.Command("grep", "apipe")
```   

首先，设置 cmd1 的 Stdout 字段，然后启动 cmd1：   

```go
var outputBuf1 bytes.Buffer
cmd1.Stdout = &outputBuf1

if err := cmd1.Start(); err != nil {
  log.Fatalf("Error: The first command can not be startup: %s\n", err)
}

if err := cmd1.Wait(); err != nil {
  log.Fatalf("Error: Couldn't wait for the first command: %s\n", err)
}
```   

接下来，再设置 cmd2 的 Stdin 和 Stdout 字段，启动 cmd2:   

```go
cmd2.Stdin = &outputBuf1
var outputBuf2 bytes.Buffer
cmd2.Stdout = &outputBuf2
if err := cmd2.Start(); err != nil {
  log.Fatalf("Error: The second command can not be startup: %s\n", err)
}

if err := cmd2.Wait(); err != nil {
  log.Fatalf("Error: Couldn't wait for the second command: %s\n", err)
}

fmt.Printf("%s\n", outputBuf2.Bytes())
```   

这里其实好像并没有用到管道相关的 API 啊，只是一种重定向的方案吧。   

这个程序不知道哪出了错，因为把 grep 的参数由 apipe 换成 bash 就可以，但是 apipe 会在第二个
命令的等待时候报错。    

上面所讲的管道也叫做匿名管道，与此相对的是命名管道。与匿名管道不同的是，任何进程都可以通过
命名管道交换数据。实际上，命名管道以文件的形式存在于文件系统中，使用它的方式与文件很类似。
Linux 支持通过 shell 命令创建和使用命名管道：   

```shell
$ mkfifo -m 644 myfifo1
$ tee dst.log < myfifo1 &
[1] 3456
$ cat src.log > myfifo1
```    

命名管道默认是阻塞式的，只有在对这个命名管道的读操作和写操作都已准备就绪之后，数据才开始
流转。    

在 os 包中，包含了可以创建这种独立管道的 API：   

```go
reader, writer, err := os.Pipe()
```   

函数 `os.Pipe()` 的结果中，第一个结果值是代表了该管道输出端的 `*os.File` 类型的值，而
第二个结果值则代表了该管道输入端的 `os.File` 类型的值。    

### 3.2.2 信号

![linux-signal](https://raw.githubusercontent.com/temple-deng/markdown-images/master/other/linux-signal.png)  

Linux 支持的信号有 62 种，其中，编号从 1 到 31 的信号属于标准信号（也称不可靠信号），而
编号从 34 到 64 的信号属于实时信号（也称为可靠信号）。对于同一个进程来说，每种标准信号
只会被记录并处理一次。并且，如果发送给某一个进程的标准信号的种类有多个，那么它们的处理
顺序也是完全不确定的。而实时信号解决了标准信号的这两个问题，即多个同种类的实时信号都可以
记录在案，并且它们可以按照信号的发送顺序被处理。   

Linux 对每一个标准信号都有默认的操作方式，针对不同种类的标准信号，其默认的操作方式一定
会是以下操作之一：终止进程、忽略该信号、终止进程并保存内存信息、停止进程、恢复进程。    

Go 命令会对其中的一些以键盘输入为来源的标准信号做出响应，这是通过标准库代码包 os/signal
中的一些 API 实现的。更具体地讲，Go 命令指定了需要被处理的信号并用一种很优雅的方式来监听
信号的到来。    

```go
type Signal interface{
  String() string
  Signal()  // to distinguish from other Stringers
}
```    

从接口的声明可知，其中的 Signal 方法的声明并没有实际意义。它只是作为 os.Signal 接口类型
的一个标识。   

在 Go 标准库中，已经包含了与不同操作系统的信号相对应的程序实体。具体来说，标准库代码包
syscall 中有与不同操作系统所支持的每一个标准信号对应的同名常量。这些信号常量的类型都是
`syscall.Signal`。`syscall.Signal` 是 `os.Signal` 接口的一个实现类型，同时也是一个 int
类型的别名类型。也就是说，每一个信号常量都隐含着一个整数值，并且都与它所表示的信号在所属
操作系统中的编号一致。    

代码包 os/signal 中的 `Notify` 函数用来当操作系统发送指定信号时发出通知，先来看该函数的
声明：   

`func Notify(c chan<- os.Signal, sig ...os.Signal)`    

`signal.Notify` 函数会把当前进程接收到的指定信号放入参数 c 代表的通道中，这样该函数的
调用方就可以从这个 signal 接收通道中按顺序获取操作系统发来的信号并进行相应的处理。    

第二个参数是一个可变长的参数，这意味着我们再调用 `signal.Notify` 函数时，可以在第一个参数
值之后再附加任意个 `os.Signal` 类型的参数值。参数 sig 代表的参数值包含我们希望自行处理
的所有信号。接收到需要自行处理的信号后，os/signal 包中的程序会把它封装成 `syscall.Signal`
类型的值并放入到 signal 接收通道中。    

```go
sigRecv := make(chan os.Signal, 1)
sigs := []os.Signal{syscall.SIGINT, syscall.SIGQUIT}
signal.Notify(sigRect, sigs...)
for sig := range sigRecv {
	fmt.Printf("Received a signal: %s\n", sig)
}
```    

那这里看， `Notify` 更像是一个绑定监听函数的函数，我们调用这个函数，传入我们想要接收信号
值的通道，以及我们想要监听的信号。然后可能操作系统在收到我们想要监听的信号时，会以某种方式
将信号放入通道中，那这时如果我们正在尝试从通道中取出信号时，就可以获取到信号了。   

完整代码：   

```go
package main

import (
	"os"
	"fmt"
	"os/signal"
	"syscall"
)

func main() {
	sigRecv := make(chan os.Signal)
	sigs := []os.Signal{syscall.SIGINT, syscall.SIGQUIT}
	signal.Notify(sigRecv, sigs...)
	for sig := range sigRecv {
		fmt.Printf("Received siganl: %s\n", sig)
	}
}
```   

还真的可以运行，发送 SIGINT, SIGQUIT 的时候进程不会退出。    

除了自行处理信号之外，还可以在之后的任意时刻恢复对它们的系统默认操作，这需要用到 os/signal
包中的 Stop 函数，其声明如下：   

```go
func Stop(c chan<- os.Signal)
```    

函数 siganl.Stop 会取消掉在之前调用 signal.Notify 函数时告知 signal 处理程序需要自行
处理的若干信号的欣慰。只有把当初传递给 signal.Notify 函数的那个 signal 接收通道作为调用
signal.Stop 函数时的参数值，才能如愿以偿地取消掉之前的行为。   

```go
sigRecv := make(chan os.Signal, 1)
signal.Notify(sigRecv)
for sig := range sigRecv {
	fmt.Printf("Received a signal: %s\n", sig)
}
```    

### 3.2.3 socket

在 Linux 系统中，存在一个名为 `socket` 的系统调用，其声明如下：   

`int socket(int domain, int type, int protocol)`    

三个参数分别是通信域、类型和所用协议。    

每个 socket 都必将存在于一个通信域当中，而通信域决定了该 socket 的地址格式和通信范围：   


通信域 | 含义 | 地址形式 | 通信范围
---------|----------|---------|---------
 AF_INET | IPv4域 | IPv4地址 | 在基于 IPv4 协议的网络中任意两台计算机之上的两个应用程序
 AF_INET6 | IPv6域 | IPv6地址 | 在基于 IPv6 协议的网络中任意两台计算机之上的两个应用程序
 AF_UNIX | UNIX 域 | 路径名称 | 在同一台计算机上的两个应用程序    

socket 的类型有很多，包括 SOCK_STREAM, SOCK_DGRAM, 更底层的 SOCK_RAW，以及针对某个新兴
数据传输技术的 SOCK_SEQPACKET。    

在调用系统调用 socket 的时候，一般会把 0 作为它的第三个参数值，其含义是让操作系统内核根据
第一个参数和第二个参数的值自行决定所使用的协议。   

![socket](https://raw.githubusercontent.com/temple-deng/markdown-images/master/other/socket.png)  

`func Listen(net, laddr string) (Listener, error)`   


Go 中 `net.Listen` 函数用于获取监听器。在建立监听器之后，就可以等待客户端的连接请求了：   

`conn, err := listener.Accept()`    

当调用监听器的 `Accept` 方法时，流程会被阻塞，直到某个客户端程序与当前程序建立 TCP 连接。
参数 `conn` 代表了当前 `TCP` 连接的 `net.Conn` 类型。    

首先需要说明的是，Go 的 socket 编程 API 程序在底层获取的是一个非阻塞式的 socket 实例，
这意味着在该实例上的数据读取操作也都是非阻塞式的。在应用程序试图通过系统调用 `read` 从
socket 的接收缓冲区中读取数据时，即使接收缓冲区中没有任何数据，操作系统内核也不会使系统
调用 read 进入阻塞状态，而是直接返回一个错误码为 EAGAIN 的错误。但是，应用程序并不应该
视此为一个真正的错误，而是应该忽略它，然后稍等片刻再去尝试读取。如果在读取数据的时候接收缓冲
区有数据，那么系统调用 read 就会携带这些数据立即返回。即使当时的接收缓冲区中只包含了一个
字节的数据，也会是这样。    

另一方面，在应用程序试图向 socket 的发送缓冲区中写入一段数据时，即使发送缓冲区已满，系统
调用 write 也不会被阻塞，而是直接返回一个错误码为 EAGAIN 的错误。同样应用程序应该忽略该
错误并稍后再尝试写入数据。   

`net.Conn` 是一个接口类型，定义了如下的方法：   

+ Read: `Read(b []byte) (n int, err error)` 从 socket 的接收缓冲区中读数据。    

如果 socket 编程 API 程序在从接收缓冲区中读取数据类型时发现 TCP 连接已经被另一端关闭了，就
会立刻返回一个 error 类型值。这个 error 类型值与 io.EOF 变量的值是相等的。    

```go
var dataBuffer bytes.Buffer
b := make([]byte, 10)
for {
	n, err := conn.Read(b)
	if err != nil {
		if err == io.EOF {
			fmt.Println("The connection is closed.")
			conn.Close()
		} else {
			fmt.Printf("Read Error: %s\n", err)
		}
		break
	}

	dataBuffer.Write(b[:n])
}
```    

后面的不想写了。    



