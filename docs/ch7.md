# 第七章 接口

接口类型表达了对其他类型行为的抽象和概括；因为接口类型不会和特定的实现细节绑定在一起，通过
这种抽象的方式我们可以让我们的函数更加灵活和更具有适应能力。     

Go 语言中接口类型的独特之处在于它是满足隐式实现的。也就是说，我们没有必要对于给定的具体
类型定义所有满足的接口类型；简单地拥有一些必需的方法。这种设计可以让你创建一个新的接口类型
满足已经存在的具体类型却不会去改变这些类型的定义。      

## 7.1 接口约定

目前为止，我们看到的类型都是具体的类型。一个具体的类型可以准确的描述它所代表的值并且展示
出对类型本身的一些操作方式，就像数字类型的算术操作、slice 的索引、append 和取范围操作。
具体的类型还可以通过它的方法提供额外的行为操作。总的来说，当你拿到一个具体的类型时，你就
知道它的本身是什么和你可以用它来做什么。     

接口类型时一种抽象的类型。他不会暴露出它所代表的对象的内部值的结构和这个对象支持的基础
操作的集合；它们只会展示出它们自己的方法。也就是说当你看到一个接口类型的值时，你不知道它
是什么，唯一知道的就是可以通过它的方法来做什么。      

`Fprintf` 函数中的第一个参数不是一个文件类型。它是一个 `io.Writer` 类型，这个类型是一个
接口类型，定义如下：    

```go
package io

type Writer interface {
  Write(p []byte) (n int, err error)
}
```    

`io.Writer` 接口类型定义了 `Fprintf` 和其调用者之间的一份约定。一方面，约定要求调用者
提供具体类型的值就像 `*os.File` 和 `*bytes.Buffer`，这些类型都有一个有着特定签名和行为
的叫做 Writer 的方法。另一方面这个约定保证了 Fprintf 接受任何满足 `io.Writer` 接口的值
都可以正常工作。     

给一个类型定义 `String` 方法，可以让它满足最广泛使用的接口类型之一 `fmt.Stringer`:   

```go
package fmt

type Stringer interface {
  String () string
}
```     

感觉接口呢，就是方便我们在使用时，并不指定具体的类型而只指定类型所具有的行为。     

## 7.2 接口类型

接口类型具体描述了一系列方法的集合，一个实现了这些方法的具体类型是这个接口类型的实例。   

`io.Writer` 类型是用的最广泛的接口之一，因为它提供了所有的类型写入 bytes 的抽象，包括文件
类型、内存缓冲区、网络连接、HTTP 客户端、压缩工具、哈希等。io 包中定义了很多其他有用的
接口类型。`Reader` 可以代表任意可以读取 bytes 的类型，`Closer` 可以是任意可以关闭的值。    

```go
package io

type Reader interface {
  Read(p []byte) (n int, err error)
}

type Closer interface {
  Close() error
}
```    

我们可以使用组合已有的接口来定义新的接口类型：    

```go
type ReadWriter interface {
  Reader
  Writer
}

type ReadWriteCloser interface {
  Reader
  Writer
  Closer
}
```    

上面用到的语法和结构内嵌相似，我们可以用这种方式以一个简写命名另一个接口，而不用声明它所有的
方法。这种方式称为接口内嵌。甚至可以使用混合的风格：    

```go
type ReadWriter interface {
  Read (p []byte) (n int, err error)
  Writer
}
```    

## 7.3 实现接口的条件

一个类型如果拥有一个接口需要的所有方法，那么这个类型就实现了这个接口。     

接口指定的规则非常简单：表达一个类型属于某个接口只要这个类型实现了这个接口。所以：    

```go
var w io.Writer
w = os.Stdout
w = new(Bytes.Buffer)
```    

这个规则甚至适用于等式右边本身也是一个接口类型。     

在接下来讲之前，有必要解释一下 **一个类型有一个方法意味着什么**。在6.2 节中讲到，对
于每一个命名的具体类型T；它一些方法的 receiver 是类型 T 本身，而另一些则是一个指向 T 
的指针。使用一个 T 类型参数调用一个 *T 上的方法也是合法的，**只要这个参数是一个变量**。
编译器会隐式地获取它的地址，我觉得之所以必须是变量是因为变量才是可以取地址的，而一个字面量
是没有地址的。这仅仅是一个语法糖。T类型的值不拥有所有 *T 指针的方法，结果就是它可能只实现
了一小部分的接口。    

```go
type IntSet struct { /**/ }
func (*IntSet) String() string
var _ = IntSet{}.String()   // compile error

var s IntSet
var a = s.String()   // ok

var b fmt.Stringer = &s  // ok
var c fmt.Stringer = s   // compile error
```    

`interface{}` 被称为空接口类型，并且是无可或缺的。空接口类型对实现它的类型没有要求，所以
我们可以将任意一个值赋给空接口类型。    

```go
var any interface{}
any = true
any = 123.34
any = "hello"
```    

## 7.5 接口值

概念上讲一个接口类型的值，或者叫接口值，由两个部分组成：一个具体的类型和一个类型的值。它们被称为
接口的动态类型和动态值。对于像 Go 语言这种静态类型的语言，类型是编译器的概念，因此一个
类型不是一个值。在我们的概念模型中，一些提供每个类型信息的值称为类型描述符，比如类型的名称
和方法。在一个接口值中，类型部分代表与之相关类型的描述符。     

下面4个语句中，变量 w 得到了3个不同值。（开始和最后的值是相同的）    

```go
var w io.Writer
w = os.Stdout
w = new(bytes.Buffer)
w = nil
```    

让我们进一步观察在每一个语句后的 w 变量的值和动态行为。第一个语句定义了变量 w:    

`var w io.Writer`     

在 Go 语言中，变量总是被一个定义明确的值初始化，即使接口类型也不例外。对于一个接口的零值
就是它的类型和值的部分都是 nil。    

![nil-interface](https://github.com/temple-deng/learning-repo/tree/master/pics/nil-interface.png)     

一个接口值基于它的动态类型被描述为空或非空，所以这是一个空的接口值。你可以通过使用 w==nil
或者 w!=nil 来判断接口值是否为空。调用一个空接口值上的任意方法都会产生 panic:   

`w.Writer([]byte("hello"))  // panic: nil pointer dereference`     

第二个语句将一个 *os.File 类型的值赋给变量 w:    

`w = os.Stdout`    

这个赋值过程调用了一个具体类型到接口类型的隐式转换，这和显式的使用 `io.Writer(os.Stdout)`
是等价的。这类转换不管是显式的还是隐式的，都会刻画出操作到的类型和值。这个接口值的动态类型
被设为 `*os.File` 指针的类型描述符，它的动态值持有 `os.Stdout` 的拷贝；这是一个代表
处理标准输出的 `os.File` 类型变量的指针。     

![os-file-interface](https://github.com/temple-deng/learning-repo/tree/master/pics/os-file-interface.png)     

调用一个包含 `*os.File` 类型指针的接口值的 Write 方法，使得`(*os.File).Write` 方法被
调用。这个调用输出"hello"。    

`w.Write([]byte("hello"))`      

通常在编译期，我们不知道接口值的动态类型是什么，所以一个接口上的调用必须使用动态分配。因为
不是直接进行调用，所以编译器必须生成代码从类型描述符上获得名字为 Write 的方法的地址，
然后间接调用那个地址。这个调用的接收者是一个接口动态值的拷贝，os.Stdout。效果和下面的
直接调用一样：    

`os.Stdout.Write([]byte("hello"))`     

第三个语句给接口值赋了一个 `*bytes.Buffer` 类型的值：    

`w = new(bytes.Buffer)`    

现在动态类型是 `*bytes.Buffer` 并且动态值是一个指向新分配的缓冲区的指针。    

最后，第四个语句将 nil 赋给了接口值。这个重置将它所有的部分都设为 nil 值。    

接口值可以使用 == 和 != 来进行比较。两个接口值相等仅当它们都是 nil 值或者它们的动态类型相同
并且动态值也根据这个动态类型的 == 操作相等。    

### 7.5.1 警告：一个包含 nil 指针的接口不是 nil 接口

一个不包含任何值的 nil 接口值和一个刚好包含 nil 指针的接口值是不同的。     

## 7.8 error 接口

error 类型实际上就是一个接口类型，这个类型有一个返回错误信息的单一方法：   

```go
type error interface {
  Error() string
}
```   

创建一个 error 最简单的方法就是调用 errors.New 函数，他会根据传入的错误信息返回一个新的
error。整个 errors 包仅只有4行：    

```go
package errors

func New(text string) error {
  return &errorString{text}
}

type errorString struct {
  text string
}

func (e *errorString) Error() string {
  return e.text
}
```     

调用 `errors.New` 函数是非常稀少的，因为有一个方便的封装函数 fmt.Errorf:    

```go
package fmt

import "errors"

func Errorf(format string, args ...interface{}) error {
  return errors.New(Sprintf(format, args...))
}
```    

虽然 *errorString 可能是最简单的错误类型，但远非只有他一个。例如，`syscall` 包提供了
Go 语言底层系统调用 API。在多个平台上，它定义一个实现 error 接口的数字类型 Errno，并且
在 Unix 平台上，Errno 的 Error 方法会从一个字符串表中查找错误消息。     

## 7.10 类型断言

类型断言是一个应用在一个接口值上的操作。语法上它看起来像 x.(T)，x 是一个接口类型的表达式，
T 是一个类型，叫做“断言”类型。类型断言检查它操作对象的动态类型是否和断言的类型匹配。    

这里有两种可能。第一种，如果断言的类型 T 是一个具体类型，然后类型断言检查 x 的动态类型
是否和 T 相同。如果这个检查成功了，类型断言的结果是 x 的动态值，当然它的类型是 T。换句话
说，具体类型的类型断言从它的操作对象中获得具体的是。如果检查失败，接下来这个操作会抛出 panic：    

```go
var w io.Writer
w = os.Stdout
f := w.(*os.File)         // success: f == os.Stdout
c := w.(*bytes.Buffer)    // panic: interface holds *os.File, not *bytes.Buffer
```     


