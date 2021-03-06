# 第六章 方法

<!-- TOC -->

- [第六章 方法](#第六章-方法)
  - [6.1 方法声明](#61-方法声明)
  - [6.2 基于指针对象的方法](#62-基于指针对象的方法)
    - [6.2.1 指针 recevier 和值 receiver](#621-指针-recevier-和值-receiver)
    - [6.2.2 nil 也是一个合法的接收器类型](#622-nil-也是一个合法的接收器类型)
  - [6.3 通过嵌入结构体来扩展类型](#63-通过嵌入结构体来扩展类型)
  - [6.4 方法值和方法表达式](#64-方法值和方法表达式)

<!-- /TOC -->

从我们的理解来讲，一个对象其实也就是一个简单的值或者一个变量，在这个对象中会包含一些
方法，而一个方法则是一个一个和特殊类型关联的函数。一个面向对象的程序会用方法来表达其属性
和对应的操作，这样使用这个对象的用户就不需要直接去操作对象，而是借助方法来做这些事情。    

## 6.1 方法声明

在函数声明时，在其名字之前放在一个变量，既是一个方法。这个附加的参数会将该函数附加到这种
类型上，即相当于为这种类型定义了一个独占的方法。     

```go
import "math"

type Point struct {
	X, Y float64
}

func Distance(p, q Point) float64 {
	return math.Hypot(q.X - p.X, q.Y - p.Y)
}

func (p Point) Distance(q Point) float64 {
	return math.Hypot(q.X - p.X, q.Y - p.Y)
}
```    

参数 p，叫做方法的接收器 receiver。注意其实 p 也是参数，只是写在了函数名的前部。    

所以方法其实就是一个带有一个 receiver 参数的函数。   

我们只能为在同一包中声明的类型的 receiver 定义方法，而不能为其他包中定义的类型的 receiver
定义方法，即使是 int 之类的内置类型也不可。简单的说就是 receiver 的类型定义和方法必须
定义在同一包内。   

在 Go 语言中，我们并不会想其他语言那样使用 this 或者 self 作为接收器，我们可以任意的选择
接收器的名字。由于接收器的名字经常会被使用到，所以保持其方法间传递时的一致性和简短性是不错
的主意。这里的建议是使用其类型的第一个字母。    

```go
p := Point{1,2}
q := Point{4,6}
fmt.Println(Distance(p,q))  // 5
fmt.Println(p.Distance(q))  // 5
```    

这种 p.Distance 的表达式叫做选择器，因为他会选择合适的对应 p 这个对象的 Distance 方法
来执行。选择器也会被用来选择一个 struct 类型的字段，比如 p.X。由于方法和字段都是在同一
命名空间，如果我们在这里声明一个 X 方法的话，编译器会报错。     

```go
type Path []Point

func (path Path) Distance() float64 {
  sum := 0.0
  for i := range path {
    if i > 0 {
      sum += path[i-1].Distance(path[i])
    }
  }
  return path
}
```    

Path 是一个命名的 slice 类型，而不是 Point 那样的 struct 类型，然而我们依然可以为它定义
方法。在能够给任意类型定义方法这一点上，Go 和很多其他的面向对象的语言不太一样。因此在 Go
语言里，我们为一些简单的数值、字符串、slice、map 来定义一些附加行为很方便。方法可以被
声明到任意类型，只要不是一个指针或者一个 interface。     

## 6.2 基于指针对象的方法

当调用一个函数时，会对其每一个参数值进行拷贝，如果一个函数需要更新一个变量，或者函数的其中
一个参数实在太大我们希望能够避免进行这种默认的拷贝，这种情况下我们就需要用到指针了。对应
到我们这里用来更新接收器的对象的方法，我们就可以用其指针而不是对象来声明方法：    

```go
func (p *Point) ScaleBy(factor float64) {
  p.X *= factor
  p.Y *= factor
}
```    

这个方法的名字是 `(*Point).ScaleBy`。这里的括号是必须的；没有括号的话这个表达式可能会被
理解为 `*(Point.ScaleBy)`。      

**问题**：是因为是结构体所以不需要加指针符吗？？？？？

只有类型(Point) 和指向他们的指针(*Point),才是可能会出现在接收器声明里的两种接收器。此外，
为了避免歧义，在声明方法时，如果一个类型名本身是一个指针的话，是不允许其出现在接收器中的：   

```go
type P *int
func (P) f() {/*...*/}  // compile error: invalid receiver type
```    

简单来说就是可以定义一个指针 receiver，不过这个指针的原本类型不能再是一个指针了。   

想要调用指针类型方法 (*Point).ScaleBy，只要提供一个 Point 类型的指针即可：    

```go
r := &Point{1,2}
r.ScaleBy(2)
fmt.Println(*r)   // {2,4}
```    

或者这样：   

```go
p := Point{1,2}
pptr := &p
pptr.ScaleBy(2)
fmt.Println(p)
```    

或者这样：    

```go
p := Point{1,2}
(&p).ScaleBy(2)
fmt.Println(p)
```     

不过后两种方法有些笨拙，go 语言本身在这种地方会帮到我们，如果接收器 p 是一个 Point 类型的
变量，并且其方法需要一个 Point 指针作为接收器，我们可以用下面这种简短的写法：    

`p.ScaleBy(2)`     

编译器会隐式地帮我们用 &p 去调用 ScaleBy 这个方法，这种简写方法只适用于变量。我们不能通过一个
无法取到地址的接收器来调用指针方法，比如临时变量的内存地址就无法获取到：    

`Point{1,2}.ScaleBy(2)   // compile error`     

这里的几个例子可能让你有些困惑，所以我们总结一下：在每一个合法的方法调用表达式中，下面三种
情况里的任意一种都是可以的：    

+ 接收器的实际参数和其接收器的形参相同，比如两者都是类型 T 或者都是类型 *T
+ 接收器形参是类型 T，但接收器实参是类型 *T，这种情况下编译器会隐式地为我们取变量的地址：   
`p.ScaleBy(2) // implicit(&p)`    
+ 接收器形参是类型 *T，实参是类型 T，编译器会隐式地为我们解引用，取到指针指向的实际变量  
`pptr.Distance(q)  // implicit(*pptr)`      

### 6.2.1 指针 recevier 和值 receiver

当一个方法定义时的 receiver 是一个指针时，那么在调用时，不管 receiver 是对应类型的指针还是
类型的的值，都可以直接调用：    

```go
func (v *Vertex) Scale(f float) {
	v.X = v.X * f
	v.Y = v.Y * f
}

func (v Vertex) Abs() float64 {
	return math.Sqrt(v.X * v.X + v.Y * v.Y)
}

var v Vertex
v.Scale(5)   // OK
p := &v
p.Scale(5)   // OK
```    

具体的原因应该是由于 Scale 方法是一个指针 receiver 的方法，那么 Go 在遇到用值调用方法时
`v.Scale(5)` 会尝试将其解释为 `(&v).Scale(5)`，当然这个会不会是隐式编译为后一种形式，这个就
不太清楚了。     

同时，在相反的方向也会有这样的效果，如果一个方法接收一个值 receiver，那么如果你使用指针去调用，
结果也是 OK 的，原理也一样，`p.Abs()` 也会被解释为 `(*p).Abs()`。  


### 6.2.2 nil 也是一个合法的接收器类型

略。    

## 6.3 通过嵌入结构体来扩展类型

略。    

## 6.4 方法值和方法表达式

我们经常选择一个方法，并且在同一表达式里执行。比如常见的 p.Distance() 形式，实际上将其
分成两步来执行也是可能的。p.Distance 叫做选择器，选择器会返回一个方法值——一个将方法绑定
到特定接收器变量的函数。这个函数可以不通过指定其接收器即可被调用；即调用时不需要指定接收器。
只要传入函数的参数即可：    

```go
p := Point{1,2}
q := Point{4,6}

distanceFromP := p.Distance
fmt.Println(distanceFormP(q))
```     

和方法值有关的还有方法表达式。当 T 是一个类型时，方法表达式可能会写作 T.f 或者 (*T).f，
会返回一个函数值，这种函数会将其第一个参数用作接收器：    

```go
p := Point{1,2}
q := Point{4,6}

distance := Point.Distance
fmt.Println(distance(p, q))
```    

