# 第二章 程序结构

<!-- TOC -->

- [第二章 程序结构](#第二章-程序结构)
  - [2.1 命名](#21-命名)
  - [2.2 声明](#22-声明)
  - [2.3 变量](#23-变量)
    - [2.3.1 简短变量声明](#231-简短变量声明)
    - [2.3.2 指针](#232-指针)
    - [2.3.3 new 函数](#233-new-函数)
    - [2.3.4 变量的生命周期](#234-变量的生命周期)
  - [2.4  赋值](#24--赋值)
    - [2.4.1 元祖赋值](#241-元祖赋值)
    - [2.4.2 可赋值性](#242-可赋值性)
  - [2.5 类型](#25-类型)
  - [2.6 包和文件](#26-包和文件)
    - [2.6.1 导入包](#261-导入包)
    - [2.6.2 包的初始化](#262-包的初始化)
  - [2.7 作用域](#27-作用域)
  - [2.8 总结](#28-总结)

<!-- /TOC -->

## 2.1 命名

Go 语言中的函数名、变量名、常量名、类型名、语句标号（就是可以控制break和continue跳出点
的那个 label 名吧）和包名都遵循一个简单的规则：名字必须以字母或下划线开头，后面可以跟
任意数量的字母、数字和下划线。     

Go 语言中类似 if 和 switch 的关键字有 25个，这 25个关键字不能用作命名，只能在特定的语法结构中使用：    

```
break    default       func    interface     select
case     defer         go      map           struct
chan     else          goto    package       switch
const    fallthrough   if      range         type
continue for           import  return        var
```     

此外还有大约30多个预定义的名字，比如 int 和 true 等，主要对应于内置的常量、类型和函数：    

```
常量：   true  false  iota   nil
类型：   int   int8   int16   int32   int64
        uint uint8   uint16   uint32   uint64   uintptr
        float32   float64   complex128    complex64
        bool    byte   rune   string   error
函数：   make   len   cap  new   append   copy  close  delete
        complex   real   imag   panic   recover
```     

这些内部预先定义的名字并不是关键字，你可以在定义中重新使用它们。在一些特殊的场景中重新
定义它们是有意义的，但是也要注意避免过度而引起语义混乱。    

如果一个实体是在函数内部定义，那么它就是一个函数局部的实体。如果定义在函数外，它在
所属包中的所有文件中都是可见的。名字的开头字母的大小写决定了名字在包外的可见性。如果一个
名字是大写字母开头的，那么它将是导出的，意味着可以被外部的包访问。     

在习惯上，Go 语言建议使用驼峰式命名。      

## 2.2 声明

声明语句定义了程序的各种实体对象，以及部分或全部的属性。Go 语言主要有四种类型的声明语句：`var, const, type, func`。     

在包一级声明语句声明的名字可在整个包对应的每个源文件中访问，而不仅仅在其声明语句所在的
源文件中访问。相比之下，局部声明的名字就只能在函数内部很小的范围被访问。    

OK，现在出现了几个问题，第一，变量重名怎么办，不允许重名吗？第二个，函数是否可以在函数内声明。     

## 2.3 变量

变量声明的语法一般如下：    

`var variableName type = expression`     

其实 type 和 = expression 两个部分可以省略其中的一个。数值类型变量对应的零值是0，
布尔类型变量对应的零值是false，字符串类型对应的零值是空字符串，接口或引用类型（包括 slice, map, chan 和函数）变量对应的零值是 nil。
数组或结构体等聚合类型对应的零值是每个元素或字段都是对应该类型的零值。     

可以在一个声明语句中同时声明一组变量，或用一组初始化表达式声明并初始化一组变量。如果省略每个变量的类型，将可以声明多个类型不同的变量：    

```go
var i, j, k int
var b, f, s = true, 2.3, "four"
```    

初始化表达式可以是字面量或任意的表达式。在包级别声明的变量会在 main 入口函数执行前完成初始化，局部变量将在声明语句被执行到的时候完成初始化。     

一组变量也可以通过调用一个函数，由函数返回的多个返回值初始化：    

`var f, err = os.Open(name)`      

### 2.3.1 简短变量声明

和 var 形式声明语句一样，简短变量声明语句也可以用来声明和初始化一组变量：    

`i, j := 0, 1`     

这里有一个比较微妙的地方：简短变量声明左边的变量可能并不是全部都是刚刚声明的。如果有一些
已经在相同的词法域声明过了，那么简短变量声明语句对这些已经声明过的变量就只有赋值行为了。    

简短变量声明语句中必须至少要声明一个新的变量，下面的代码将不能编译通过：    

```go
f, err := os.Open(infile)
// ...
f, err := os.Create(outfile)
```     

解决的方法时第二个简短变量声明语句改用普通的多重赋值语言。    

所以要分清带 var 不带 var，简短变量声明。      

简短变量声明语句只有对已经在同级词法域声明过的变量才和赋值操作语句等价，如果变量是在
外部词法域声明的，那么简短变量声明语句将会在当前词法域重新声明一个新的变量。     

### 2.3.2 指针

一个变量对应一个保存了变量对应类型值的内存空间。普通变量在声明语句创建时被绑定到一个
变量名，比如叫 x 的变量，但是还有很多变量始终以表达式方式引入，例如 x[i] 或 x.f 变量。
所有这些表达式一般都是读取一个变量的值，除非它们是出现在赋值语句的左边，这种时候是给对应变量赋予一个新的值。    

一个指针的值是另一个变量的地址。一个指针对应变量在内存中的存储位置。并不是每一个值都会
有一个内存地址，但是对于每一个变量必然有对应的内存地址。通过指针，我们可以直接读或更新对应变量的值，而不需要知道该变量的名字。    

如果用 `var x int` 声明语句声明一个 x 变量，那么 `&x` 表达式（取x变量的内存地址）将
产生一个指向该整数变量的指针，指针对应的数据类型是 `*int`，指针被称之为“指向 int 类型的指针”，
如果指针名字为 p，那么可以说“p 指针指向变量 x”，或者说“p 指针保存了 x 变量的内存地址”。
同时 `*p` 表达式对应 p 指针指向的变量的值。一般 `*p`表达式读取指针指向变量的值。这里
为 int 类型的值。同时因为 `*p` 对应一个变量，所以该表达式也可以出现在赋值语句的坐标，表示更新指针所指向的变量的值。       

```go
x := 1
p := &x
fmt.Println(*p)  // 1
*p = 2
fmt.Println(x)   // 2
```     

也就说我们可以这样理解，假设现在程序里有一张表对应了变量名和变量值的内存地址，那么普通变量
值的内存地址所指向的内存中存放的是变量值，而指针变量存放的是另一个变量的变量值的内存地址。     

对于聚合类型的每个成员——比如结构体的每个字段、或者是数组的每个元素--也都是对应一个变量，
因此可以被取地址。     

变量有时候被称为可寻址的值。即使变量由表达式临时生成，那么表达式也必须能接收 `&` 取地址
操作。     

任何类型的指针的零值都是 nil。如果 `p != nil` 测试为真，那么 p 是指向某个有效变量。指针
之间也是可以进行相等测试的，只有当它们指向同一个变量或全部是 nil 时才相等。     

在 Go 语言中，返回函数中局部变量的地址也是安全的。例如下面的代码，调用 f 函数时创建局部
变量v，在局部变量地址被返回之后依然有效，因为指针 p 依然引用这个变量。    

```go
var p = f()

func f() *int {
  v := 1
  return &v
}
```    

每次调用 f 函数都将返回不同的结果。    

因为指针包含了一个变量的地址，因此如果将指针作为参数调用函数，那将可以在函数中通过该指针
来更新变量的值。      

每次我们对一个变量取地址，或者复制指针，我们都是为原变量创建了新的别名。例如，`*p` 就是
变量 v 的别名。      

### 2.3.3 new 函数

另一个创建变量的方法是调用内置的 new 函数。表达式 `new(T)` 将创建一个 `T` 类型的匿名
变量，初始化为 T 类型的零值，然后返回变量地址，返回的指针类型为 `*T`.       

```go
p := new(int)
fmt.Println(*p)
*p = 2
fmt.Println(*p)
```    

用 new 创建变量和普通变量声明语句方式创建变量没有什么区别，除了不需要声明一个临时变量的
名字外，我们还可以在表达式中使用 `new(T)`。换言之，new 函数类似是一种语法糖，而不是一个
新的基础概念。     

由于 new 只是一个预定义的函数，它并不是一个关键字，因此我们可以将 new 名字重新定义为别的
类型。     

### 2.3.4 变量的生命周期

变量的生命周期指的是在程序运行期间变量有效存在的时间间隔。对于在包一级声明的变量来说，它们的
生命周期和整个程序的运行周期是一致的。而相比之下，局部变量的生命周期则是动态的：从每次
创建一个新变量的声明语句开始，直到该变量不再被引用为止，然后变量的存储空间可能被回收。
函数的参数变量和返回值变量都是局部变量。它们在函数每次被调用的时候创建。     

那么 Go 语言的自动垃圾收集器是如何直到一个变量是何时可以被回收的呢？这里我们可以避开完整
的技术细节，基本的实现思路是，从每个包级的变量和每个当前运行函数的每一个局部变量开始，通过
指针或引用的的访问路径遍历，是否可以找到该变量。如果不存在这样的访问路径，那么说明该变量
是不可达的，也就是说它是否存在并不会影响程序后续的计算结果。     

因为一个变量的有效周期只取决于是否可达，因此一个循环迭代内部的局部变量的生命周期可能超出
其局部作用域。同时，局部变量可能在函数返回之后依然存在。   

编译器会自动选择在栈上还是堆上分配局部变量的存储空间，但可能令人惊讶的是，这个选择并不是由
用 var 还是 new 声明变量的方式决定的。    

```go
var global *int

func f() {
  var x int
  x = 1
  global = &x
}

func g() {
  y := new(int)
  *y = 1
}
```    

f 函数里的 x 变量必须在堆上分配，因为它在函数退出后依然可以通过包一级的 global 变量找到，
虽然它是在函数内部定义的，用 Go 语言的术语说，这个 x 局部变量从函数 f 中逃逸了。相反，
`*y`并没有从函数 g 中逃逸，编译器可以选择在栈上分配 `*y` 的存储空间。      

## 2.4  赋值

### 2.4.1 元祖赋值

元祖赋值是一种形式的赋值语句，它允许同时更新多个变量的值。在赋值之前。赋值语句右边的所有
表达式将会先进行求值，然后再统一更新左边对应变量的值。这对于处理有些同时出现在元祖赋值语句
左右两边的变量很有帮助，例如我们可以这样交换两个变量的值：    

```go
x, y = y, x
a[i], a[j] = a[j], a[i]
```    

### 2.4.2 可赋值性

赋值语句是显式的赋值形式，但是程序中还有很多地方会发生隐式的赋值行为：函数调用会隐式地将
调用参数的值赋值给函数的参数变量，一个返回语句隐式地将返回操作的值赋值给结果变量。     

不管是隐式还是显式地赋值，在赋值语句左边的变量和右边最终的求到的值必须有相同的数据类型。
更直白地说，只有右边的值对于左边的变量是可赋值的，赋值语句才是允许的。     

可赋值性的规则对于不同类型有着不同要求，对每个新类型特殊的地方我们会专门解释。对于目前我们
已经讨论过的类型，它的规则是简单的：类型必须完全匹配，nil 可以赋值给任何指针或引用类型的变量。
常量则有更灵活的赋值规则，因为这样可以避免不必要的显式的类型转换。      

对于两个值是否可以用 ==  或 !=  进行相等比较的能力也和可赋值能力有关系：对于任何类型
的值的相等比较，第二个值必须是对第一个值类型对应的变量是可赋值的，反之依然。     

## 2.5 类型

变量或表达式的类型定义了对应存储值的属性特征，例如数值在内存的存储大小，它们在内部是如何
表达的，是否支持一些操作符，以及它们自己关联的方法集等。    

一个类型声明语句创建了一个新的类型名称，和现有类型具有相同的底层结构。新命名的类型提供了
一个方法，用来分隔不同概念的类型，这样即使它们底层类型相同也是不兼容的。     

`type 类型名字 底层类型`     

类型声明语句一般出现在包一级，因此如果新创建的类型名字的首字符大写，则在外部包也可以使用。     

```go
package tempconv

type Celsius float64
type Fahrenheit float64

const (
	AbsoluteZeroC Celsius = -273.15
	FreezingC Celsius = 0
	BoilingC Celsius = 100
)

func CToF (c Celsius) Fahrenheit {
	return Fahrenheit(c * 9 / 5 + 32)
}

func FToC (f Fahrenheit) Celsius {
	return Celsius((f - 32) * 5 / 9)
}
```     

我们在这个包声明了两种类型：Celsius 和 Fahrenheit 分别对应不同的温度单位。它们虽然有着
相同的底层类型 float64，但是它们是不同的数据类型，因此它们不可以被相互比较或混在一个表达
式运算。需要一个类似 `Celsius(t)` 或 `Fahrenheit(t)` 形式的显式转型操作才能将 `float64`
转为对应的类型。`Celsius(t)` 或 `Fahrenheit(t)` 是类型转换操作，它们并不是函数调用。
类型转换并不会改变值本身，但是会使它们的语义发生变化。     

对于每一个类型 T，都有一个对应的类型转换操作 `T(x)`，用于将 x 转为 T 类型。只有当两个
类型的底层基础类型相同时，才允许这种转型操作，或者是两者都是指向相同底层结构的指针类型，
这些转换只改变类型而不会影响值本身。     

数值类型之间的转换也是允许的，并且在字符串和一些特定类型的 slice 之间也是可以转换的。
这类转换可能改变值的表现，例如将一个浮点数转为整数将丢弃小数部分。在任何情况下，运行时不会
发生转换失败的错误。      

比较运算符例如 `==` 和 `<` 可以用来比较一个命名类型的值和另一个相同类型的值，或者一个
命名类型的值和一个其底层类型的值，但是两个值如果有着不同命名类型，也就是说，要么比较的
两者类型相同，要么一边是另一边的底层类型值：    

```go
var c Celsius
var f Fahrenheit

fmt.Println(c == 0)  // "true"
fmt.Println(f >= 0)  // "true"
fmt.Println(c == f)  // compile error: type mismatch
fmt.Println(c == Celsius(f))
```     

**问题**：Go 中到底包含有多少运算符，看样子是否支持算术运算取决与类型的底层类型，但是比较
运算符的话就要求两边是同一类型，或者一边是底层类型。    

命名类型还可以为该类型的值定义新的行为。这些行为表示为一组关联到该类型的函数集合，我们
称为类型的方法集。我们会在第6章讨论细节，不过下面可以先简单介绍一下：    

下面的声明语句，Celsius 类型的参数 c 出现在了函数名的前面，表示声明的是Celsius 类型的
一个名叫 `String` 的方法，该方法返回该类型对象 c 带着 `°C` 温度单位的字符串：     

```go
func (c Celsius) String() string {
  return fmt.Sprintf("%g°C", c)
}
```     

许多类型都会定义一个 `String` 方法，因为当使用 fmt 包的打印方法时，将会优先使用该类型
对应的 `String` 方法返回的结果打印。      

```go
c := FToC(212.0)          
fmt.Println(c.String())   // 100°C
fmt.Printf("%v\n", c)     // 100°C
fmt.Printf("%s\n", c)     // 100°C
fmt.Println(c)            // 100°C
fmt.Printf("%g\n", c)     // 100
fmt.Println(float64(c))   // 100
```     

## 2.6 包和文件

通常一个包所在的目录路径的后缀是包的导入路径。    

每个包都对应一个独立的命名空间。    

包还可以让我们通过控制哪些名字是外部可见的来隐藏内部实现信息。在 Go 语言中，一个简单的
规则是：如果一个名字是大写字母开头的，那么该名字是可导出的。      

包级别的名字，例如在一个文件声明的类型和常量，在同一个包的其他源文件也是可以直接访问的，
就好像所有代码都在一个文件一样。     

在每个源文件的包声明前仅跟着的注释是包注释。通常，包注释的第一句应该显示包的功能概要。
一个包通常只有一个源文件有包注释。如果包注释很大，通常会放到一个独立的 doc.go 文件中。    

### 2.6.1 导入包

在 Go 语言程序中，每个包都有一个全局唯一的导入路径。除了包的导入路径，每个包还有一个包名。
按照惯例，一个包的名字和包的导入路径的最后一个字段相同。     

导入本地包好像必须是相对路径，且是相对于当前目录的。     

### 2.6.2 包的初始化

包的初始化是按照包级变量的声明顺序初始化变量，不过首先要解决了依赖关系：    

```go
var a = b + c       // a 第三个初始化，值为 3
var b = f()         // b 第二个初始化，值为 2
var c = 1           // c 第一个初始化，值为 1

func f() { return c + 1 }
```   

如果包中含有多个 .go 源文件，它们将按照发给编译器的顺序进行初始化，Go 语言的构建工具首先
会将 .go 文件根据文件名排序，然后依次调用编译器编译。     

对于在包级别声明的变量，如果有初始化表达式则用表达式初始化，还有一些没有初始化表达式的，
例如某些表格数据初始化并不是一个简单的赋值过程。在这种情况下，我们可以用一个特殊的 init
初始化函数来简化初始化工作。每个文件都可以包含多个 `init` 初始化函数。     

`func init() { /*....*/ }`     

这样的 init 初始化函数除了不能被调用或引用外，其他行为和普通函数类似。在每个文件中的 init
初始化函数，在程序开始执行时按照它们声明的顺序被自动调用。     

每个包在解决依赖的前提下，以导入声明的顺序初始化，每个包只会被初始化一次。因此，如果一个
p 包带入了 q 包，那么在 p 包初始化的时候可以认为 q 包必然已经初始化过了。初始化工作是
自下而上进行的，main 包最后被初始化。以这种方式，可以确保在 main 函数执行之前，所有依赖
的包都已经完成初始化工作了。      


## 2.7 作用域

一个声明语句将程序中的实体和一个名字关联，比如一个函数或一个变量。声明语句的作用域是指
源代码中可以有效使用这个名字的范围。     

不要将作用域和生命周期混为一谈。声明语句的作用域对应的是一个源代码的文本区域；它是一个编译
时的属性。一个变量的生命周期是指程序运行时变量存在的有效时间段，在此时间区域内它可以被程序
的其他部分引用；是一个运行时的概念。     

语法块是由花括号所包含的一系列语句，就像函数体或循环体花括号对应的语法块那样。语法块内部
声明的名字是无法被外部语法块访问的。语法块决定了内部声明的名字的作用域范围。有一个语法块
为整个源代码，称为全局语法块；然后是每个包的包语法块；每个 for、if 和 switch 语句的语法
块；每个 switch 或 select 也有独立的语法块；当然也包括显式书写的语法块。        

对于导入的包，例如 tempconv 导入的 fmt 包，则是对应源文件级的作用域，因此只能在当前的
文件中访问导入的 fmt 包，当前包的其他源文件无法访问在当前源文件导入的包。     

控制流标号，就是 break, continue 或 goto 语句后面跟着的那种标号，则是函数级的作用域。     

```go
func main() {
	x := "hello!"
	for i := 0; i < len(x); i++ {
		x := x[i]
		if x != '!' {
			x := x + 'A' - 'a'
			fmt.Printf("%c", x) // "HELLO" (one letter per iteration)
		}
	}
}
```     

正如上面例子所示，並不是所有的词法域都显式地对应到由花括号包含的语句；还有一些隐含的规则。
上面的 for 语句创建了两个词法域：花括号包含的显式部分是 for 循环体部分词法域，另外一个
隐式的部分则是循环的初始化部分，比如用于迭代变量 i 的初始化。隐式的词法域部分的作用域还
包含条件测试部分和循环后的迭代部分，当然也包含循环体词法域。      

和 for 循环类似，if 和 switch 语句也会在条件部分创建隐式词法域，还有它们对应的执行体词法域：   

```go
if x := f(); x == 0 {
  fmt.Println(x)
} else if y := g(x); x == y {
  fmt.Println(x, y)
} else {
  fmt.Println(x, y)
}
fmt.Println(x, y)  // compile error: x and y are not visible here
```    

第二个 if 语句嵌套在第一个内部，因此第一个 if 语句条件初始化词法域声明的变量在第二个 if 
中也可以访问。switch 语句的每个分支也有类似的词法域规则：条件部分为一个隐式词法域，
然后每个是每个分支的词法域。     

在包级别，声明的顺序并不会影响作用域范围，因此一个先声明的可以引用它自身或者是引用后面的
一个声明。    


## 2.8 总结

命名规则：函数、变量、常量、类型、包、statement labels 的命名必须以字母或下划线开头，后
接若干字母、数字、字符串。包级命名开头字母的大小写决定了名字是否是可导出的。    


一个声明定义了一个程序实体以及其部分或全部的属性。Go 中包括4种类型的声明：变量、常量、
类型及函数。    

变量声明创建了一个特定类型的变量，并将其与一个名字绑定起来，设置它的初始值。每条声明都有
如下的形式：    

`var name type = expression`    

其中 `type` 或 ` = expression` 部分是可以省略的，但不能两部分同时省略。如果省略了
初始化表达式，初始值就是对应类型的零值，数值是0，字符串是""，布尔是 false，接口和引用
类型(slice, map, pointer, channel, function)是 nil，聚合类型数组和结构体是元素各个
类型的零值。     

包级变量会在 `main` 函数执行前就初始化。      

包级变量的生命周期是程序的整个执行期间。局部变量的生命周期则是动态的：每次声明语句执行的
时候都会创建一个新的实例，生命周期直到变量是不可达的。   

每个算术运算符和位操作运算符都有对应的赋值运算符版本。    

对于每一个类型 T，都有一个对应的类型转换操作 T(x)。只有当两个类型的底层基础类型相同时，
才允许这种转型操作，或者两者都是指向相同底层结构的指针类型。     

数值类型之间的转换是允许的，而且在字符串和一些特定类型的 slice 之间转换也是可行的。不过
这类转换可能改变值的表现，例如将一个字符串转为 `[]byte` 类型的 slice 将拷贝一个字符串
数据的副本。    

比较运算符 ==  和 <  也可以用来比较一个命名类型的变量和另一个有相同类型的变量，或有
着相同底层类型的未命名类型的值之间做比较。但是如果两个值有着不同的类型，则不能直
接进行比较。       






