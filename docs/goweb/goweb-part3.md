# Go Web 编程 Part3

<!-- TOC -->

- [Go Web 编程 Part3](#go-web-编程-part3)
- [第 7 章 Go Web 服务](#第-7-章-go-web-服务)
  - [7.1 Web 服务简介](#71-web-服务简介)
  - [7.2 基于 SOAP 的 Web 服务简介](#72-基于-soap-的-web-服务简介)
  - [7.3 基于 REST 的 Web 服务简介](#73-基于-rest-的-web-服务简介)
    - [7.3.1 将动作转换为资源](#731-将动作转换为资源)
    - [7.3.2 将动作转换为资源的属性](#732-将动作转换为资源的属性)
  - [7.4 通过 Go 分析和创建 XML](#74-通过-go-分析和创建-xml)
  - [7.5 通过 Go 分析和创建 JSON](#75-通过-go-分析和创建-json)
  - [7.5.1 分析 JSON](#751-分析-json)
- [第 8 章 应用测试](#第-8-章-应用测试)
  - [8.2 使用 Go 进行单元测试](#82-使用-go-进行单元测试)
    - [8.2.1 跳过测试用例](#821-跳过测试用例)
    - [8.2.2 以并行方式运行测试](#822-以并行方式运行测试)
    - [8.2.3 基准测试](#823-基准测试)
  - [8.3 使用 Go 进行 HTTP 测试](#83-使用-go-进行-http-测试)
- [第 9 章 发挥 Go 的并发优势](#第-9-章-发挥-go-的并发优势)
  - [9.1 并发与并行的区别](#91-并发与并行的区别)
  - [9.2 goroutine](#92-goroutine)
    - [9.2.1 等待 goroutine](#921-等待-goroutine)
- [第 10 章 Go 的部署](#第-10-章-go-的部署)
  - [10.1 将应用部署到独立服务器](#101-将应用部署到独立服务器)

<!-- /TOC -->

# 第 7 章 Go Web 服务

## 7.1 Web 服务简介

Web服务是一个软件系统,它的目的是为网络上进行的可互操作机器间交互(interoperable machine-to-machine	interaction)
提供支持。每个Web服务都拥有一套自己的接口,这些接口由一种名为Web服务描述语言
(web	service	description	language,WSDL)的机器可处理格式描述。其他系统需要根据Web服务的
描述,使用SOAP消息与Web服务交互。为了与其他Web相关标准实现协作,SOAP消息通常会被序列化为XML并通过HTTP传输。    

似乎所有Web服务都应该基于SOAP来实现,但实际中却存在着多种不同类型的Web服务,其中包括基于SOAP的、
基于REST的以及基于XML-RPC的。    

基于SOAP的Web服务出现的时间较早,W3C工作组已经对其进行了标准化,与之相关的文档和资料也非常丰富。
基于SOAP的服务不仅健壮、能够使用WSDL进行明确的描述、拥有内置的错误处理机制,而且还可以通过UUDI
(Universal Description,Discovery,	and	Integration,统一描述、发现和集成)(一种目录服务)规范发布。    

SOAP的缺点也是非常明显的:它不仅笨重,而且过于复杂。SOAP的XML报文可能会变得非常冗长,导致难以调试,
使用户只能通过其他工具对其进行管理,而基于SOAP的 Web 服务可能会因为额外的资源损耗而无法高效地运行。
此外,WSDL 虽然在客户端和服务器之间提供了坚实的契约,但这种契约有时候也会变成一种累赘:为了对Web
服务进行更新,用户必须修改WSDL,而这种修改又会引起SOAP客户端发生变化。    

跟基于SOAP的Web服务比起来,基于REST的Web服务就显得灵活多了。REST本身并不是一种结构,而是一种
设计理念。很多基于 REST 的 Web 服务都会使用像 JSON 这样较为简单的数据格式而不是XML,从
而使Web服务可以更高效地运行。    

## 7.2 基于 SOAP 的 Web 服务简介

SOAP是一种协议,用于交换定义在XML里面的结构化数据,它能够跨越不同的网络协议并在不同的编程模型中使用。
SOAP原本是Simple Object Access Protocol(简单对象访问协议)的首字母缩写。    

因为SOAP不仅高度结构化,而且还需要严格地进行定义,所以用于传输数据的XML可能会变得非常复杂。
WSDL 是客户端与服务器之间的契约,它定义了服务提供的功能以及提供这些功能的方式,服务的每
个操作以及输入/输出都需要由 WSDL 明确地定义。    

## 7.3 基于 REST 的 Web 服务简介

在使用 REST 设计的情况下,一个应用要如何才能激活一个用户的账号呢?因为REST只允许用户使用指定的几个
HTTP 方法操纵资源,而不允许用户对资源执行任意的动作,所以应用是无法发送像下面这样的请求的：   

`ACTIVATE /user/456 HTTP/1.1`    

有一些办法可以绕过这个问题,下面是最常用的两种方法:   

1. 把过程具体化,或者把动作转换成名词,然后将其用作资源
2. 将动作用作资源的属性   

### 7.3.1 将动作转换为资源

对于上面列举的例子,我们可以把对用户的激活动作转换为对资源的激活动作,然后通过向资源发送 HTTP 
方法来执行激活动作,这样一来,我们就可以通过以下方法激活指定的用户:    

```
POST /user/456/activation HTTP/1.1

{ "date": "2015-05-15T13:05:05Z" }
```    

### 7.3.2 将动作转换为资源的属性

如果用户的激活与否可以通过用户账号的一个状态来确定,那么我们只需要将激活动作用作资源的属性,然后
通过HTTP的 PATCH 方法对该资源进行部分更新即可,就像这样:   

```
PATCH /user/456 HTTP/1.1

{ "active": "true" }
```    

## 7.4 通过 Go 分析和创建 XML

在 Go 语言里面，用户首先需要将 XML 的分析结果存储到一些结构里面，然后通过访问这些结构来获取 XML
记录的数据。下面是分析 XML 时常见的两个步骤：   

1. 创建一些用于存储 XML 数据的结构
2. 使用 `xml.Unmarshal` 将 XML 数据解封到结构里面

```xml
<?xml version="1.0" encoding="utf-8"?>
<post id="1">
  <content>Hello World!</content>
  <author id="2">Sau Sheong</author>
</post>
```   

```go
package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type Post struct {
	XMLName xml.Name `xml:"post"`
	Id string `xml:"id,attr"`
	Content string `xml:"content"`
	Author Author `xml:"author"`
	XML string `xml:",innerxml"`
}

type Author struct {
	Id string `xml:"id,attr"`
	Name string `xml:",chardata"`
}

func main() {
	xmlFile, err := os.Open("post.xml")
	if err != nil {
		fmt.Println("Error opening XML file:", err)
		return
	}
	defer xmlFile.Close()
	xmlData, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		fmt.Println("Error reading XML data:", err)
		return
	}

	var post Post
	xml.Unmarshal(xmlData, &post)
	fmt.Println(post)
}
```   

结构标签是一些跟在字段后面，使用字符串表示的键值对：它的键是一个不能包含空格、引号或者冒号的
字符串，而值则是一个被双引号包围的字符串。在处理 XML 时，结构标签的键总是为 xml。   

出于创建映射的需要, xml 包要求被映射的结构以及结构包含的所有字段都必须是公开的,也就是,它们的
名字必须以大写的英文字母开头。    

下面是 XML 结构标签中的一些使用规则：    

1. 通过创建一个名字为 XMLName、类型为 xml.Name 的字段,可以将XML元素的名字存储在这个字段里面
2. 通过创建一个与XML元素属性同名的字段,并使用 `xml:"<name>,attr"` 作为该字段的结构标签,
可以将元素的 `<name>` 属性的值存储到这个字段里面。
3. 通过创建一个与XML元素标签同名的字段,并使用 `xml:",chardata"` 作为该字段的结构标签,可以
将XML元素的字符数据存储到这个字段里面。
4. 通过定义一个任意名字的字段,并使用 `xml:",innerxml"` 作为该字段的结构标签,可以将XML元素
中的原始XML存储到这个字段里面。
5. 没有模式标志(如 ,attr、,chardata 或者 ,innerxml )的结构字段将与同名的XML元素匹配。
6. 使用 `xml:"a>b>c"` 这样的结构标签可以在不指定树状结构的情况下直接获取指定的XML元素,其中
a 和 b 为中间元素,而 c 则是想要获取的节点元素。    

简单解释一下，分析程序定义了与XML元素 post 同名的 Post 结构,虽然这种做法非常常见,但是在某些
时候,结构的名字与XML元素的名字可能并不相同,这时用户就需要一种方法来获取元素的名字。为此, xml 包提供了
一种机制,使用户可以通过定义一个名为 XMLName 、类型为 xml.Name 的字段,并将该字段映射至元素
自身来获取XML元素的名字。    

那这里其实是通过 `XMLName` 这一字段获取到了 XML 中的 post 元素，然后第二条、第三条规则创建
一个与元素属性或标签同名的 tag，这里其实在字段和标签中并没有出现元素 post 或者 author 的名字，
这里应该是因为我们已经通过 XMLName 以及 Author 将元素映射到当前的结构中了，所以其实这个结构体
中的所有 tag 都是针对当前结构体所映射的 XML 元素展开的。    

上面的程序打印出的结果是这样的：   

```
{{ post} 1 Hello World! {2 Sau Sheong}
	<content>Hello World!</content>
	<author id="2">Sau Sheong</author>
}
```    

我们上面的做法虽然能够很好地处理体积较小的XML文件,但是却无法高效地处理以流(stream)方式传输的
XML文件以及体积较大的XML文件。为了解决这个问题,我们需要使用 Decoder 结构来代替 Unmarshal
函数,通过手动解码XML元素的方式来解封XML数据。    

对XML进行解码首先需要创建一个 Decoder ,这一点可以通过调用 NewDecoder 并向其传递一个
io.Reader 来完成。在上面展示的代码清单中,程序就把 os.Open 打开的 xmlFile 文件传递给了 NewDecoder。   

在拥有了解码器之后,程序就会使用 Token 方法来获取XML流中的下一个 token: 在这种情景下,token
实际上就是一个表示XML元素的接口。为了从解码器里面取出所有token,程序使用一个无限 for 循环包裹
起了从解码器里面获取token的相关动作。当解码器包含的所有token都已被取出时, Token 方法将返回
一个表示文件数据或数据流已被读取完毕的 io.EOF 结构作为结果,并将返回值中的 err 变量的值设置为 nil 。   

也就是处理这些编码的内容，有两种方案，一种是是文件整体传入，直接封装 Marshal 或者解封 Unmarshal，
另一种就是使用解码器/编码器以流的方式处理。    

## 7.5 通过 Go 分析和创建 JSON

## 7.5.1 分析 JSON

跟映射XML相比,把结构映射至JSON要简单得多,后者只有一条通用的规则:对于名字为 &lt;name&gt;
的JSON键,用户只需要在结构里创建一个任意名字的字段,并将该字段的结构标签设置为 'json:"&lt;name&gt;"',
就可以把JSON键 &lt;name&gt; 的值存储到这个字段里面。     

```json
{
  "id": 1,
  "content": "Hello World!",
  "author": {
    "id": 2,
    "name": "Sau Sheong"
  },
  "comments": [
    {
      "id": 3,
      "content": "Have a great day!",
      "author": "Adam"
    },
    {
      "id": 4,
      "content": "How are you today?",
      "author": "Betty"
    }
  ]
}
```   

```go
package main

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"os"
)

type Post struct {
	Id  int `json:"id"`
	Content string `json:"content"`
	Author Author `json:"author"`
	Comments []Comment `json:"comments"`
}

type Author struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type Comment struct {
	Id int `json:"id"`
	Content string `json:"content"`
	Author string `json:"author"`
}

func main() {
	jsonFile, err := os.Open("post.json")
	if err != nil {
		fmt.Println("There is an error on Opending File:", err)
		return
	}

	defer jsonFile.Close()
	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading JSON Data:", err)
		return
	}

	var post Post
	json.Unmarshal(jsonData, &post)
	fmt.Println(post)
}
```

使用解码器的方案：   

```go
jsonFile, err := os.Open("post.json")
if err != nil {
	fmt.Println("Error opening JSON file:", err)
	return
}
defer jsonFile.Close()
decoder := json.NewDecoder(jsonFile)
for {
	var post Post
	err := decoder.Decode(&post)
	if err == io.EOF {
		break
	}
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	fmt.Println(post)
}
```    

感觉这里的循环没必要啊，又不是多个文件，明明一次就完了，循环个篮子。    

# 第 8 章 应用测试

## 8.2 使用 Go 进行单元测试

```go
// main.go
package main

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"os"
)

type Post struct {
	Id  int `json:"id"`
	Content string `json:"content"`
	Author Author `json:"author"`
	Comments []Comment `json:"comments"`
}

type Author struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type Comment struct {
	Id int `json:"id"`
	Content string `json:"content"`
	Author string `json:"author"`
}

func main() {
	_, err := decode("post.json")
	if err != nil {
		fmt.Println("Error:", err)
	}
}


func decode(filename string) (post Post, err error){
	jsonFile, err := os.Open(filename)
	if err != nil {
		fmt.Println("There is an error on Opending File:", err)
		return
	}

	defer jsonFile.Close()
	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading JSON Data:", err)
		return
	}

	json.Unmarshal(jsonData, &post)
	return
}
```   

json 的内容还是上一章的那些。   

```go
// main_test.go
package main

import (
	"testing"
)

func TestDecode(t *testing.T) {
	post, err := decode("post.json")
	if err != nil {
		t.Error(err)
	}

	if post.Id != 1 {
		t.Error("Wrong id, was expecting 1 but got", post.Id)
	}
	if post.Content != "Hello World!" {
		t.Error("Wrong content, was expecting 'Hello World!' bug got", post.Content)
	}
}

func TestEncode(t *testing.T) {
	t.Skip("Skiping encoding for now")
}
```    

testing.T 结构拥有几个非常有用的函数：   

+ **Log** - 将给定的文本记录到错误日志里面，与 fmt.Println 类似
+ **Logf** - 根据给定的格式，将给定的文本记录到错误日志里面，与 fmt.Printf 类似
+ **Fail** - 将测试函数标记为 “已失败”，但允许测试函数继续执行
+ **FailNow** - 将测试函数标记为 “已失败” 并停止执行测试函数

除以上4个函数之外, testing.T 结构还提供了如下的一些便利函数(convenience	function),这些
便利函数都是由以上4个函数组合而成的。    


&nbsp; | Log | Logf
---------|----------|---------
 Fail | Error | Errorf
 FailNow | Fatal | Fatalf   

无论是 Fail 还是 FailNow ,它们都只会对自己所处的测试用例产生影响,比如,在上面的例子中,
TestDecode 调用的 Error 函数就只会对 TestDecode 本身产生影响。   

测试的时候使用 -v(verbose) 选项获取更详细的信息，通过覆盖率标志 -cover 来获知测试用例对代码
的覆盖率： `go test -v -cover`   

### 8.2.1 跳过测试用例

除了可以使用 Skip 函数直接跳过整个测试用例，用户还可以通过向 go test 命令传入短暂标志 -short,
并在测试用例中使用某些条件逻辑来跳过测试中的指定部分。注意,这种做法跟在 go test 命令中通过选项
来选择性地执行指定的测试不一样:选择性执行只会执行指定的测试,并跳过其他所有测试,而 -short 标志
则会根据用户编写测试代码的方式,跳过测试中的指定部分或者跳过整个测试用例。    

```go
func TestLongRunningTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping long-running test in short mode")
	}
	time.Sleep(10*time.Second)
}
```    

### 8.2.2 以并行方式运行测试

```go
func TestParallel_1(t *testing.T) {
	t.Parallel()
	time.Sleep(1 * time.Second)
}

func TestParallel_2(t *testing.T) {
	t.Parallel()
	time.Sleep(2 * time.Second)
}

func TestParallel_3(t *testing.T) {
	t.Parallel()
	time.Sleep(3 * time.Second)
}
```    

现在,我们只要在终端中执行以下命令,Go就会以并行的方式运行测试:   

`go test -v -short -parallel 3`    

### 8.2.3 基准测试

```go
// main_bench_test.go
package main

import (
	"testing"
)

func BenchmarkDecode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		decode("post.json")
	}
}
```   

为了运行基准测试用例,用户需要在执行 go test 命令时使用基准测试标志 -bench ,并将一个正则
表达式用作该标志的参数,从而标识出自己想要运行的基准测试文件。当我们需要运行目录下的所有基准测试
文件时,只需要把点. 用作 -bench 标志的参数即可:   

`go test -v -cover -short -bench .`

![benchmark-test](https://github.com/temple-deng/markdown-images/blob/master/uncategorized/go-benchmark-test.png)   

结果中的 100000 为测试时 b.N 的实际值。需要注意的是,在进行基准测试时,测试用例的迭代次数是由
Go 自行决定的,虽然用户可以通过限制基准测试的运行时间达到限制迭代次数的目的,但用户是无法直接指定
迭代次数的——测试程序将进行足够多次的迭代,直到获得一个准确的测量值为止。在 Go 1.5 中, test
子命令拥有一个 -test.count 标志,它可以让用户指定每个测试以及基准测试的运行次数,该标志的默认值为1。   

注意,上面的命令既运行了基准测试,也运行了功能测试。如果需要,用户也可以通过运行标志 -run 来忽略
功能测试。 -run 标志用于指定需要被执行的功能测试用例,如果用户把一个不存在的功能测试名字用作
-run 标志的参数,那么所有功能测试都将被忽略。比如,如果我们执行以下命令:   

`go test -run x -bench .`    

那么由于我们的测试中不存在任何名字为 x 的功能测试用例,因此所有功能测试都不会被运行。   

## 8.3 使用 Go 进行 HTTP 测试

对 Go Web 应用的单元测试可通过 testing/httptest 包来完成。这个包提供了模拟一个 Web 服务器
所需的设施，用户可以利用 net/http 包中的客户端函数向这个服务器发送 HTTP 请求，然后获取模拟
服务器返回的 HTTP 响应。    

# 第 9 章 发挥 Go 的并发优势

## 9.1 并发与并行的区别

并发指的是两个或多个任务在同一时间段内启动、运行并结束，并且这些任务可能会互动。与并发形式执行的
多个任务会同时存在。    

并行与并发是两个看上去相似但实际上却截然不同的概念，因为并发和并行都可以同时运行多个任务，所以
很多人都把这两个概念混淆了。对于并发来说，多个任务并不需要同时开始或同时结束——这些任务的执行
过程在时间上是相互重叠。并发执行的多个任务会被调度，并且它们会通过通信分享数据并协调执行时间（
不过这种通信不是必须的）。    

在并行中，多个任务将同时启动并执行。并行通常会把一个大任务分割成多个更小的任务，然后通过同时执行
这些小任务来提高性能。并行通常需要独立的资源（如CPU），而并发则会使用和分享相同的资源。因为并行
考虑的是同时启动和执行多个任务，所以它在直觉上更易懂一些。并行，正如它的名字所昭示的那样，是
一系列相互平行、不会重叠的处理过程。    

并发指的是同时处理多项任务，而并行指的是同时执行多项任务。   

尽管并发和并行在概念上并不相同,但它们并不相互排斥,比如Go语言就可以创建出同时具有并发和并行这两种
特征的程序。    

其实并发更像是一组功能相同的任务，然后这组任务可能会在不同的时间启动，但是执行时间上可能会有所
重叠，但并行更像是一组功能不同的任务，但是需要在同一时间启动。并发更像是一种我们模拟出来的多任务
同时执行的一种状态。并行却是实实际际的多个任务在同一时刻都有在执行，所以并行要求我们在机器上就有
多个相同的资源，例如多个 CPU，但是并发却不一定，我们可以在同一资源上执行这多个任务，通过使用合适
的调度方案，来实现一种多任务并发执行的效果。     

## 9.2 goroutine

当一个 goroutine 被阻塞时，它也会阻塞所复用的操作系统线程，而运行时环境则会把位于被阻塞线程上的
其他 goroutine 移动到其他未阻塞线程上继续执行。   

等等啊，这里阻塞系统线程，不就会引起操作系统的线程调度嘛，那运行时环境只有在下次进程内某个线程可以
执行时才能运行 goroutine 的调度吧。    

### 9.2.1 等待 goroutine

Go 语言在 sync 包中提供了一种名为 WaitGroup 的机制，它的运作方式非常简单直接：   

+ 声明一个 WaitGroup
+ 使用 Add 方法为 WaitGroup 的计数器设置值
+ 当一个 goroutine 完成它的工作时，使用 Done 方法对 WaitGroup 的计数器执行减一操作
+ 调用 Wait 方法，该方法将一直阻塞，知道等待组计数器的值变为 0    

# 第 10 章 Go 的部署

## 10.1 将应用部署到独立服务器

当我们在服务器上启动程序后，此时 Web 服务是在前台运行，所以在服务运行期间，我们将无法执行其他操作。
于此同时，我们也无法简单地通过 `&` 或者 `bg` 命令在后台运行这个服务，因为这样做的话，一旦用户
登出，Web 服务就会被杀死。   

避免上述问题的一种方法就是使用 `nohup` 命令,让操作系统在用户注销时,把发送至Web服务的 HUP
(hangup,挂起)信号忽略掉:    

`nohup ./ws-s &`   

除 nohup 之外,持续运行 Web 服务的另一种方法是使用Upstart或者systemd这样的 init 守护进程:
init 进程是类Unix系统在启动时运行的第一个进程,该进程由内核负责启动,它会一直运行直到系统关闭为止,
并且它还是其他所有进程直接或间接的祖先。    

Upstart是由Ubuntu创建的一个基于事件的 init 替代品。为了使用Upstart,用户首先需要创建一个对应
的 Upstart 任务配置文件,并将该文件放到 etc/init 目录里面。    

```shell
# ws.conf
respawn
respawn limit 10 5

setuid sausheong
setgid sausheong

exec /go/src/github.com/sausheong/ws-s/ws-s
```    

这个 Upstart 任务配置文件非常简单和直接。文件中的每个 Upstart 任务都由一个或任意多个称为节
(stanzas)的命令块组成。第一节 respawn 指示当任务失效(fail)时,Upstart将对其实施重新派生
(respawn)或者重新启动。第二节 respawn limit 10 5 为 respawn 设置了参数,它指示Upstart
最多只会尝试重新派生该任务10次,并且每次尝试之间会有 5s 的间隔;在用完了10次重新派生的机会之后,
Upstart将不再尝试重新派生该任务,并将该任务视为已失效。第三节和第四节负责设置运行进程的
用户以及用户组,而最后一节则是Upstart在启动任务时需要运行的可执行文件。    

为了启动上述Upstart任务,我们需要在终端里面执行以下命令:   

```shell
$ sudo start ws
ws start/running, process 2011
```   

这个命令将触发Upstart读取 /etc/init/ws.conf 任务配置文件并启动任务。    

