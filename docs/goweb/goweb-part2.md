# Go Web 编程 Part2

<!-- TOC -->

- [Go Web 编程 Part2](#go-web-编程-part2)
- [第 5 章 内容展示](#第-5-章-内容展示)
  - [5.1 模板引擎](#51-模板引擎)
  - [5.2 Go 的模板引擎](#52-go-的模板引擎)
    - [5.2.1 对模板进行语法分析](#521-对模板进行语法分析)
    - [5.2.2 执行模板](#522-执行模板)
  - [5.3 动作](#53-动作)
    - [5.3.1 条件动作](#531-条件动作)
    - [5.3.2 迭代动作](#532-迭代动作)
    - [5.3.3 设置动作](#533-设置动作)
    - [5.3.4 包含动作](#534-包含动作)
  - [5.4 参数、变量和管道](#54-参数变量和管道)
  - [5.5 函数](#55-函数)
  - [5.6 上下午感知](#56-上下午感知)
    - [5.6.1 不对 HTML 进行转义](#561-不对-html-进行转义)
  - [5.7 嵌套模板](#57-嵌套模板)
  - [5.8 通过块动作定义默认模板](#58-通过块动作定义默认模板)
- [第 6 章 存储数据](#第-6-章-存储数据)
  - [6.2 文件存储](#62-文件存储)
    - [6.2.1 读取和写入 CSV 文件](#621-读取和写入-csv-文件)
    - [6.2.2 gob 包](#622-gob-包)
  - [6.3 Go 与 SQL](#63-go-与-sql)

<!-- /TOC -->

# 第 5 章 内容展示

## 5.1 模板引擎

我们可以把模板引擎划分为两种理想的类型：    

+ 无逻辑模板引擎 —— 将模板中指定的占位符替换成相应的动态数据。这种模板引擎只进行字符串替换，而
不执行任何逻辑处理。
+ 嵌入逻辑模板引擎 —— 将变成语言代码嵌入模板当中，并在模板引擎渲染模板时，由模板引擎执行这些
代码并进行相应的字符串替换工作。    

Go标准库提供的模板引擎功能大部分都定义在了 text/template 库当中,而小部分与HTML相关的功能
则定义在了 html/template 库里面。这两个库相辅相成:用户可以把这个模板引擎当做无逻辑模板引擎使用,
但与此同时,Go也提供了足够多的嵌入式模板引擎特性,使这个模板引擎用起来既有趣又强大。   

## 5.2 Go 的模板引擎

模板中的动作默认使用两个大括号 `{{` 和 `}}` 包着，当然也可以通过模板引擎提供的方法自行指定其他
定界符。    

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta http-equiv="X-UA-Compatible" content="ie=edge">
  <title>Document</title>
</head>
<body>
  {{ . }}
</body>
</html>
```    

使用 Go 的 Web 模板引擎需要以下两个步骤：　　

1. 对文本格式的模板源进行语法分析，创建一个经过语法分析的模板结构，其中模板源既可以是一个字符串，
也可以是模板文件中包含的的内容。
2. 执行经过语法分析的模板，将 `ResponseWriter` 和模板所需的动态数据传递给模板引擎，被调用的模板
引擎会把经过语法分析的模板和传入的数据结合起来，生成最终的 HTML，并将这些 HTML 传递给 ResponseWriter   

```go
package main

import (
  "net/http"
  "html/template"
)

func process(w http.ResponseWriter, r *http.Request) {
  t, _ := template.ParseFiles("tmpl.html")
  t.Execute(w, "Hello World!")
}

func main() {
  server := http.Server{
    Addr: "127.0.0.1:8080",
  }

  http.HandleFunc("/process", process)
  server.ListenAndServe()
}
```    

### 5.2.1 对模板进行语法分析

`ParseFiles` 是一个独立的函数，它可以对模板文件进行语法分析，并创建出一个经过语法分析的模板
结构以供 `Execute` 方法执行。实际上，`ParseFiles` 只是为了方便地调用 `Template` 结构体
的 `ParseFiles` 方法而设置的一个函数，当用户调用 `ParseFiles` 函数的时候，Go 会创建一个
新的模板，并将用户给定的模板文件的名字用作这个新模板的名字：　　　

`t, _ := template.ParseFiles("tmpl.html")`    

这相当于创建一个新模板，然后调用它的 `ParseFiles`　方法：   

```go
t := template.New("tmpl.html")
t, _ := t.ParseFiles("tmpl.html")
```   

无论是 `ParseFiles` 函数还是 `Template` 结构的 `ParseFiles` 方法，它们都可以接受一个或多个
文件名作为函数，换句话说，这两个函数/方法都是可变参数函数/方法。但于此同时，它们都只返回一个模板。    

当用户向 `ParseFiles` 函数或 `ParseFiles` 方法传入多个文件时, `ParseFiles` 只会返回用户
传入的第一个文件的已分析模板,并且这个模板也会根据用户传入的第一个文件的名字进行命名;至于其他传入
文件的已分析模板则会被放置到一个映射里面,这个映射可以在之后执行模板时使用。   

对模板文件进行语法分析的另一种方法是使用 `ParseGlob` 函数,跟 `ParseFiles` 只会对给定文件
进行语法分析的做法不同, `ParseGlob` 会对匹配给定模式的所有文件进行语法分析。举个例子,如果目录里面只
有 tmpl.html 一个HTML文件,那么语句 `t, _ := template.ParseFiles("tmpl.html")`
和 `t, _ := template.ParseGlob("*.html")` 效果一致。其实就是文件通配符把。    

虽然Go语言的一般做法是手动地处理错误,但Go也提供了另外一种机制,专门用于处理分析模板时出现的错误:   

`t := template.Must(template.ParseFiles("tmpl.html"))`     

`Must` 函数可以包裹一个函数，被包裹的函数会返回一个指向模板的指针和一个错误，如果这个错误不是
nil，那么 `Must` 函数将产生一个 panic。    

### 5.2.2 执行模板

执行模板最常用的方法是调用模板的 `Execute` 方法。在只有一个模板的情况先，上面提到的这种方法总是
可行的，但如果模板不止一个，那么当对模板集合调用调用 `Execute` 方法的时候，`Execute` 方法
只会执行模板集合中的第一个模板。这时，我们就需要使用 `ExecuteTemplate`。    

```go
t, _ := template.ParseFiles("t1.html", "t2.html")
t.ExecuteTemplate(w, "t2.html", "Hello World")
```    

## 5.3 动作

Go 模板的动作就是一些嵌入在模板里面的命令：   

+ 条件动作
+ 迭代动作
+ 设置动作
+ 包含动作
+ 定义动作

. 也是一个动作，它代表的是传递给模板的数据，其他动作和函数基本上都会对这个动作进行处理。    

### 5.3.1 条件动作

条件动作会根据参数的值来决定对多条语句中的哪一条语句进行求值：    

```
{{ if arg }}
  some content
{{ end }}
```   

这个动作的另一种形式如下：   

```
{{ if arg }}
  some content
{{ else }}
  other content
{{ end }}
```     

### 5.3.2 迭代动作

迭代动作可以对数组、切片、映射或者通道进行迭代，而在迭代循环的内部，. 则会被设置为当前迭代的元素：   

```
{{ range array }}
  Dot is set to the element {{ . }}
{{ end }}
```    

下面是迭代动作的一个变种，这个变种允许用户在被迭代的数据结构为空时，显示一个备选的结果。   

```
{{ range . }}
<li>{{ . }}</li>
{{ else }}
<li>Nothing to show</li>
{{ end }}
```   

模板里面介于 {{ else }} 和 {{ end }} 之间的内容将在 . 为 nil 的时候显示。    


### 5.3.3 设置动作

设置动作允许用户在指定的范围之内为 . 设置值：   

```
{{ with arg }}
  Dot is set to arg
{{ end }}
```    

介于 `{{ with arg }}` 和 `{{ end }}` 之间的点将被设置为参数 arg 的值。    

跟迭代动作一样，设置动作也拥有一个能够提供备选方案的变种：   

```
{{ with arg }}
  Dot is set to arg
{{ else }}
  Fallback if arg is empty
{{ end }}
```    

### 5.3.4 包含动作

包含动作允许用户在一个模板里面包含另一个模板，从而构建出嵌套的模板。包含动作的格式为 `{{ template "name" }}`
，其中 name 参数为被包含模板的名字。    

```html
<!-- t1.html -->
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta http-equiv="X-UA-Compatible" content="ie=edge">
  <title>Document</title>
</head>
<body>
  <div>This is t1.html before</div>
  <div>This is the value of the dot in t1.html - [{{ . }}] </div>
  <hr />
  {{ template "t2.html" }}
  <hr />
  <div>This is t1.html after</div>
</body>
</html>

<!-- t2.html -->
<div style="background-color: yellow">
  This is t2.html<br />
  This is the value of the dot in t2.html - [{{ . }}]
</div>
```   

```go
func process(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("t1.html", "t2.html")
	t.Execute(w, "Hello World!")
}
```   

在执行嵌套模板时，我们必须对涉及的所有模板文件都进行语法分析。    

因为模板 t1.html 并没有把字符串 "Hello World!" 也传递为空字符串。为了向被嵌套的模板传递
数据，用户可以使用包含动作的变种 `{{ template "name" arg }}`，其中 arg 就是用户想要传递
给被嵌套模板的数据。   

## 5.4 参数、变量和管道

一个参数就是模板中的一个值。它可以是布尔值、整数、字符串等字面量，也可以是结构、结构中的一个字段
或者数组中的一个键。除此之外，参数还可以是一个变量、一个方法（这个方法必须只返回一个值，或者只
返回一个值和一个错误）或者一个函数。最后，参数也可以是一个点 .，用于表示处理器向模板引擎传递的
数据。    

比如说，在以下这个例子中，arg 是一个参数：   

```go
{{ if arg }}
  some content
{{ end }}
```   

除了参数之外，用户还可以在动作中设置变量。变量以美元符号开头： `$variable := value`。   

```go
{{ range $key, $value := . }}
  The key is {{ $key }} and the value {{ $value }}
{{ end }}
```    

模板中的管道是多个有序地串联起来的参数、函数和方法：    

`{{ p1 | p2 | p3 }}`    

这里的 p1, p2, p3 可以是参数或函数。管道允许用户将一个参数的输出传递给下一个参数，而各个参数之间
则使用 | 分隔。    

## 5.5 函数

Go 模板引擎内置了一些非常基础的函数，其中包括为 `fmt.Sprint` 的不同变种创建的几个别名函数，
用户可以自行定义自己想要的函数。    

需要注意的是,Go的模板引擎函数都是受限制的:尽管这些函数可以接受任意多个参数作为输入,但它们只能
返回一个值,或者返回一个值和一个错误。    

为了创建一个自定义模板函数，用户需要：   

1. 创建一个名为 `FuncMap` 的映射，并将映射的键设置为函数的名字，而映射的值则设置为实际定义的函数
2. 将 `FuncMap` 与模板进行绑定    

```go
func formatDate(t time.Time) string {
	layout := "2006-01-02"
	return t.Format(layout)
}

func process(w http.ResponseWriter, r *http.Request) {
	funcMap := template.FuncMap{"fdate": formatDate}
	t := template.New("tmpl.html").Funcs(funcMap)
	t, _ = t.ParseFiles("tmpl.html")
	t.Execute(w, time.Now())
}
```   

```
<div>The date/time is {{ . | fdate }}</div>
```   

除此之外，我们也可以像调用普通函数一样，将点 . 作为参数传递给 fdate 函数：   

`<div>The date/time is {{ fdate . }}</div>`    

## 5.6 上下午感知

上下文感知的一个显而易见的用途就是对被显示的内容实施正确的转义(escape):这意味着,如果模板显示
的是HTML格式的内容,那么模板将对其实施HTML转义;如果模板显示的是JavaScript格式的内容,那么模板
将对其实施JavaScript转义;诸如此类。除此之外,Go模板引擎还可以识别出内容中的URL或者CSS样式。   

上下文感知特性主要用于实现自动的防御编程,并且它使用起来非常方便。通过根据上下文对内容进行修改,
Go模板可以防止某些明显并且低级的编程错误。    

### 5.6.1 不对 HTML 进行转义

只要将不想被转义的内容传给 template.HTML 函数，模板引擎就不会对其进行转义。   

```go
func process(w http.ResponseWriter, r *http.Request) {
  t, _ := template.ParseFiles("tmpl.html")
  t.Execute(w, template.HTML(r.FormValue("comment")))
}
```    

## 5.7 嵌套模板

可以在同一个模板文件中定义多个不同的模板：   

```html
{{ define "layout" }}
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta http-equiv="X-UA-Compatible" content="ie=edge">
  <title>Document</title>
</head>
<body>
  {{ template "content" }}
</body>
</html>

{{ end }}

{{ define "content" }}
  Hello World
{{ end }}
```   

```go
func process(w http.ResponseWriter, r *http.Request) {
  t, _ := template.ParseFiles("layout.html")
  t.ExecuteTemplate(w, "layout", "")
}
```   

因为 layout 模板嵌套了 content 模板, 所以程序只需要执行 layout 模板就可以在浏览器中得到 content 模板产生的 Hello World! 输出了。    

用户除可以在同一个模板文件里面定义多个不同的模板之外,还可以在不同的模板文件里面定义同名的模板。   

## 5.8 通过块动作定义默认模板

Go 1.6 引入了一个新的块动作，这个动作允许用户定义一个模板并立即使用。块动作看上去是下面这个样子的:   

```go
{{ block arg }}
  Dot is set to arg
{{ end }}
```   

```go
{{ block "content" . }}
  ....
{{ end }}
``` 

当 layout 模板被执行时,如果模板引擎没有找到可用的 content 模板,那么它就会使用块动作中定义的 content 模板。    

# 第 6 章 存储数据

## 6.2 文件存储

把数据存储到非易失存储器里面同样也有多种方法可选,而本节要介绍的是把数据存储到文件系统里面的相关
技术。说得更具体一点,我们将要学习的是如何通过 Go 语言以两种不同的方式将数据存储到文件里面:第一种
方式需要用到通用的CSV(comma-separated value,逗号分隔值)文本格式,而第二种方法则需要用到Go
语言特有的 gob 包。   

CSV是一种常见的文件格式,用户可以通过这种格式向系统传递数据。当你需要用户提供大量数据,但是却因为
某些原因而无法让用户把数据填入你提供的表单时,CSV格式就可以派上用场了:你只需要让用户使用电子表格
程序(spreadsheet)输入所有数据,然后将这些数据导出为CSV文件,并将其上传到你的Web应用中,这样就
可以在获得CSV文件之后,根据自己的需要对数据进行解码。同样地,你的Web应用也可以将用户的数据打包成
CSV文件,然后通过向用户发送CSV文件来为他们提供数据。    

gob是一种能够存储在文件里面的二进制格式,这种格式可以快速且高效地将内存中的数据序列化到一个或多个文件里面。    

```go
package main


import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	data := []byte("Hello World!\n")
	err := ioutil.WriteFile("data1", data, 0644)
	if err != nil {
		panic(err)
	}

	read1, _ := ioutil.ReadFile("data1")
	fmt.Print(string(read1))

	file1, _ := os.Create("data2")
	defer file1.Close()

	bytes, _ := file1.Write(data)
	fmt.Printf("Wrote %d bytes to file\n", bytes)

	file2, _ := os.Open("data2")
	defer file2.Close()

	read2 := make([]byte, len(data))
	bytes, _ = file2.Read(read2)

	fmt.Printf("Read	%d	bytes	from	file\n",	bytes)
  fmt.Println(string(read2))
}
```    

在这个代码清单里面,程序使用了两种不同的方法来对文件进行写入和读取。第一种方法非常简单直接,它使用
的是 ioutil 包中的 `WriteFile` 函数和 `ReadFile` 函数:在写入文件时,程序会将文件的名字、
需要写入的数据以及一个用于设置文件权限的数字用作参数调用 `WriteFile` 函数;而在读取文件时,程序
只需要将文件的名字用作参数,然后调用 `ReadFile` 函数即可。此外,无论是传递给 `WriteFile` 的
数据,还是 `ReadFile` 返回的数据,都是一个由字节组成的切片。    

### 6.2.1 读取和写入 CSV 文件

对Go语言来说,CSV文件可以通过 encoding/csv 包进行操作：    

```go
package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type Post struct {
	Id		int
	Content  string
	Author	 string
}

func main() {
	csvFile, err := os.Create("posts.csv")
	if err != nil {
		panic(err)
	}

	defer csvFile.Close()

	allPosts := []Post{
		Post{1, "Hello World!", "Sau Sheong"},
		Post{2, "Bonjour Monde!", "Pierre"},
		Post{3, "Hola Mundo!", "Pedro"},
		Post{4, "Greetings Earthlings", "Sau Sheong"},	// 为什么这里也需要逗号
	}

	writer := csv.NewWriter(csvFile)

	for _, post := range allPosts {
		line := []string{strconv.Itoa(post.Id), post.Content, post.Author}
		err := writer.Write(line)
		if err != nil {
			panic(err)
		}
	}
	writer.Flush()

	file, err := os.Open("posts.csv")
	if err != nil {
		panic(err)
	}

	defer file.Close()
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	record, err := reader.ReadAll()
	if err != nil {
		panic(nil)
	}

	var posts []Post
	for _, item := range record {
		id, _ := strconv.ParseInt(item[0], 0, 0)
		post := Post{Id: int(id), Content: item[1], Author: item[2],}
		posts = append(posts, post)
	}

	fmt.Println(posts[0].Id)
	fmt.Println(posts[0].Content)
	fmt.Println(posts[0].Author)
}
```    


程序将读取器的 `FieldsPerRecord` 字段的值设置为负数,这样的话,即使读取器在读取时发现记录
(record)里面缺少了某些字段,读取进程也不会被中断。反之,如果 `FieldsPerRecord` 字段的值为正数,
那么这个值就是用户要求从每条记录里面读取出的字段数量,当读取器从CSV文件里面读取出的字段数量少于
这个值时,Go就会抛出一个错误。最后,如果 `FieldsPerRecord` 字段的值为 0 ,那么读取器
就会将读取到的第一条记录的字段数量用作 FieldsPerRecord 的值。    

### 6.2.2 gob 包

encoding/gob 包用于管理由gob组成的流(stream),这是一种在编码器(encoder)和解码器(decoder)
之间进行交换的二进制数据,这种数据原本是为序列化以及数据传输而设计的,但它也可以用于对数据进行持久
化,并且为了让用户能够方便地对文件进行读写,编码器和解码器一般都会分别包裹起程序的写入器以及读取器。    

```go
package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
)

type Post struct {
	Id   int
	Content  string
	Author   string
}

func store(data interface{}, filename string) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(data)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(filename, buffer.Bytes(), 0600)
	if err != nil {
		panic(err)
	}
}

func load(data interface{}, filename string) {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(nil)
	}

	buffer := bytes.NewBuffer(raw)
	dec := gob.NewDecoder(buffer)
	err = dec.Decode(data)
	if err != nil {
		panic(err)
	}
}

func main() {
	post := Post{Id: 1, Content: "Hello World!", Author: "Sau Sheong"}
	store(post, "post1")
	var postRead Post
	load(&postRead, "post1")
	fmt.Println(postRead)
}
```    

## 6.3 Go 与 SQL

```go
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Post struct {
	Id int
	Content string
	Author string
}

var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("postgres", "user=gwp dbname=gwp password=gwp sslmode=disable")
	if err != nil {
		panic(err)
	}
}

func Posts(limit int) (posts []Post, err error) {
	rows, err := Db.Query("select id, content, author from posts limit $1", limit)
	if err != nil {
		return
	}

	for rows.Next() {
		post := Post{}
		err = rows.Scan(&post.Id, &post.Content, &post.Author)
		if err != nil {
			return
		}
		posts = append(posts, post)
	}
	rows.Close()
	return
}

func GetPost(id int) (post Post, err error) {
	post = Post{}
	err = Db.QueryRow("select id, content, author from posts where id = $1", id).Scan(&post.Id, &post.Content, &post.Author)
	return
}

func (post *Post) Create() (err error) {
	statement := "insert into posts (content, author) values($1, $2) returning id"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(post.Content, post.Author).Scan(&post.Id)
	return
}

func (post *Post) Update() (err error) {
	_, err = Db.Exec("update posts set content = $2, author = $3 where id = $1", post.Id, post.Content, post.Author)
	return
}

func (post *Post) Delete() (err error) {
	_, err = Db.Exec("delete from posts where id = $1", post.Id)
	return
}

func main() {
	post := Post{Content: "Hello World!", Author: "Sau Sheong"}

	fmt.Println(post)
	post.Create()
	readPost, _ := GetPost(post.Id)
	fmt.Println(readPost)

	readPost.Content = "Bonjour Monde!"
	readPost.Author = "Pierre"
	readPost.Update()

	posts, _ := Posts(1)
	fmt.Println(posts)

	readPost.Delete()
}

```   

程序在对数据库执行任何操作之前,都需要先与数据库进行连接。sql.DB 结构是一个数据库句柄(handle),
它代表的是一个包含了零个或任意多个数据库连接的连接池(pool),这个连接池由 sql 包管理。程序可以
通过调用 Open 函数,并将相应的数据库驱动名字(driver name)以及数据源名字(data	source name)
传递给该函数来建立与数据库的连接。    

需要注意的是, Open 函数在执行时并不会真正地与数据库进行连接,它甚至不会检查用户给定的参数: 
Open 函数的真正作用是设置好连接数据库所需的各个结构,并以惰性的方式,等到真正需要时才建立相
应的数据库连接。     

Create 方法做的第一件事是定义一条SQL预处理语句,一条预处理语句(prepared statement)就是一个
SQL语句模板,这种语句通常用于重复执行指定的SQL语句,用户在执行预处理语句时需要为语句中的参
数提供实际值。     

`stmt, err := db.Prepare(statement)`    

这行代码会创建一个指向 sql.Stmt 接口的引用,这个引用就是上面提到的预处理语句。之后,程序会调用
预处理语句的 QueryRow 方法,并把来自接收者的数据传递给该方法,以此来执行预处理语句。     