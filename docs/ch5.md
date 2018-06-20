# 第五章 函数

<!-- TOC -->

- [第五章 函数](#第五章-函数)
  - [5.1 函数声明](#51-函数声明)
  - [5.2 递归](#52-递归)
  - [5.3 多返回值](#53-多返回值)
  - [5.4 错误](#54-错误)
    - [5.4.1 错误处理策略](#541-错误处理策略)
    - [5.4.2 文件结尾错误（EOF）](#542-文件结尾错误eof)
  - [5.5 函数值](#55-函数值)
    - [5.5.1 再谈类型](#551-再谈类型)
  - [5.6 匿名函数](#56-匿名函数)
    - [5.6.1 警告：捕获迭代变量](#561-警告捕获迭代变量)
  - [5.7 可变参数](#57-可变参数)
  - [5.8 Deferred 函数](#58-deferred-函数)
  - [5.9 Panic 异常](#59-panic-异常)
  - [5.10 Recover 捕获异常](#510-recover-捕获异常)

<!-- /TOC -->

## 5.1 函数声明

函数声明包括函数名、形参列表、返回值列表（可省略）以及函数体。    

```go
func name(parameter-list) (result-list) {
  body
}
```    

返回值列表描述了函数返回值的变量名以及类型。如果函数返回一个无名变量或者没有返回值，返回
值列表的括号是可以省略的。如果一个函数声明不包括返回值列表，那么函数体执行完毕后，不会
返回任何值。    

```go
func hypot(x, y float64) float64 {
  return math.Sqrt(x*x + y*y)
}
fmt.Println(hypot(3,4))   // 5
```     

返回值也可以像形式参数一样被命名。在这种情况下，每个返回值被声明成一个局部变量，并根据该
返回值的类型，将其初始化为0。如果一个函数在声明时，包含返回值列表，该函数必须以 return 
语句结尾，除非函数明显无法运行到结尾处。例如函数在结尾时调用了 panic 异常或函数中存在
无限循环。    

正如 hypot 一样，如果一组形参或返回值有相同的类型，我们不必为每个形参都写出参数类型：   

```go
func f(i, j, k int, s, t string) { /**/ }
func f(i int, j int, k int, s string, t string) { /**/}
```    

下面，我们给出过4种方法生命拥有2个 int 型参数和1个 int 型返回值的函数。_可以强调某个
参数未被使用。    

```go
func add(x int, y int) int { return x + y}
func sub(x, y int) (z int) { z=x -y; return}  // 注意这个，看样子返回值可能在函数调用时
                                              // 就被初始化了
func first(x int, _ int) int { return x}
func zero(int, int) int { return 0}
```     


函数的类型被称为函数的标识符。如果有两个函数形式参数列表和返回值列表中的变量类型一一
对应，那么这两个函数被认为有相同的类型和标识符。形参和返回值的变量名不影响函数标识符也
不影响它们是否可以以省略参数类型的形式提示。      

每一次函数调用都必须按照声明顺序为所有参数提供实参。在函数调用时，Go 语言没有默认参数值，
也没有任何方法可以通过参数名指定形参，因此形参和返回值的变量名对于函数调用者而言没有意义。   


在函数体中，函数的形参作为局部变量，被初始化为调用者提供的值。函数的形参和有名函数值作为
函数最外层的局部变量，被存储在相同的词法块中。     

实参通过值的方式传递，因此函数的形参是实参的拷贝。对形参进行修改不会影响实参。但是，如果
实参包括引用类型，如指针、slice、map、function、channel 等类型，实参可能会被函数修改。    

可能会偶尔遇到没有函数体的函数声明，这表示函数不是以 Go 实现。这样的声明定义了函数标识符：    

```go
package math
func Sin(x float64) float
```     

## 5.2 递归

略。    

## 5.3 多返回值

如果一个函数所有的返回值都命名了，那么该函数的return 语句可以省略操作数，这叫做 bare return。    

```go
func CountWordsAndImages(url string) (words, images int, err error) {
  resp, err := http.Get(url)
  if err != nil {
    return
  }
  doc, err := html.Parse(resp.Body)
  resp.Body.Close()
  if err != nil {
    err = fmt.Errorf("parsing HTML: %s", err)
    return
  }
  words, images = countWordsAndImages(doc)
  return
}
```     

按照返回值列表的次序，返回所有的返回值，在上面的例子中，每一个 return 语句等价于：   

`return words, images, err`      

## 5.4 错误

panic 是来自被调用函数的信号，表示发生了某个已知的 bug。     

对于那些将运行失败看作是预期结果的函数，它们会返回一个额外的返回值，通常是最后一个，来
传递错误信息。如果导致失败的原因只有一个，额外的返回值可以是一个布尔值，通常被命名为 ok。    

通常，导致失败的原因不止一种，尤其是对I/O操作而言，用户需要了解更多的错误信息。因此，额外
的返回值不再是简单的布尔类型，而是 error 类型。     

内置的 error 是接口类型。error 类型可能是 nil 或者 non-nil，nil 意味着函数运行成功，
non-nil 表示失败。对于 non-nil 的error 类型，我们可以通过调用 error 的 Error 函数或者
输出函数获得字符串类型的错误信息。    

通常，当函数返回 non-nil 的error 时，其他的返回值是未定义的，这些未定义的返回值应该被
忽略。然而，有少部分函数在发生错误时，仍然会返回一些有用的返回值。比如，当读取文件发生错误
时，Read 函数会返回可以读取的字节数以及错误信息。     

在 Go 中，函数运行失败时会返回错误信息，这些错误信息被认为是一种预期的值而非异常 exception，
这使得 Go 有别于那些将函数运行失败看做是异常的语言。虽然 Go 有各种异常机制，但这些机制
仅被使用在处理那些未被预料到的错误，即bug，而不是在健壮程序中应该被避免的程序错误。    

### 5.4.1 错误处理策略

当一次函数调用返回错误时，调用者应该选择怎样处理错误。根据情况的不同，有很多处理方式，让我们
看看常用的五种方式。     

最常用的方式是传播错误。这意味着函数中某个子程序的失败，会变成该函数的失败。    

第二种策略，如果错误的发生是偶然性的，或由不可预知的问题导致的。一个明智的选择是重新尝试
失败的操作。在重试时，我们需要限制重试的时间间隔或重试的次数，防止无限制的重试。    

如果错误发生后，程序无法继续运行，我们就可以采用第三种策略：输出错误信息并结束程序，需要
注意的是，这种策略只应在 main 中执行。对库函数而言，应仅向上传播错误。     

第四种策略，有时，我们只需要输出错误信息就足够了，不需要中断程序的运行，我们可以通过 log
包提供函数：    

```go
if err := Ping(); err != nil {
  log.Printf("ping failed: %v; networking disabled", err)
}
```    

第五种策略，就是直接忽略掉错误。     

在 Go 中，错误处理有一套独特的编码风格。检查某个子函数是否失败后，我们通常将处理失败的逻辑
代码放在处理成功的代码之前。如果某个错误会导致函数返回，那么成功时的逻辑代码不应放在 else
语句块中，而应直接放在函数体中。      

### 5.4.2 文件结尾错误（EOF）

io 包保证任何由文件结束引起的读取失败都返回同一个错误--io.EOF，错误在 io 包中定义。     

## 5.5 函数值

在 Go 中，函数被看做是第一类值first-class values：函数像其他值一样，拥有类型，可以被
赋值给其他变量，传递给函数，从函数返回。     

```go
func square(n int) int {
	return n*n
}

func negative(n int) int {
	return -n
}

func product(m, n int) int {
	return m*n
}

func main() {
	f := square
	fmt.Println(f(3))    // 9

	f = negative
	fmt.Println(f(3))    // -3
  fmt.Printf("%T\n", f)  // func(int) int
  
  f = product  // compile error: can't assign func(int, int) int to func(int) int
}
```    

函数类型的零值是 nil，调用值为 nil 的函数值会引起 panic 错误：   

```go
var f func(int) int
f(3)
```     

函数值可以与 nil 比较。但是函数值之间是不可比较的。     

### 5.5.1 再谈类型

看情况，数组、结构体、slice、map 和函数，这都是一个笼统的类型名称，简单介绍了具有相同
特征一组类型的集合，像一个数组的具体类型由元素个数和元素类型两部分组成，而函数由其参数
列表和其返回值列表部分决定。    


## 5.6 匿名函数

拥有函数名的函数只能在包级语法块中被声明，通过函数字面量，我们可以绕过这一限制，在任何
表达式中表示一个函数值。函数字面量的语法和函数声明类似，区别在于 func 关键字后没有函数
名。函数值字面量是一种表达式，它的值被称为匿名函数。    

函数字面量允许我们在使用函数时，再定义它：    

```go
strings.Map(func(r rune) rune {return r + 1}, "HAL-9000")
```     

更为重要的是，通过这种方式定义的函数可以访问完整的词法环境，这意味着在函数中定义的内部
函数可以引用该函数的变量，如下所示：     

```go
func squares() func() int {
	var x int
	return func() int {
		x++
		return x*x
	}
}

func main() {
	f := squares()
	fmt.Println(f())    //  1
	fmt.Println(f())    //  4
	fmt.Println(f())    //  9
	fmt.Println(f())    // 16
}
```     

没有错，就是闭包技术。     

### 5.6.1 警告：捕获迭代变量

```go
var rmdirs []func()

for _, d := range tempDirs() {
	dir := d
	os.MkdirAll(dir, 0755)
	rmdirs = append(rmdirs, func() {
		os.RemoveAll(dir)
	})
}

for _, rmdir := range rmdirs {
	rmdir()
}
```    

你可能会感到困惑，为什么要在循环体中用循环变量 d 赋值一个新的局部变量，而不是像下面的代码
一样直接使用循环变量 dir。需要注意，下面的代码是错误的。     

```go
var rmdirs []func()

for _, dir := range tempDirs() {
  os.MkdirAll(dir, 0755)
  rmdirs = append(rmdirs, func() {
    os.RemoveAll(dir)
  })
}
```     

问题的原因在于循环变量的作用域。在上面的程序中，for 循环语句引入了新的词法块，循环变量
dir 在这个词法块中被声明。在该循环中生成的所有函数值都共享相同的循环变量。需要注意，函数值
中记录的是循环变量的内存地址，而不是循环变量某一时刻的值。以 dir 为例，后续的迭代会不断
更新 dir 的值，当删除操作执行时，for 循环已完成，dir 中存储的值等于最后一次迭代的值。
这意味着，每次对 os.RemoveAll 的调用删除的都是相同的目录。     

## 5.7 可变参数

参数数量可变的函数称之为可变参数函数。在声明可变参数函数时，需要在参数列表的最后一个参数
类型之前加上省略符号...，这表示该函数会接收任意数量的该类型参数。     

```go
func sum(vals ...int) int {
	total := 0
	for _, val := range vals {
		total += val
	}
	return total
}

func main() {
	fmt.Println(sum(1, 2, 3, 4))    //  10
	fmt.Println(sum(5, 6, 7))       // 18
	fmt.Println(sum())              // 0
}
```     

在函数体中，vals 被看作是类型为 []int 的切片。    

在上面的代码中，调用者隐式的创建一个数组，并将原始参数复制到数组中，再把数组的一个切片作为
参数传给被调函数。如果原始参数已经是切片类型，只需要在最后一个参数后加上省略符。     

```go
values := []int{1,2,3,4}
fmt.Println(sum(values...))
```     

虽然在可变参数函数内部，...int 型参数的行为看起来很像切片类型，但实际上，可变参数函数和
以切片作为参数的函数是不同的。    

```go
func f(...int) {}
func g([]int) {}
fmt.Printf("%T\n", f)   // func(...int)
fmt.Printf("%T\n", g)   // func([]int)
```    

## 5.8 Deferred 函数

你只需要在调用普通函数或方法前加上关键字 defer，就完成了 defer 所需要的语法。当defer 
语句被执行时，跟在 defer 后面的函数会被延迟执行。直到包含该 defer 语句的函数执行完毕时，
defer 后的函数才会被执行，无论包含 defer 语句的函数是通过 return 正常结束，还是由于 panic
导致的异常结束。你可以在一个函数中执行多条 defer 语句，它们的执行顺序与声明顺序相反。   

defer 语句经常被用于处理成对的操作，如打开、关闭、连接、断开连接、加锁、释放锁。通过
defer 机制，不论函数逻辑多复杂，都能保证在任何执行路径下，资源被释放。释放资源的 defer
应该直接跟在请求资源的语句后。     

defer语句中的函数会在return语句更新返回值变量后再执行，又因为在函数中定义
的匿名函数可以访问该函数包括返回值变量在内的所有变量，所以，对匿名函数采用defer机
制，可以使其观察函数的返回值。      

被延迟执行的匿名函数甚至可以修改函数返回给调用者的返回值：   

```go
func triple(x int) (result int) {
  defer func() { result += x}
  return double(x)
}
fmt.Println(triple(4))   // 12
```     

## 5.9 Panic 异常

Go 的类型系统会在编译时捕获很多错误，但有些错误只能在运行时检查，如数组访问越界、空指针
引用等。这些运行时错误会引起 panic 异常。    

一般而言，当 panic 异常发生时，程序会中断运行，并立即执行在该 goroutine(可以暂时理解为
线程)中被延迟的函数（defer机制）。随后，程序崩溃并输出日志信息。日志信息包括 panic value
和函数调用的堆栈跟踪信息。panic value通常是某种错误信息。      

不是所有的 panic 异常都来自运行时，直接调用内置的 panic 函数也会引发 panic 异常；panic
函数接受任何值作为参数。当某些不应该发生的场景发生时，我们就应该调用 panic。    

虽然 Go 的 panic 机制类似于其他语言的异常，但 panic 的使用场景有一些不同。由于 panic 会
引起程序的崩溃，因此 panic 一般用于严重错误，如程序内部的逻辑不一致。     

## 5.10 Recover 捕获异常

通常来说，不应该对 panic 异常做任何处理，但有时，也许我们可以从异常中恢复，至少我们可以
在程序崩溃前，做一些操作。     

如果在 deferred 函数中调用了内置函数 recover，并且定义该 defer 语句的函数发生了 panic
异常，recover 会使程序从panic中恢复，并返回 panic value。导致 panic 异常的函数不会
继续运行，但能正常返回，在未发生 panic 时调用 recover，recover 会返回 nil。     

