# Tutorial

<!-- TOC -->

- [Tutorial](#tutorial)
  - [前言](#前言)
  - [1.1 Hello world](#11-hello-world)
  - [1.2 命令行参数](#12-命令行参数)
  - [1.3 查找重复的行](#13-查找重复的行)
  - [1.4 GIF 动画](#14-gif-动画)
  - [1.5 获取 URL](#15-获取-url)
  - [1.6 并发获取多个 URL](#16-并发获取多个-url)
  - [1.7 Web 服务器](#17-web-服务器)
  - [1.8 本章要点](#18-本章要点)
  - [1.9 总结](#19-总结)

<!-- /TOC -->

## 前言

在高级语言中，Go 出现的较晚，因而有一定后发优势，它的基础部分实现得不错：有垃圾收集、包系统、
一等函数、词法作用域、系统调用接口、还有不可变的、默认用 UTF-8 编码的字符串。但相对来说，
它的语言特性不多，而且不太会增加新特性了。比如说，它没有隐式数值类型转换，没有构造和析构函数，
没有运算符重载，没有形参默认值，没有继承，没有泛型，没有异常，没有宏，没有函数注记，没有线程
局部存储。    

这样看的话，Go 还真是只包含了写一个稳定的程序所需要内容的最小子集，而且还摒弃了其他语言中那些
复杂容易出错的地方。    

## 1.1 Hello world

`run` 子命令会将一个或多个以 .go 结尾的源文件进行编译，链接依赖的库，并运行最终的可执行文件。    

```go
package main

import "fmt"

func main() {
  fmt.Println("Hello World")
}
```     

Go 原生支持 Unicode。     

`build` 子命令可以将编译后的可执行文件保存下来。     

Go 的代码都组织在包里，这一点与其他语言中的库或者模块类似。每个包通常包含一个或多个 go 源
代码文件，这些文件都放置在一个定义了这个包做什么事情的文件夹中。每个源代码文件都以 `package`
包声明开始，表明了文件属于哪个包，上面的例子中就是 `package main`。后面跟着一个其导入包的
文件列表。      

`Println` 是 `fmt` 包中一个基本的输出函数，它会打印一到多个值，每个值用空格分隔，结尾是一个
换行符。     

`main` 包比较特殊。它定义了一个独立的可执行程序，而不是一个库。`main` 包中的 `main` 函数也是
比较特殊的，它是整个程序执行的起始点。    

我们必须准确地导入文件中需要的每个包。如果我们少导入了需要的包，或者导入了一些没有用到的包，程序
是无法编译通过的。       

`import` 声明必须跟在`package` 声明之后。之后，程序可能会有函数声明、变量声明、常量声明以及
类型声明。这些声明的顺序大多不重要。  

函数声明包含关键词 `func`、函数的名字，参数列表，一个返回值列表，以及函数的主体。      

Go 语言不要求在声明或语句的结尾加分号，除非一行中出现两个及两个以上的声明和语句。实际上，编译器会主动把
特定符号后的换行符转换为分号，因此换行符添加的位置会影响 Go 代码的正确解析。举个例子，
函数的左括号 `{` 必须和 `func` 函数声明在同一行上，且位于末尾，不能独占一行，而在表达式
`x + y` 中，可在 `+` 后换行，不能在 `+` 前换行。      

Go 语言在代码格式上采取了很强硬的态度。`gofmt` 工具把代码格式化为标准的格式，go 工具的
`fmt` 子命令会对指定包中的所有文件，或者默认情况下是当前目录的所有文件使用 `gofmt` 工具格式化。    

有一个相关的工具 `goimports` 可以根据代码需要，自动地添加或删除 `import` 声明，这个工具并
没有包含在标准的分发包中，可以用下面的命令安装：    

`$ go get golang.org/x/tools/cmd/goimports`     

## 1.2 命令行参数

`os` 包以跨平台的方式，提供了一些与操作系统交互的函数和变量。程序的命令行参数可以 `os` 包的
`Args` 变量获取；`os` 包外部使用 `os.Args` 访问该变量。     

变量`os.Args` 是一个字符串的切片 `slice`。我们现在可以把切片 `s` 理解为一个动态大小的数组，
用 `s[i]` 访问单个元素，用 `s[m:n]`获取子序列。元素的个数可以通过 `len(s)` 获取到。和大部分
语言类似，区间索引时，Go 也采用左闭右开形式，即，区间包括第一个索引元素，不包括最后一个。
如果省略切片表达式的 m 或 n，会默认传入 0 或 `len(s)`。        

`os.Args` 中的第一个元素 `os.Args[0]`，是命令自身；其他的元素则是程序启动时传给它的参数。   

下面是一份 Unix 中 `echo` 命令的实现，这个命令会将其命令行参数打印在一行中。程序导入了两个包。
括号把它们括起来写成列表形式，而没有分开写成独立的 `import` 声明。两种形式都是合法的，但是通常
使用列表形式。导入的顺序并不重要。`gofmt` 工具会将包按名字的字母顺序排序。     

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	var s, sep string
	for i := 1; i < len(os.Args); i++ {
		s += sep + os.Args[i]
		sep = " "
	}
	fmt.Println(s)
}
```    

`var` 声明声明了两个变量 `s` 和 `sep`，都是字符串类型的。一个变量可以在其声明时同时进行
初始化。如果没有显示初始化，那么会被隐式地声明为其类型的**零值**，对于数值类型就是 0，字符串类型就是空字符串 `""`。     

`:=` 符号是 _短变量声明_ 的一部分，短变量声明可以声明一个或多个变量，并基于其初始值给变量
赋予合适的类型。    

`i++` 和 `i--` 都是语句，而不是表达式。这也就意味 `j = i++` 是不合法的，并且只能是
后缀形式，`--i` 也是不合法的。    

Go 语言中只有 `for` 循环这一种循环语句。不过它有很多种形式。上面例子中的形式类似下面这样：    

```go
for initializtion; condition; post {
  // 零到多个语句
}
```     

循环中的三个组件部分不需要使用括号包裹，不过大括号是必须的，并且必须和 `post` 在同一行。    

initialization 语句是可选的，在循环开始前执行。initialization 如果存在，必须是一条
简单语句simple statement，即，短变量声明、自增语句、赋值语句或函数调用。condition 是
一个布尔表达式，其值在每次循环迭代开始时计算。如果为 true 则执行循环体语句。post 在
循环体执行结束后执行，之后再次对 condition 求值。condition 为 false 时，循环结束。      


for 循环的这三个部分每个都可以省略，如果省略了 initialization 和 post，那么分号也就可以省略了：   

```go
// 传统的 while 循环
for condition {

}
```    

如果连 condition 也省略了：    

```go
for {

}
```    

类似于无限循环，但是还是可以用 `break` 或者 `return` 终止循环。     

`for` 循环的另一种形式，在某种数据类型的区间 range 上遍历，如字符串或切片：   

```go
// echo2
func main() {
  var s, sep string
  for _, arg := range os.Args[1:] {
    s += sep + arg
    sep = " "
  }

  fmt.Println(s)
}
```     

每次循环迭代，`range` 产生一对值；索引以及在该索引处的元素值。这个例子不需要索引，
但 `range` 语法要求，要处理元素，必须处理索引。一种思路是把索引赋值给一个临时变量，
如 `temp`，然后忽略它的值，但 Go 语言不允许使用无用的局部变量 local variable，这会导致编译错误。    

Go 语言中这种情况的解决方法是用 _空标识符_，即 `_`。空标识符可用于任何语法需要变量名
但程序逻辑不需要的时候，例如，在循环里，丢弃不需要的循环索引，保留元素值。    

注意上面的循环中，索引还是从 0 开始的。     

声明一个变量有好几种方式，下面这些都等价：    

```go
s := ""
var s string
var s = ""
var s string = ""
```     

第一种形式，短变量声明，是最简洁的，但是只能用在一个函数中，而不能用作包一级的变量声明。
第三种形式用的比较少，除非是声明多个变量一次性。第四种形式显示地声明了变量的类型，当
变量的初始值与变量类型一致的时候，这时是有些冗余的，但是当初始值与变量类型不一致的时候，
这种形式就有必要了。通常来说，我们大多数会使用前两种形式。      

## 1.3 查找重复的行

```go
import (
	"fmt"
	"os"
	"bufio"
)

func main() {
	counts := make(map[string]int)
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		counts[input.Text()]++
	}

	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}
```    

与 `for` 循环类似，`if` 语句 condition 两边也可以不加括号，但是主体的大括号必须有。同时，
`if` 也可以有一个类似  initialization 的部分，同时也是一个简短语句。简短语句中声明的
变量在 if 及 else 中都是可用的。        


`map` 是键值对的集合，并且提供了对集合元素的常数时间的存取以及测试操作。键可以是任何
能够与 `==` 进行比较的类型值，通常都是字符串。值也可以是任何类型。内置的 `make`函数可以
创建一个新的空 map，不过它还有别的作用。     

`map` 遍历的顺序是不确定的，随机的。这种设计是有意为之的，因为能防止程序依赖特定的
遍历顺序，而这是无法保证的。     

`Scanner` 是 `bufio` 包提供的最有用的特性之一，它读取输入并将其拆成行或单词；通常是
处理行形式的输入最简单的方法。     

input 变量从程序的标准输入中读取内容。每次调用 `input.Scan()`，即读入下一行，并移除行
末的换行符；读取的内容可以调用 `input.Text()` 得到。`Scan` 函数在没有输入时返回 false。      

`fmt.Printf` 函数对一些表达式产生格式化输出。该函数的首个参数是个格式字符串，指定后续
参数改如何格式化。每个参数的格式取决于“转换字符”，形式为百分号后跟一个字母。     

下面的表格是部分可用的转换字符：   

```
%d            十进制整数
%x, %o, %b    十六进制、八进制、二进制整数
%f, %g, %e    浮点数: 3.141593 3.1415926.....   3.141593e+00
%t						布尔
%c						字符
%s						字符串
%q						带双引号的字符串"abc"或带单引号的字符'c'   说实话这个看不懂
%v						变量的自然形式
%T						变量的类型
%%						字面上的百分号标志
```     

默认情况下，`Printf`不会换行。按照惯例，以字母 `f` 结尾的格式化函数，如 `log.Printf`
或 `fmt.Errorf`，都采用 `fmt.Printf` 的格式化准则。而以 `ln` 结尾的格式化函数，
则遵循 `Println` 的方式，以跟 `%v` 差不多的方式格式化参数，并在最后添加一个换行符。`f`指 format, `ln` 指 line。      

```go
package main

import (
	"os"
	"fmt"
	"bufio"
)

func main() {
	counts := make(map[string]int)
	files := os.Args[1:]

	if len(files) == 0 {
		countLines(os.Stdin, counts)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				// 错误处理
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			countLines(f, counts)
			f.Close()
		}
	}

	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}

func countLines(f *os.File, counts map[string]int) {
	input := bufio.NewScanner(f)
	
	for input.Scan() {
		counts[input.Text()]++
	}
}
```     

`os.Open` 函数返回两个值。第一个值是被打开的文件，第二个值是内置 `error` 类型的值。
如果 `err` 等于内置值 `nil`，那么文件被成功打开。读取文件，直到文件接收，然后调用
`Close` 关闭该文件，并释放占用的所有资源。相反的话，说明打开文件时出错了。     

函数和包级别的变量可以任意顺序声明，并不影响其调用。    

`map` 是一个由 `make` 函数创建的数据结构的引用。`map` 作为参数传递给某函数时，
该函数接收这个引用的一份拷贝，被调用函数对 `map` 底层数据结构的任何修改，调用者函数都
可以通过持有的 `map` 引用看到。注意是引用的拷贝，所以类似与引用传值。       

`dup` 的前两个版本以“流”模式读取输入，并根据需要拆分成多个行。理论上，这些程序可以处理
任意数量的输入数据。还有另一个方法，就是一个口气把全部输入数据读到内存中，一次分割为多行，
然后处理它们。下面这个版本就是这么操作的，这个例子引入了 `ReadFile` 函数，其读取指定文件的全部内容。      

```go
package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"strings"
)

func main() {
	counts := make(map[string]int)
	for _, filename := range os.Args[1:] {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "dup3: %v\n", err)
			continue
		}

		for _, line := range strings.Split(string(data), "\n") {
			counts[line]++
		}
	}

	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}
```    

`ReadFile` 函数返回字节切片 byte slice，必须把它转换为 `string`，才能用 `strings.Split` 分割。      

在底层实现上，`bufio.Scanner`, `ioutil.ReadFile`和 `ioutil.WriteFile` 都使用
`* os.File` 的 `Read` 和 `Write` 方法，但是，大多数程序员很少需要直接调用那些低级函数。     

## 1.4 GIF 动画

```go
package main

import (
	"image"
	"image/color"
	"image/gif"
	"io"
	"math"
	"math/rand"
	"os"
)
```    

这个例子暂时不解释。     

当我们 `import` 了一个包路径包含有多个单词的包时，比如 `image/color`，通常我们只需要用最后那个单词表示这个包就可以。    

## 1.5 获取 URL

```go
import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	for _, url := range os.Args[1:] {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			os.Exit(1)
		}

		b, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
			os.Exit(1)
		}

		fmt.Printf("%s", b)
	}
}
```    

http.Get 函数是创建 HTTP 请求的函数，如果获取过程没有出错，那么会在 resp 这个结构体
中得到访问的请求结果。resp的Body字段包括一个可读的服务器响应流。ioutil.ReadAll 函数从
response 中读取到全部内容。注意上面的操作好像是流操作。    

## 1.6 并发获取多个 URL

```go
import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"net/http"
	"time"
)

func main() {
	start := time.Now()
	ch := make(chan string)
	for _, url := range os.Args[1:] {
		go fetch(url, ch)
	}

	for range os.Args[1:] {
		fmt.Println(<-ch)
	}

	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}

	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}

	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %7d %s", secs, nbytes, url)
}
```   

goroutine 是一种函数的并发执行方式，而 channel 是用来在 goroutine 之间进行参数传递。
main 函数本身也运行在一个 goroutine中，而 go function 则表示创建一个新的 goroutine，
并在这个新的 goroutine 中执行这个函数。    

main 函数中用 make 函数创建了一个传递 string 类型参数的 channel，对每一个命令行参数，
我们都用 go 这个关键字来创建一个 goroutine，并且让函数在这个 goroutine 异步执行 http.Get 
方法，这个程序里的 io.Copy 会把响应的 Body 内容拷贝到 ioutil.Discard 输出流中。每当
请求返回内容时，fetch 函数都会往 ch 这个 channel 里写入一个字符串，由 main 函数里的
第二个 for 循环来处理并打印channel 里的这个字符串。    

当一个 goroutine 尝试在一个 channel 上做 send 或 receive 操作时，这个 goroutine 会
阻塞在调用处，直到另一个 goroutine 往这个 channel 里写入或者接收值，这样两个 goroutine 
才会继续执行 channel 操作之后的逻辑。在这个例子中，每一个fetch 函数在执行时都会往 channel 
里发送一个值，main 函数负责接收。    

## 1.7 Web 服务器

略。   


## 1.8 本章要点

switch 选择：    

```go
switch coinfilp() {
	case "heads":
		heads++
	case "tails":
		tails++
	default:
		fmt.Println("landed on edge")
}
```   

Go 语言不需要显示地在每一个 case 后写 break，语言默认执行完 case 后的逻辑语句会自动退出。
当然了，如果你想要相邻的几个 case 都执行同一逻辑的话，需要自己显示地写上一个 fallthrough
语句来覆盖这种默认行为。    

Go 语言的 switch 还可以不带操作对象，不带操作对象时默认用 true 值代替，然后将每个 case 的表达式和 true 值进行比较：    

```go
func Signum(x int) int {
	switch {
		case x > 0:
			return +1
		default:
			return 0
		case x < 0:
			return -1
	}
}
```    

像 for 和 if 控制语句一样，switch 也可以紧跟一个简短的变量声明、一个自增语句、赋值语句，或者一个函数调用。    

## 1.9 总结

每个文件的结构：包声明、可选的依赖导入、其他的各种声明。     

短变量声明也是语句 statement，后缀表达式 i--/i++ 也是语句。      

for 循环的 initialization 部分、if 和 switch 的操作数部分必须是一个简短语句 simple
statement，即短变量声明、自增或赋值语句、函数调用。注意 if 的条件部分好像前面也可以
添加一个 initialization 部分，switch 的操作数如果不存在等价于 `switch true`。      

map 的键名的类型必须可以用 == 进行比较。    

Printf 支持的动词：    

```go
%d
%x, %o, %b
%f, %g, %e
%t
%s
%c
%q
%v
%T
```   

以 ln 结尾的格式化函数会将参数以 %v 的形式展示。     

用 make 创建的 map 是对数据结构的引用，所以如果将 map 类型当做参数传入函数，函数参数
其实是一个引用的副本。     

常量的值只能是数值类型、字符串或布尔型。    

表达式 `[]color.Color{....}` 和 `gif.Gif{...}` 叫做复合字面值 composite literals。
这是实例化 Go 语言里复合类型的一种写法。    
