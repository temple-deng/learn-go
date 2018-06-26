# 第四章 复合数据类型

<!-- TOC -->

- [第四章 复合数据类型](#第四章-复合数据类型)
  - [4.1 数组](#41-数组)
  - [4.2 slice](#42-slice)
    - [4.2.1 append 函数](#421-append-函数)
    - [4.2.2 slice 内存技巧](#422-slice-内存技巧)
  - [4.3 map](#43-map)
  - [4.4 结构体](#44-结构体)
    - [4.4.1 结构体字面值](#441-结构体字面值)
    - [4.4.2 结构体比较](#442-结构体比较)
    - [4.4.3 结构体嵌入和匿名成员](#443-结构体嵌入和匿名成员)
    - [4.4.4 数组、结构体、slice、map 比较](#444-数组结构体slicemap-比较)
  - [4.5 JSON](#45-json)
  - [4.6 文本和 HTML 模板](#46-文本和-html-模板)

<!-- /TOC -->

注意标题是复合数据类型 composite types，包括数组、slice、map和结构体。数组和结构体是
聚合类型，而slice 和map 是引用类型。区别是数组和结构体都是有固定内存大小的数据结构。
相比之下，slice 和 map 则是动态的数据结构，它们将根据需要动态增长。    

## 4.1 数组

数组是一个固定长度的由零个或多个特定类型元素组成的序列。由于其长度是固定的，因此在 Go
中很少直接使用数组，相反的，slice 的长度是可变的，用的更多。     

数组的每个元素都可以通过索引下标来访问，内置 len 函数返回数组中元素的个数：    

```go
var a [3]int
fmt.Println(a[0])
fmt.Println(a[len(a)-1])

for i, v := range a {
  fmt.Printf("%d %d\n", i, v)
}
```     

默认情况下,数组的每个元素都被初始化为元素类型对应的零值,对于数字类型来说就是0。
我们也可以使用数组字面值语法用一组值来初始化数组:    

```go
var q [3]int = [3]int{1,2,3}
var r [3]int = [3]int{1,2}
fmt.Println(r[2])   //0
```     

在数组字面值中,如果在数组的长度位置出现的是“...”省略号,则表示数组的长度是根据初始
化值的个数来计算。因此,上面q数组的定义可以简化为:    

```go
q := [...]int{1,2,3}
```    

注意省略号的语法只能用在数组字面值写法里面，不能用来类型定义的部分。但是如果类型
定义即声明了个数，右边也还是可以用省略号的语法，不过这时候两边的长度必须相等。     

数组的长度是数组类型的一个组成部分,因此[3]int和[4]int是两种不同的数组类型。数组的长
度必须是常量表达式,因为数组的长度需要在编译阶段确定。     

```go
q := [3]int{1,2,3}
q = [4]int{1,2,3,4}     // compile error: cannot assign [4]int to [3]int
```     

数组、slice、map和结构体字面值的写法都很相似。之前好像在第一章提到过，这种叫复合字面量
composite literals，上面的形式是直接提供顺序初始化值序列，但是也可以指定一个索引和对应值列表的方式初始化，就像下面这样：    

```go
type Currency int

const (
  USD Currency = iota
  EUR
  GBP
  RMB
)

symbol := [...]string{USD: "$", EUR: "€", GBP: "£", RMB: "¥"}

fmt.Println(RMB, symbol[RMB])
```    

在这种形式的数组字面值形式中，初始化索引的顺序是无关紧要的，而且没用到的索引可以
省略。    

如果一个数组的元素类型是可以相互比较的,那么数组类型也是可以相互比较的,这时候我
们可以直接通过==比较运算符来比较两个数组,只有当两个数组的所有元素都是相等的时候
数组才是相等的。不相等比较运算符!=遵循同样的规则。那岂不是可以比较的数组可以用作
map 的键名（经测试的确是可以的）。      

当调用一个函数的时候,函数的每个调用参数将会被赋值给函数内部的参数变量,所以函数
参数变量接收的是一个复制的副本,并不是原始调用的变量。因为函数参数传递的机制导致
传递大的数组类型将是低效的,并且对数组参数的任何的修改都是发生在复制的数组上,并
不能直接修改调用时原始的数组变量。在这个方面,Go语言对待数组的方式和其它很多编程
语言不同,其它编程语言可能会隐式地将数组作为引用或指针对象传入被调用的函数。     

## 4.2 slice

slice 代表变长的序列，序列中每个元素都有相同的类型。一个 slice 类型一般写作 []T，其中
T 代表 slice 中元素的类型。     

数组和 slice 之间有着紧密的联系。一个 slice 是一个轻量级的数据结构，提供了访问
数组子序列（或者全部）元素的功能，而且 slice 的底层确实引用了一个数组对象。一个slice
由三个部分构成：指针、长度和容量。指针指向第一个 slice 元素对应的底层数组元素的地址，
要注意的是 slice 的第一个元素并不不一定是数组的第一个元素。长度对应 slice 中元素
的数目；长度不能超过容量，容量一般是从 slicee 的开始位置到底层数据的结尾位置。内置
的len 和 cap 函数分别返回 slice 的长度和容量。    

多个 slice 之间可以共享底层的数据，并且引用的数组部分区间可能重叠：   

```go
months := [...]string{
  1: "January",
  2: "February",
  3: "March",
  4: "April",
  5: "May",
  6: "June",
  7: "July",
  8: "August",
  9: "September",
  10: "October",
  11: "November",
  12: "December"
}
```   

slice 的切片操作 `s[i:j]`，其中 0 &lt;= i &lt;= j &lt;= cap(s)，用于创建
一个新的 slice，j 位置的索引被省略就使用 `len(s)`:    
```go
Q2 := months[4:7]
summer := months[6:9]
```    

如果切片操作超出 `cap(s)` 的上限将导致一个 panic 异常，但是超出 `len(s)` 则是意味着
扩展了 slice，因为新 slice 的长度会变大。    

```go
fmt.Println(summer[:20])    // panic: out of range

endlessSummer := summer[:5]
fmt.Println(endlessSummer)     // [June July August September October]
```     

另外，字符串的切片操作和[]byte 字节类型切片的切片操作是类似的。它们都写作x[m:n]，并且
都是返回一个原始字节系列的子序列，底层都是共享之前的底层数组，因此切片操作对应常量时间
复杂度。x[m:n]切片操作对于字符串则生成一个新字符串，如果x是[]byte的话则生成
一个新的[]byte。      

因为 slice 值包含指向第一个 slice 元素的指针，因此向函数传递 slice 将允许在函数
内部修改底层数组的元素。换句话说，复制一个 slice 只是对底层的数组创建了一个新的 slice
别名。    

和数组不同的是,slice之间不能比较,因此我们不能使用==操作符来判断两个slice是否含有
全部相等元素。slice唯一合法的比较操作是和nil比较。   

一个零值的slice等于nil。一个nil值的slice并没有底层数组。一个nil值的slice的长度和容量都
是0，但是也有非nil值的slice的长度和容量也是0的，例如[]int{}或make([]int, 3)[3:]。与任意类
型的nil值一样，我们可以用[]int(nil)类型转换表达式来生成一个对应类型slice的nil值。  

```go
var s []int    // len(s) == 0, s == nil
s = nil        // len(s) == 0, s == nil
s = []int(nil) // len(s) == 0, s == nil
s = []int{}    // len(s) == 0, s != nil
```   

内置的 `make` 函数创建一个指定元素类型、长度和容量的 slice。容量部分可以省略，在
这种情况下，容量将等于长度。    

```go
make([]T, len)
make([]T, len, cap)
```     

在底层，`make` 创建了一个匿名的数组变量，然后返回一个 slice。       

来看一下我们目前知道创建 slice 的方法：   

+ 通过字面量直接声明一个 slice
+ 通过 `make` 函数
+ 通过切割 slicing 一个数组或 slice    

### 4.2.1 append 函数

内置的append 函数用于向 slice 追加元素：    

```go
var runes []rune
for _, r := range "Hello, 世界" {
  runes = append(runes, r)
}

fmt.Printf("%q\n", runes)  // ['H' 'e' 'l' 'l' 'o' ',' ' ' '世' '界']
```     

当然对应这个特殊的问题我们可以通过Go语言内置的[]rune("Hello, 世界")转换操作完成。    

```go
func appendInt(x []int, y int) []int {
	var z []int
	zlen := len(x) + 1
	if zlen <= cap(x) {
		z = x[:zlen]
	} else {
		zcap := zlen

		if zcap < zlen * 2 {
			zcap = zlen * 2
		}

		z = make([]int, zlen, zcap)
		copy(z, x)
	}
	z[len(x)] = y
	return z
}
```    

内置的copy函数可以方便地将一个slice复制另一个相同类型的slice。copy函数的第一个参数
是要复制的目标slice,第二个参数是源slice,目标和源的位置顺序和dst	=	src赋值语句是一致的。
copy 函数将返回成功复制的元素的个数，等于两个 slice 中较小的长度。      

内置的append函数可能使用比appendInt更复杂的内存扩展策略。因此,通常我们并不知道
append调用是否导致了内存的重新分配,因此我们也不能确认新的slice和原始的slice是否引
用的是相同的底层数组空间。同样,我们不能确认在原先的slice上的操作是否会影响到新的
slice。因此,通常是将append返回的结果直接赋值给输入的slice变量:    

`runes	=	append(runes,	r)`      

内置的append函数可以追加多个元素,甚至追加一个slice：    

```go
var x []int
x = append(x, 1)
x = append(x, 2, 3)
x = append(x, 4, 5, 6)
x = append(x, x...)
```      

需要注意的是，当 append 操作后的 slice 的容量没有超过底层数组的容量时，append 操作是会修改
底层数据的：     

```go
primes := [6]int{2, 3, 5, 7, 11, 13}

var s []int = primes[1:4]
s = append(s, 8, 9)
fmt.Println(s)        // [3, 5, 7, 8, 9]
fmt.Println(primes)   // [2, 3, 5, 7, 8, 9]
```    

当时当 append 操作导致内存的重新分配，那么其操作就不再影响之前的数组了：　　　　

```go
primes := [6]int{2, 3, 5, 7, 11, 13}

var s []int = primes[1:4]
s = append(s, 8, 9, 10)
fmt.Println(s)         // [3, 5, 7, 8, 9, 10]
fmt.Println(primes)    // [2, 3, 5, 7, 11, 13]
```   

当重新分配内存后，slice 的容量是不确定的，并不会只是简单的与长度相等，通常可能会更长一些，以免
下次进行 append 操作时又要分配内存操作。    


### 4.2.2 slice 内存技巧

给定一个字符串列表，下面的 nonempty 函数将在原有 slice 内存空间之上返回不包含
空字符串的列表：     

```go
func nonempty(strings []string) []string {
	i := 0
	for _, s := range strings {
		if s != "" {
			strings[i] = s
			i++
		}
	}
	return strings[:i]
}
```     

比较微妙的地方是，输入的 slice 和输出的 slice 共享一个底层数组，这可以避免分配
另一个数组，不过原来的数组将可能被覆盖：    

```go
data := []string{"one", "", "three"}

fmt.Printf("%q\n", nonempty(data))
fmt.Printf("%q\n", data)
```     

## 4.3 map

在 Go 语言中，一个 map 就是一个哈希表的引用，map 类型可以写为 map[k]v。map 中所有的
key 都有相同的类型，所有的 value 也有着相同的类型。key 必须是支持 == 比较运算符的数据
类型。     

内置的 make 函数可以创建一个 map:    

`ages := make(map[string]int)`

好像 slice 和 map 都可以通过 make 函数创建，但是数组不可以。    

使用复合字面值创建 map 的语法如下：    

```go
ages := map[string]int{
  "alice": 31,
  "charlie": 35
}
```    

map 中的元素通过 key 对应的下标语法访问：    

```go
ages["alice"] = 32
fmt.Println(ages["alice"])
```     

使用内置的 delete 函数可以删除元素：    

`delete(ages, "alice")`      

所有这些操作是安全的，即使这些元素不在map中也没有关系；如果一个查找失败将返回
value类型对应的零值。但是map中的元素并不是一个变量，因此我们不能对map的元素进行取址操作：    

`_ = &ages["alice"]  // compile error: cannot take address of map element`     

禁止对map元素取址的原因是map可能随着元素数量的增长而重新分配更大的内存空间，从而
可能导致之前的地址无效。    

同理，估计对 slice 元素也不能取地址，因为 slice 也可能由于元素个数的增加而重新分配内存，
但是结构体和数组的元素就可以取地址，因为他们的内存大小都一定的。    

map类型的零值是nil，也就是没有引用任何哈希表。    

```go
var ages map[string]int
fmt.Println(ages == nil)    // true
fmt.Println(len(ages) == 0) // true
```     

map 上的大部分操作，包括查找、删除、len 和 range 循环都可以安全工作在 nil 值的map 上，
它们的行为和一个空的 map 类型。但是向一个 nil 值的 map 存入元素将导致一个 panic 异常。    

这一点和 slice 有区别，如果 slice 是 nil，是无法访问其中的类型，但是可以使用 range 循环。    

在向 map 存数据前必须先创建 map。     

通过key作为索引下标来访问map将产生一个value。如果key在map中是存在的，那么将得到
与key对应的value；如果key不存在，那么将得到value对应类型的零值，正如我们前面看到的
ages["bob"]那样。这个规则很实用，但是有时候可能需要知道对应的元素是否真的是在map
之中。例如，如果元素类型是一个数字，你可以需要区分一个已经存在的0，和不存在而返回
零值的0，可以像下面这样测试：     

```go
age, ok := ages["bob"]
if !ok {

}
```     

在这种场景下，map 的下标语法将产生两个值；第二个是一个布尔值，用于报告元素是否真的存在。     

和 slice 一样，map 之间也不能进行相等比较；唯一的例外是和 nil 进行比较。      

## 4.4 结构体

```go
type Employee struct {
  ID         int
  Name       string
  Address    string
  DoB        time.Time
  Position   string
  Salary     int
  ManagerID  int
}

var dilbert Employee
```    

结构体变量的成员可以通过点操作符访问，也可以直接对每个成员赋值：    

`dilbert.Salary -= 5000`      

或者是对成员取地址，然后通过指针访问：    

```go
position := &dilbert.Position
*position = "Senior " + *position
```    

点操作符也可以和指向结构体的指针一起工作：    

```go
var employeeOfTheMonth *Employee = &dilbert
employeeOfTheMonth.Position += " (proactive team player)"
```    

相当于下面的语句：    

`(*employeeOfTheMonth).Position += " (proactive team player)"`      


上面的意思是这样的：如果我们有一个指向结构体的指针 p，那么可以通过 (*p).X 来访问其字段 X。
不过这么写太啰嗦了，所以语言也允许我们直接写 p.X 就可以，而不需要明确的解引用。    

下面的EmployeeByID函数将根据给定的员工ID返回对应的员工信息结构体的指针。我们可以
使用点操作符来访问它里面的成员：    

```go
func EmployeeByID(id int) *Employee { /* ... */ }
fmt.Println(EmployeeByID(dilbert.ManagerID).Position) // "Pointy-haired boss"
id := dilbert.ID
EmployeeByID(id).Salary = 0 // fired for... no real reason
```     

后面的语句通过EmployeeByID返回的结构体指针更新了Employee结构体的成员。如果将
EmployeeByID函数的返回值从 *Employee  指针类型改为Employee值类型，那么更新语句将
不能编译通过，因为在赋值语句的左边并不确定是一个变量（译注：调用函数返回的是值，
并不是一个可取地址的变量）。不是很理解这里的操作了。    

通常一行对应一个结构体成员，成员的名字在前类型在后，不过如果相邻的成员类型相同的话可以被
合并到一行：    

```go
type Employee struct {
  ID       int
  Name, Address string
  DoB      time.Time
  Position  string
  Salary    int
  ManagerID  int
}
```    

结构体成员的输入顺序也有重要的意义。我们也可以将Position成员合并（因为也是字符串类
型），或者是交换Name和Address出现的先后顺序，那样的话就是定义了不同的结构体类
型。纳尼？     

如果结构体成员名字是以大写字母开头的，那么该成员就是导出的；这是Go语言导出规则决
定的。一个结构体可能同时包含导出和未导出的成员。     

**问题**：结构体的成员可导出的意义是什么？    

一个命名为S的结构体类型将不能再包含S类型的成员：因为一个聚合的值不能包含它自身。
（该限制同样适应于数组。）但是S类型的结构体可以包含 *S  指针类型的成员，这可以让我
们创建递归的数据结构，比如链表和树结构等。     

### 4.4.1 结构体字面值

结构体也可以使用字面值表示：    

```go
type Point struct {
  X, Y int
}

p := Point{1, 2}
```     

有两种形式的结构体字面值语法，上面的是第一种写法，要求以结构体成员定义的顺序为每个结构体
成员指定一个值。这个的话要求我们记住字段的顺序，很明显是不合理的。    

更常用的是第二种写法，以成员的名字和值来初始化，可以包含部分或全部的成员。     

两种形式的写法不能混合使用。     

结构体可以作为函数的参数和返回值。    

如果要在函数内部修改结构体成员的话，用指针传入是必须的，因为在 Go 语言中，**所有的函数**
**参数都是值拷贝传入的，函数参数将不再是函数调用时的原始变量**      

### 4.4.2 结构体比较

如果结构体的全部成员是可以比较的，那么结构体也是可以比较，那么结构体也是可以用 == 或 !=
比较的。     

### 4.4.3 结构体嵌入和匿名成员

Go 语言有一个特性让我们只声明一个成员对应的数据类型而不知名成员的名字；这类成员就叫匿名
成员。匿名成员的数据类型必须是命名的类型或指向一个命名的类型的指针。下面的代码中，Circle
和 Wheel 各自都有一个匿名成员。我们可以说 Point 类型被嵌入到了 Circle 结构体，同时
Circle 类型被嵌入到了 Wheel 结构体。     

```go
type Circle struct {
  Point
  Radius int
}

type wheel struct {
  Circle
  Spokes int
}
```     

得意于匿名嵌入的特性，我们可以直接访问叶子属性而不需要给出完整的路径：    

```go
var w Wheel
w.X = 8         // 等价于 w.Circle.Point.X = 8
w.Y = 8
w.Radius = 5
w.Spokes = 20
```     

这个东西有点像 JS 中的继承属性。     

不幸的是，结构体字面值并没有简短表示匿名成员的语法，因此下面的语句都不能编译通过：    

```go
w = Wheel{8,8,5,20}
w = Wheel{X:8, Y:8, Radius:5, Spokes: 20}
```     


结构体字面值必须遵循形状类型声明时的结构，所以我们只能用下面的两种语法，它们彼此是等价的：    

```go
w = Wheel{Circle{Point{8,8}, 5}, 20}

w = Wheel {
  Circle: Circle{
    Point: Point {X: 8, Y: 8},
    Radius: 5,
  },
  Spokes: 20,
}
```    

上面的第二种分号有点多余啊。     

因为匿名成员也有一个隐式的名字，因此不能同时包含两个类型相同的匿名成员，这会导致名字冲突。
同时，因为成员的名字是由其类型隐式地决定的，所有匿名成员也有可见性的规则约束，在上面的
例子中，Point 和 Circle 匿名成员都是导出的。即使它们不导出（比如改成小写字母开头的
point 和 circle），我们依然可以用简短形式访问匿名成员嵌套的成员：    

`w.X = 8   // 等价于 w.circle.point.X = 9`    

但是在包外部，因为 circle 和 point 没有导出，不能访问它们的成员。因此简短的匿名成员访问
语法也是禁止的。     

说了半天也没弄懂这个可见性的问题。     

简短的点运算符语法可以用于选择匿名成员嵌套的成员，也可以
用于访问它们的方法。实际上，外层的结构体不仅仅是获得了匿名成员类型的所有成员，而
且也获得了该类型导出的全部的方法。       

### 4.4.4 数组、结构体、slice、map 比较


类型 | 数组 | 结构体 | slice | map
---------|----------|---------|---------|--------- 
 所属大类 | 聚合类型 | 聚合类型 | 引用类型 | 引用类型
 是否可比较 | 如果元素类型可比较，则数组可比较 | 如果全体成员可比较，则结构体可比较 | 除与 nil 外皆不可 | 除与 nil 外皆不可
 所占内存空间是否一定 | 一定 | 一定 | 不一定 | 不一定
 是否可以使用 make 函数创建 | 不可 | 不可 | 可以 | 可以

## 4.5 JSON

```go
type Movie struct {
  Title string
  Year int `json:"released"`
  Color bool `json:"color,omitempty"`
  Actors []string
}
var movies = []Movie{
  {Title: "Casablanca", Year: 1942, Color: false,
    Actors: []string{"Humphrey Bogart", "Ingrid Bergman"}},
  {Title: "Cool Hand Luke", Year: 1967, Color: true,
    Actors: []string{"Paul Newman"}},
  {Title: "Bullitt", Year: 1968, Color: true,
    Actors: []string{"Steve McQueen", "Jacqueline Bisset"}},
  // ...
}
```    

将一个 Go 语言中类似 movies 的结构体 slice 转为 JSON 的过程叫编组 marshaling，编组
通过调用 `json.Marshal` 函数完成：     

```go
data, err := json.Marshal(movies)
if err != nil {
  log.Fatalf("JSON marshaling failed: %s", err)
}
fmt.Printf("%s\n", data)
```     

Marshal 函数返回一个编码后的字节 slice，包含很长的字符串，并且没有空白缩进（难道不就是
字符串，还字节 slice？）：     

```
[{"Title":"Casablanca","released":1942,"Actors":["Humphrey Bogart","Ingr
id Bergman"]},{"Title":"Cool Hand Luke","released":1967,"color":true,"Ac
tors":["Paul Newman"]},{"Title":"Bullitt","released":1968,"color":true,"
Actors":["Steve McQueen","Jacqueline Bisset"]}]
```     

这种紧凑的表示形式虽然包含了全部的信息，但是很难阅读。为了生成便于阅读的格式，另
一个json.MarshalIndent函数将产生整齐缩进的输出。该函数有两个额外的字符串参数用于表
示每一行输出的前缀和每一个层级的缩进：    

```go
data, err := json.MarshalIndent(movies, "", "  ")
if err != nil {
  log.Fatalf("JSON marshaling failed: %s", err)
}
fmt.Printf("%s\n", data)
```    

上面的代码将产生这样的输出：    

```json
[
  {
    "Title": "Casablanca",
    "released": 1942,
    "Actors": [
      "Humphrey Bogart",
      "Ingrid Bergman"
    ]
  },
  {
    "Title": "Cool Hand Luke",
    "released": 1967,
    "color": true,
    "Actors": [
      "Paul Newman"
    ]
  },
  {
    "Title": "Bullitt",
    "released": 1968,
    "color": true,
    "Actors": [
      "Steve McQueen",
      "Jacqueline Bisset"
    ]
  }
]
```     

一个 JSON 数组可以通过将 Go 语言的数组和 slice编码得到。JSON 的对象类型可以通过 Go
语言的 map 类型（key类型是字符串）和结构体编码得到。     

在编码结构体时，默认使用 Go 语言结构体的成员名字作为 JSON 对象的键名。只有导出的
结构体成员才会被解码，也就是说只有大写字母开头的成员名字才会被 JSON 编码把。    

细心的读者可能已经注意到，其中Year名字的成员在编码后变成了released，还有Color成员
编码后变成了小写字母开头的color。这是因为结构体成员Tag所导致的。一个结构体成员Tag是和
在编译阶段关联到该成员的元信息字符串：    

```go
Year int `json:"released"`
Color bool `json:"color, omitempty"`
```     

结构体的成员 Tag 可以是任意的字符串字面值，但是通常是一系列用空格分隔的key: "value"
键值对序列；因为值中含有双引号字符，因此成员 Tag 一般用原生字符串字面值的形式书写。
json 开头键名对应的值用于控制 encoding/json 包的编码和解码的行为，并且 encoding
下其他的包也遵循这个约定。   

成员 Tag 中 json 对应值的第一部分用于指定 JSON 字段的名字，比如将 Go 语言中的
TotalCount 成员对应到 JSON 中的 total_count 字段。Color 成员的 Tag 还带了
一个额外的 omitempty 选项，表示当 Go 语言结构体成员为空或零值时不生成 JSON
字段。     

编码的逆操作是解码，对应将 JSON 数据解码为 Go 语言的数据结构，Go 语言中一般叫
unmarshaling，通过 json.Unmarshal 函数完成。通过定义合适的 Go 语言数据结构，我们
可以选择性地解码 JSON 中感兴趣的成员。当 Unmarshal 函数调用返回，slice 将被只含有
Title 信息值填充，其他 JSON 成员将被忽略。     

```go
var titles []struct{ Title string }
err := json.Unmarshal(data, &titles)

if err != nil {
  log.Fatalf("JSON unmarshaling failed: %s", err)
}

fmt.Println(titles)   // [{Casablanca} {Cool Hand Luke} {Bullitt}]
```     

```go
package github

import (
	"time"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const IssuesURL = "https://api.github.com/search/issues"

type IssuesSearchResult struct {
	TotalCount  int  `json:"total_count"`
	Items				[]*Issue
}

type Issue struct {
	Number 		int
	HTMLURL		string  `json:"html_url"`
	Title     string
	State			string
	User			*User
	CreatedAt time.Time `json:"created_at"`
	Body			string
}

type User struct {
	Login		string
	HTMLURL string  `json:"html_url"`
}

func SearchIssues(terms []string) (*IssuesSearchResult, error) {
	q := url.QueryEscape(strings.Join(terms, " "))
	resp, getErr := http.Get(IssuesURL + "?q=" + q)
	if getErr != nil {
		return nil, getErr
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	var result IssuesSearchResult
	parseErr := json.NewDecoder(resp.Body).Decode(&result)
	if parseErr != nil {
		resp.Body.Close()
		return nil, parseErr
	}

	resp.Body.Close()
	return &result, nil
}
```     

在早些的例子中，我们使用了 json.Unmarshal 函数来将 JSON 格式的字符串解码为字节
slice。但是在这个例子中，我们使用了基于流式的解码器 json.Decoder，它可以从一个
输入流解码 JSON 数据，尽管这不是必须的。     

## 4.6 文本和 HTML 模板

text/template 和 html/template 等模板包提供了一种将变量值填充到一个文本或 HTML
格式的模板的机制。     

一个模板是一个字符串或一个文件，里面包含了一个或多个由双花括号包含的 `{{action}}`
对象。大部分的字符串只是按字面值打印，但是对于 actions 部分将触发其他的行为。每个
actions 都包含了一个用模板语言书写的表达式，一个 action 虽然简短但是可以输出复杂的
打印值，模板语言包含通过选择结构体的成员、调用函数或方法、表达式控制流和 range 循环
语句，以及其他实例化模板等诸多特性。      

```go
const templ = `{{ .TotalCount }} issues:
{{ range .Items }}-------------------------
Number:  {{ .Number }}
User:    {{ .User.Login }}
Title:   {{ .Title | printf "%.64s" }}
Age:     {{ .CreatedAt | daysAgo }} days
{{ end }}`
```    

对于每一个action，都有一个当前值的概念，对应点操作符。当前值最初被初始化为调用模板时
的参数，在当前例子中对应 github.IssuesSearchResult 类型的变量。模板中
`{{range .Items}}` 和 `{{end}}` 对应一个循环 action，循环每次迭代的当前中对应
当前的 Items 元素的值。     

在一个 action 中， `|` 操作符表示将前一个表达式的结果作为后一个函数的输入。在 Title
这一行的 action 中，第二个操作是一个 `printf` 函数，是一个基于 `fmt.Sprintf` 实现
的内置函数，所有模板都可以直接使用。对于 Age 部分，第二个动作是一个叫 `daysAgo` 的
函数，通过 `time.Slice` 函数将 `CreatedAt` 转换为过去的时间长度：    

```go
func daysAgo(t time.Time) int {
	return int(time.Slice(t).Hours() / 24)
}
```       

生成模板输出需要两个处理步骤。第一步是要分析模板并转为内部表示，然后基于指定的输入执行模板。分析
模板部分一般只需要执行一次。下面的代码创建并分析上面定义的模板 templ。注意方法调用链的顺序：
`template.New` 先创建并返回一个模板；Funcs 方法将 daysAgo 等自定义函数注册到模板中，并返回
模板；最后调用 Parse 函数分析模板。     

```go
report, err := template.New("report").
	Funcs(template.FuncMap{"daysAgo": daysAgo}).
	Parse(templ)
if err != nil {
	log.Fatal(err)
}
```     

生成模板输出需要两个处理步骤。第一步是要分析模板并转为内部表示，然后基于指定的输入执行模板。分析
模板部分一般只需要执行一次。下面的代码创建并分析上面定义的模板 templ。注意方法调用链的顺序：
`template.New` 先创建并返回一个模板；Funcs 方法将 daysAgo 等自定义函数注册到模板中，并返回
模板；最后调用 Parse 函数分析模板。     

```go
report, err := template.New("report").
	Funcs(template.FuncMap{"daysAgo": daysAgo}).
	Parse(templ)
if err != nil {
	log.Fatal(err)
}
```     

```go
var report = template.Must(template.New("issuelist").
	Funcs(template.FuncMap{"daysAgo": daysAgo}).
	Parse(templ))

func main() {
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	err = report.Execute(os.Stdout, result)
	if err != nil {
		log.Fatal(err)
	}
}
```    