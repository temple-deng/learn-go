# Go Web 编程 Part1

<!-- TOC -->

- [Go Web 编程 Part1](#go-web-编程-part1)
- [第 1 章 Go 与 Web 应用](#第-1-章-go-与-web-应用)
  - [1.1 使用 Go 语言构建 Web 应用](#11-使用-go-语言构建-web-应用)
    - [1.1.1 可扩展](#111-可扩展)
    - [1.1.2 模块化](#112-模块化)
    - [1.1.3 可维护](#113-可维护)
    - [1.1.4 高性能](#114-高性能)
  - [1.2 Web 应用的诞生](#12-web-应用的诞生)
- [第 2 章 ChitChat 论坛](#第-2-章-chitchat-论坛)
  - [2.1 应用设计](#21-应用设计)
  - [2.4 请求的接收和处理](#24-请求的接收和处理)
    - [2.4.1 多路复用器](#241-多路复用器)
    - [2.4.2 服务静态文件](#242-服务静态文件)
    - [2.4.3 处理器函数](#243-处理器函数)
    - [2.4.4 使用 cookie](#244-使用-cookie)
  - [2.5 使用模板](#25-使用模板)
- [第 3 章 接收请求](#第-3-章-接收请求)
  - [3.1 使用 Go 构建服务器](#31-使用-go-构建服务器)
  - [3.2 通过 HTTPS 提供服务](#32-通过-https-提供服务)
  - [3.3 处理器和处理器函数](#33-处理器和处理器函数)
    - [3.3.1 处理请求](#331-处理请求)
    - [3.3.2 使用多个处理器](#332-使用多个处理器)
    - [3.3.3 处理器函数](#333-处理器函数)
    - [3.3.4 串联多个处理器和处理器函数](#334-串联多个处理器和处理器函数)
    - [3.3.5 ServeMux 和 DefaultServeMux](#335-servemux-和-defaultservemux)
    - [3.3.6 使用其他多路复用器](#336-使用其他多路复用器)
  - [3.4 使用 HTTP/2](#34-使用-http2)
- [第 4 章 处理请求](#第-4-章-处理请求)
  - [4.1 请求和响应](#41-请求和响应)
    - [4.1.1 Request 结构](#411-request-结构)
    - [4.1.2 请求 URL](#412-请求-url)
    - [4.1.3 请求首部](#413-请求首部)
    - [4.1.4 请求主体](#414-请求主体)
  - [4.2 Go 与 HTML 表单](#42-go-与-html-表单)
    - [4.2.1 Form 字段](#421-form-字段)
    - [4.2.2 PostForm 字段](#422-postform-字段)
    - [4.2.3 MultipartForm 字段](#423-multipartform-字段)
    - [4.2.4 文件](#424-文件)
  - [4.3 ResponseWriter](#43-responsewriter)
  - [4.4 Cookie](#44-cookie)
    - [4.4.1 Go 与 cookie](#441-go-与-cookie)
    - [4.4.2 将 cookie 发送至浏览器](#442-将-cookie-发送至浏览器)

<!-- /TOC -->

# 第 1 章 Go 与 Web 应用

## 1.1 使用 Go 语言构建 Web 应用

在开发大规模Web应用方面,Go语言提供了一种不同于现有语言和平台但又切实可行的方案。大规模可扩展的
Web应用通常需要具备以下特质:    

+ 可扩展
+ 模块化
+ 可维护
+ 高性能    

### 1.1.1 可扩展

大规模的Web应用应该是可扩展的(scalable),这意味着应用的管理者应该能够简单、快速地提升应用的性
能以便处理更多请求。如果一个应用是可扩展的,那么它就是线性的,这意味着应用的管理者可以通过添加更
多硬件来获得更强的请求处理能力。    

有两种方式可以对性能进行扩展:    

+ 一种是垂直扩展(vertical	 scaling),即提升单台设备的CPU数量或者性能;
+ 另一种则是水平扩展(horizontal	 scaling),即通过增加计算机的数量来提升性能。     

因为Go语言拥有非常优异的并发编程支持,所以它在垂直扩展方面拥有不俗的表现:一个Go Web应用只需要
使用一个操作系统线程(OS thread),就可以通过调度来高效地运行数十万个goroutine。  

跟其他Web应用一样,Go也可以通过在多个Go Web应用之上架设代理来进行高效的水平扩展。因为Go Web
应用都会被编译为不包含任何动态依赖关系的静态二进制文件,所以我们可以把这些文件分发到没有安装Go语言
的系统里,从而以一种简单且一致的方式部署Go	Web应用。    

### 1.1.2 模块化

尽管Go是一门静态类型语言,但用户可以通过它的接口机制对行为进行描述,以此来实现动态类型匹配
(dynamic typing)。Go语言的函数可以接受接口作为参数,这意味着用户只要实现了接口所需的方法,
就可以在继续使用现有代码的同时向系统中引入新的代码。与此同时,因为 Go 语言的所有类型都实现了
空接口,所以用户只需要创建出一个接受空接口作为参数的函数,就可以把任何类型的值用作该函数的实际参数。
此外,Go语言还实现了一些在函数式编程中非常常见的特性,其中包括函数类型、使用函数作为值以及闭包,这些特性允许用户使用已有的函数来构建新的函数,从而帮助用户构建出更为模块化的代码。    

### 1.1.3 可维护

因为 Go 语言希望文档可以和代码一同演进,所以它的文档工具 godoc 会对 Go 源代码及其注释进行语法
分析,然后以HTML、纯文本或者其他多种格式创建出相应的文档。godoc 的使用方法非常简单,开发者只需要
把文档写到源代码里面,godoc 就会把这些文档以及与之相关联的代码提取出来,生成相应的文档文件。    

### 1.1.4 高性能

略。   

## 1.2 Web 应用的诞生

在万维网出现不久之后,人们开始意识到一点:尽管使用 Web 服务器处理静态 HTML 文件这个主意非常棒,
但如果 HTML 里面能够包含动态生成的内容,那么事情将会变得更加有趣。其中,通用网关接口
(Common Gateway Interface, CGI)就是在早期尝试动态生成 HTML 内容的技术之一。    

1993年,美国国家超级计算应用中心(National Center for Supercomputing Applications, NCSA)
编写了一个在 Web 服务器上调用可执行命令行程序的规范(specification), 他们把这个规范命名为 CGI,
并将它包含在了NCSA开发的广受欢迎的HTTPd服务器里面。不过NCSA制定的这个规范最终并没有成为正式的
互联网标准,只有 CGI 这个名字被后来的规范沿用了下来。    

CGI 是一个简单的接口,它允许 Web 服务器与一个独立运行于 Web 服务器进程之外的进程进行对接。
通过 CGI 与服务器进行对接的程序通常被称为 CGI 程序,这种程序可以使用任何编程语言编写——这也是我们
把这种接口称之为“通用”接口的原因,不过早期的 CGI 程序大多数都是使用Perl语言编写的。向 CGI 程序
传递输入参数是通过设置环境变量来完成的, CGI 程序在运行之后将向标准输出(stand output)返回结果,而
服务器则会将这些结果传送至客户端。    

与 CGI 同期出现的还有服务器端包含(server-side includes, SSI)技术,这种技术允许开发者在 HTML
文件里面包含一些指令(directive): 当客户端请求一个 HTML 文件的时候,服务器在返回这个文件之前,会
先执行文件中包含的指令,并将文件中出现指令的位置替换成这些指令的执行结果。SSI 最常见的用法是在 HTML
文件中包含其他被频繁使用的文件,又或者将整个网站都会出现的页面首部(header)以及尾部(footer)的
代码段嵌入 HTML 文件中。   

SSI 技术的最终演化结果就是在 HTML 里面包含更为复杂的代码,并使用更为强大的解释器(interpreter)。
这一模式衍生出了PHP、ASP、JSP 和 ColdFusion 等一系列非常成功的引擎,开发者通过使用这些引擎能
够开发出各式各样复杂的Web应用。    

# 第 2 章 ChitChat 论坛

## 2.1 应用设计

请求的格式通常是由应用自行决定的, 比如, ChitChat 的请求使用的是以下格式:
http://<服务器名><处理器名>?<参数> 。    

服务器名(server name)是ChitChat服务器的名字,而处理器名(handler name)则是被调用的处理器的名字。
处理器的名字是按层级进行划分的:位于名字最开头是被调用模块的名字,而之后跟着的则是被调用子模块的名字,
以此类推,位于处理器名字最末尾的则是子模块中负责处理请求的处理器。    

当请求到达服务器时,多路复用器(multiplexer)会对请求进行检查,并将请求重定向至正确的处理器进行
处理。处理器在接收到多路复用器转发的请求之后,会从请求中取出相应的信息,并根据这些信息对请求进行
处理。在请求处理完毕之后,处理器会将所得的数据传递给模板引擎,而模板引擎则会根据这些数据生成将要
返回给客户端的 HTML。    

其实这里多路复用器有点路由器角色的样子啊。    

## 2.4 请求的接收和处理

### 2.4.1 多路复用器

```go
package main

import (
  "net/http"
)

func main() {

  mux := http.NewServeMux()
  files := http.FileServer(http.Dir("/public"))
  mux.Handle("/static/", http.StripPrefix("/static/", files))

  // 注意这里代码应该是简化过的，所以这里的 index 具体的定义应该是被省略了
  mux.HandleFunc("/", index)

  server := &http.Server{
    Addr: "0.0.0.0:8080",
    Handler: mux,
  }

  server.ListenAndServe()
}
```     

net/http 标准库提供了一个默认的多路复用器，这个多路复用器可以通过调用 `NewServeMux` 函数来
创建：   

`mux := http.NewServeMux`    

为了将发送至根 URL 的请求重定向到处理器，程序使用了 `HandleFunc` 函数：`mux.HandleFunc("/", index)`。   

`HandleFunc` 函数接受一个 URL 和一个处理器的名字作为参数，并将针对给定 URL 的请求转发至指定
的处理器进行处理，因此对上述调用来说，当有针对根 URL 的请求到达时，该请求会被重定向到名为 index
的处理器函数。此外，因为所有的处理器都接受一个 `ResponseWriter` 和一个指向 `Request` 结构的
指针作为参数，并且所有请求参数都可以访问 Request 结构得到，所以程序并不需要向处理器显示地传入
任何请求参数。   

### 2.4.2 服务静态文件

除了负责将请求重定向到相应的处理器，多路复用器还需要为静态文件提供服务。为了做到这一点，程序使用
`FileServer` 函数创建了一个能够为指定目录中的静态文件服务器的处理器，并将这个处理器传递给了
多路复用器的 `Handle` 函数。除此之外，程序还使用 `StripPrefix` 函数去移除请求 URL 中的
指定前缀：    

```go
files := http.FileServer(http.Dir("/public"))
mux.Handle("/static/", http.StripPrefix("/static/", files))
```    

### 2.4.3 处理器函数

```go
func index(w http.ResponseWriter, r *http.Request) {
  files := []string{"templates/layout.html",
                  "templates/navbar.html",
                  "templates/index.html"}
  
  templates := template.Must(template.ParseFiles(files...))
  threads, err := data.Threads();
  if err == nil {
    templates.ExecuteTemplate(w, "layout", threads)
  }
}
```    

看这个函数的意思是，`data.Threads` 读出所有帖子的数据，然后把数据注入到 layout.html 的模板
中把，然后写入到 http.ResponseWriter 里返回。   

### 2.4.4 使用 cookie

```go
func authenticate(w http.ResponseWriter, r *http.Request) {
  r.ParseForm()
  user, _ := data.UserByEmail(r.PostFormValue("email"))

  if user.Password == data.Encrypt(r.PostFormValue("password")) {
    session := user.CreateSession()
    cookie := http.Cookie{
      Name: "_cookie",
      Value: session.Uuid,
      HttpOnly: true,
    }
    http.SetCookie(w, &cookie)
    http.Redirect(w, r, "/", 302)
  } else {
    http.Redirect(w, r, "/login", 302)
  }
}
```     

## 2.5 使用模板

程序调用 `ParseFiles` 函数对这些模板文件进行语法分析，并创建出响应的模板。为了捕获语法分析过程
中可能会产生的错误，程序使用了 `Must` 函数区包裹 `ParseFiles` 函数的执行结果，这样当 `ParseFiles`
返回错误的时候，`Must` 函数就会向用户返回相应的错误报告。   

每个模板文件都定义了一个模板，但这种做法并不是强制的，用户也可以在一个目标文件里面定义多个模板。
我们在模板中使用了 `define` 动作，这个动作通过文件开头的 `{{ define "layout" }}` 和文件
末尾的 `{{ end }}`，把被包围的文本块定义成了 layout 模板的一部分。   

```html
<!-- layout.html -->
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
  {{ template "navbar" . }}
  <div class="container">
    {{ template "content" . }}
  </div>
</body>
</html>
{{ end }}
```    

除了 define 动作之外, layout.html 模板文件里面还包含了两个用于引用其他模板文件的 template
动作。跟在被引用模板名字之后的点( . )代表了传递给被引用模板的数据,比如 {{ template "navbar" . }}
语句除了会在语句出现的位置引入 navbar 模板之外,还会将传递给 layout 模板的数据传递给 navbar 模板。   

程序通过调用 ExecuteTemplate 函数,执行(execute)已经经过语法分析的 layout 模板。
执行模板意味着把模板文件中的内容和来自其他渠道的数据进行合并,然后生成最终的HTML内容。   

# 第 3 章 接收请求

## 3.1 使用 Go 构建服务器

创建一个服务器的步骤非常简单，只要调用 `ListenAndServe` 并传入网络地址以及负责处理请求的处理器
作为参数就可以了。如果网络地址参数为空字符串,那么服务器默认使用80端口进行网络连接;如果处理器参数
为 nil ,那么服务器将使用默认的多路复用器 `DefaultServeMux`。    

```go
package main

import (
  "net/http"
)

func main() {
  http.ListenAndServe("", nil)
}
```   

用户除了可以通过 `ListenAndServe` 的参数对服务器的网络地址和处理器进行配置之外,还可以通过
`Server` 结构对服务器进行更详细的配置,其中包括为请求读取操作设置超时时间、为响应写入操作设置超时
时间以及为 `Server` 结构设置错误日志记录器等。    

```go
package main
import (
  "net/http"
)

func main() {
  server := http.Server{
    Addr: "127.0.0.1:8080",
    Handler: nil,
  }
  server.ListenAndServe()
}
```   

下面展示了 `Server` 结构所有可选的配置选项：    

```go
type Server struct {
  Addr string
  Handler Handler
  ReadTimeout time.Duration
  WriteTimeout time.Duration
  MaxHeaderBytes int
  TLSConfig *tls.Config
  TLSNextProto map[string]func(*Server, *tls.Conn, Handler)
  ConnState func(net.Conn, ConnState)
  ErrorLog *log.Logger
}
```    

## 3.2 通过 HTTPS 提供服务

```go
package main

import (
  "net/http"
)

func main() {
  server := http.Server{
    Addr: "127.0.0.1:8080",
    Handler: nil,
  }
  server.ListenAndServeTLS("cert.pem", "key.pem")
}
```    

cert.pem 文件是 SSL 证书，而 key.pem 则是服务器的私钥。生成证书的办法有很多种，其中一种就是
使用 Go 标准库中的 crypto 包。    

```go
package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"time"
)

func main() {
	max := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, _ := rand.Int(rand.Reader, max)
	subject := pkix.Name{
		Organization: []string{"Manning Publications Co."},
		OrganizationalUnit: []string{"Books"},
		CommonName: "Go Web Programming",
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: subject,
		NotBefore: time.Now(),
		NotAfter: time.Now().Add(365 * 24 * time.Hour),
		KeyUsage: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses: []net.IP{net.ParseIP("127.0.0.1")},
	}

	pk, _ := rsa.GenerateKey(rand.Reader, 2048)

	derBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, &pk.PublicKey, pk)
	certOut, _ := os.Create("cert.pem")
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	keyOut, _ := os.Create("key.pem")
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVIATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
	keyOut.Close()
}
```   

生成SSL证书和密钥的步骤并不是特别复杂。因为SSL证书实际上就是一个将扩展密钥用法(extended key usage)
设置成了服务器身份验证操作的X.509证书。    

但是这里生成的证书好像不能用啊，好像是私钥的格式有误，暂不清楚是哪部分的问题。   

首先,程序使用一个 `Certificate` 结构来对证书进行配置：    

```go
template := x509.Certificate{
  SerialNumber: serialNumber,
  Subject: subject,
  NotBefore: time.Now(),
  NotAfter: time.Now().Add(365 * 24 * time.Hour),
  KeyUsage: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
  ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
  IPAddresses: []net.IP{net.ParseIP("127.0.0.1")},
}
```    

结构中的证书序列号( SerialNumber )用于记录由CA分发的唯一号码,为了能让我们的Web应用运行起来,
程序在这里生成了一个非常长的随机整数来作为证书序列号。而结构中 KeyUsage 字段和 ExtKeyUsage
字段的值则表明了这个X.509证书是用于进行服务器身份验证操作的。最后,程序将证书设置成了只能在IP
地址127.0.0.1之上运行。     

在此之后,程序通过调用 crypto/rsa 标准库中的 `GenerateKey` 函数生成了一个RSA私钥:    

`pk, _ := rsa.GenerateKey(rand.Reader, 2048)`    

程序创建的RSA私钥的结构里面包含了一个能够公开访问的公钥(public key),这个公钥在使用
`x509.CreateCertificate` 函数创建SSL证书的时候就会用到:    

`derBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, &pk.PublicKey, pk)`     

CreateCertificate 函数接受 Certificate 结构、公钥和私钥等多个参数,创建出一个经过DER编码
格式化的字节切片。后续代码的意图也非常简单明了,它们首先使用 encoding/pem 标准库将证书编码
到 cert.pem 文件里面:     

```go
certOut, _ := os.Create("cert.pem")
pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
certOut.Close()
```    

然后继续以PEM编码的方式把之前生成的密钥编码并保存到 key.pem 文件 里面:    

```go
keyOut, _ := os.Create("key.pem")
pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVIATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(pk)})
keyOut.Close()
```    

## 3.3 处理器和处理器函数

### 3.3.1 处理请求

在 Go 语言中，一个处理器就是一个拥有 `ServeHTTP` 方法的接口，这个 `ServeHTTP` 方法需要接受
两个参数：第一个参数是一个 `ResponseWriter` 接口，而第二个参数则是一个指向 `Request` 结构
的指针。    

`DefaultServeMux` 既是 `ServeMux` 结构的实例，也是 `Handler` 结构的实例，因此 `DefaultServeMux`
不仅是一个多路复用器，它还是一个处理器。不过 `DefaultServeMux` 处理器和其他一般的处理器不同，
`DefaultServeMux` 是一个特殊的服务器，它唯一要做的就是根据请求的 URL 将请求重定向到不同的
处理器。     

```go
package main

import (
	"fmt"
	"net/http"
)

type MyHandler struct {}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func main() {
	handler := MyHandler{}
	server := http.Server{
		Addr: "127.0.0.1:8080",
		Handler: &handler,
	}

	server.ListenAndServe()
}
```   

### 3.3.2 使用多个处理器

我们希望使用多个处理器去处理不同的 URL。为了做到这一点,我们不再在 `Server` 结构的 `Handler`
字段中指定处理器,而是让服务器使用默认的 `DefaultServeMux` 作为处理器,然后通过 `http.Handle`
函数将处理器绑定至 `DefaultServeMux`。需要注意的是,虽然 `Handle` 函数来源于 http 包,但它实际上是 `ServeMux`结构的方法:这些函数是为了操作便利而创建的函数,调用它们等同于调用 `DefaultServeMux`
的某个方法。     

```go
package main

import (
	"fmt"
	"net/http"
)

type HelloHandler struct {}

func (h *HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello!")
}

type WorldHandler struct {}

func (h *WorldHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "World!")
}

func main() {
	hello := HelloHandler{}
	world := WorldHandler{}
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}

	http.Handle("/hello", &hello)
	http.Handle("/world", &world)

	server.ListenAndServe()
}
```     

### 3.3.3 处理器函数

处理器函数实际上就是与处理器拥有相同行为的函数：这些函数与 `ServeHTTP` 方法拥有相同的签名。    

```go
package main

import (
	"fmt"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello!")
}

func world(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "World!")
}

func main() {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}

	http.HandleFunc("/hello", hello)
	http.HandleFunc("/world", world)

	server.ListenAndServe()
}
```    

处理器函数的实现原理是这样的:Go语言拥有一种 HandlerFunc 函数类型,它可以把一个带有正确签名的
函数 f 转换成一个带有方法 f 的 `Handler`。换句话说,处理器函数只不过是创建处理器的一种便利的方法而已。   

### 3.3.4 串联多个处理器和处理器函数

诸如日志记录、安全检查和错误处理这样的操作通常被称为横切关注点(cross-cutting concern),虽然
这些操作非常常见,但是为了防止代码重复和代码依赖问题,我们又不希望这些操作和正常的代码搅和在
一起。为此,我们可以使用串联(chaining)技术分隔代码中的横切关注点。     

```go
package main

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")
}

func log(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
		fmt.Println("Handler function called - " + name)
		h(w, r)
	}
}

func main() {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}

	http.HandleFunc("/hello", log(hello))
	server.ListenAndServe()
}
```    

### 3.3.5 ServeMux 和 DefaultServeMux

`ServeMux` 是一个 HTTP 请求多路复用器，它负责接收 HTTP 请求并根据请求中的 URL 将请求重定向
到正确的处理器。    

正如之前所说,因为 ServeMux 结构也实现了 ServeHTTP 方法,所以它也是一个处理器。当 ServeMux
的 ServeHTTP 方法接收到一个请求的时候,它会在结构的映射里面找出与被请求URL最为匹配的URL,然后调
用与之相对应的处理器的 ServeHTTP 方法。    

DefaultServeMux 实际上是 ServeMux 的一个实例,并且所有引入了 net/http 标准库的程序都可以
使用这个实例。当用户没有为 Server 结构指定处理器时,服务器就会使用 DefaultServeMux 作
为 ServeMux 的默认实例。    

在匹配规则上，如果被绑定的 URL 不是以 / 结尾，那么它只会与完全相同的 URL 匹配；但是如果被绑定
的 URL 以 / 结尾，那么即使请求的 URL 只有前缀部分与被绑定的 URL 相同，`ServeMux` 也会认定
这两个 URL 是匹配的。    

### 3.3.6 使用其他多路复用器

因为创建一个处理器和多路复用器唯一需要做的就是实现 ServeHTTP方法,所以通过自行创建多路复用器来
代替 net/http 包中的 ServeMux 是完全可行的。话说之前并没有提到过多路复用器是实现 ServeHTTP
方法就可以了把。这里介绍一个第三方的多路复用器 - HttpRouter。    

```go
package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", p.ByName("name"))
}

func main() {
	mux := httprouter.New()
	mux.GET("/hello/:name", hello)

	server := http.Server{
		Addr: "127.0.0.1:8080",
		Handler: mux,
	}

	server.ListenAndServe()
}
```    

## 3.4 使用 HTTP/2    

在 1.6 以上版本的 Go 语言中，如果使用 HTTPS 模式启动服务器，那么服务器将默认使用 HTTP/2。    

# 第 4 章 处理请求

## 4.1 请求和响应

### 4.1.1 Request 结构

Request 结构主要由以下部分组成：   

+ URL 字段
+ Header 字段
+ Body 字段
+ Form 字段、PostForm 字段和 MultipartForm 字段    

### 4.1.2 请求 URL

URL 字段用于表示请求行中包含的 URL。这个字段是一个指向 `url.URL` 结构实例的指针。    

```go
type URL struct {
  Scheme      string
  Opaque      string
  User        *Userinfo
  Host        string
  Path        string
  RawPath     string
  ForceQuery  bool
  RawQuery    string
  Fragment    string
}
```    

URL 的一般格式为：   

`scheme://[userinfo@]host/path[?query][#fragment]`    

那些在 scheme 之后不带斜线的 URL 则会被解释为：    

`scheme:opaque[?query][#fragment]`     

虽然通过对 `RawQuery` 字段的值进行语法分析可以获取到键值对格式的查询参数，但直接使用 `Request`
结构的 `Form` 字段来获取这些键值对会更方便一些。    

### 4.1.3 请求首部

请求和响应的首部都使用 `Header` 类型描述,这种类型使用一个映射来表示HTTP首部中的多个键值对。
`Header` 类型拥有4种基本方法,这些方法可以根据给定的键执行添加、删除、获取和设置值等操作。    

一个 `Header` 类型的实例就是一个映射,这个映射的键为字符串,而键的值则是由任意多个字符串组成的切片。
为 `Header` 类型设置首部以及添加首部都是非常简单的,但了解这两种操作之间的区别有助于更好地
理解 Header 类型的构造:在对 Header 执行设置操作时,给定键的值首先会被设置成一个空白的字符串切片,
然后该切片中的第一个元素会被设置成给定的首部值;而在对 Header 执行添加操作时,给定的首部值会被
添加到字符串切片已有元素的后面。    

```go
package main

import (
  "fmt"
  "net/http"
)

func headers(w http.ResponseWriter, r *http.Request) {
  h := r.Header
  fmt.Fprintln(w, h)
}

func main() {
  server := http.Server{
    Addr: "127.0.0.1:8080",
  }
  http.HandleFunc("/headers", headers)
  server.ListenAndServe()
}
```   

如果想要获取的是某个特定的首部: `h := r.Header["Accept-Encoding"]`。输出为 `[gzip, deflate]`。   

除此之外,我们还可以使用以下语句: `h := r.Header.Get("Accept-Encoding")`，输出为
`gzip, deflate`。    

注意以上两条语句之间的区别:直接引用	 Header 将得到一个字符串切片,而在 Header 上调用 Get 
方法将返回字符串形式的首部值,其中多个首部值将使用逗号分隔。    

### 4.1.4 请求主体

请求的首部由 `Request` 结构的 `Body` 字段表示，这个字段同时实现了 `io.Reader` 和 `io.Closer`
两个接口。    

```go
package main

import (
  "fmt"
  "net/http"
)

func body(w http.ResponseWriter, r *http.Request) {
  len := r.ContentLength
  body := make([]byte, len)
  r.Body.Read(body)
  fmt.Fprintln(w, string(body))
}

func main() {
  server := http.Server{
    Addr: "127.0.0.1:8080",
  }

  http.HandleFunc("/body", body)
  server.ListenAndServe()
}
```   

## 4.2 Go 与 HTML 表单

### 4.2.1 Form 字段

通过调用 `Request` 结构提供的方法,用户可以将URL、主体又或者以上两者记录的数据提取到该结构的
`Form`、`PostForm` 和 `MultipartForm` 等字段当中。跟我们平常通过 POST 请求获取到的数据一样,
存储在这些字段里面的数据也是以键值对形式表示的。使用 `Request` 结构的方法获取表单数据的一般步骤是:  

1. 调用 `ParseForm` 方法或者 `ParseMultipartForm` 方法，对请求进行语法分析
2. 根据步骤 1 调用的方法，访问相应的 `Form` 字段、`PostForm` 字段或 `MultipartForm` 字段    

```go
package main

import (
  "fmt"
  "net/http"
)

func process(w http.ResponseWriter, r *http.Request) {
  r.ParseForm()
  fmt.Fprintln(w, r.Form)
}

func main() {
  server := http.Server{
    Addr: "127.0.0.1:8080",
  }
  http.HandleFunc("/process", process)
  server.ListenAndServe()
}
```    

然后提交如下页面的表单：   

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
  <form	action=http://127.0.0.1:8080/process?hello=world&thread=123	
    method="post"	enctype="application/x-www-form-urlencoded">
   <input	type="text"	name="hello"	value="sau sheong" />
   <input	type="text"	name="post"	value="456" />
   <input	type="submit"/>
  </form>
</body>
</html>
```    

打印的内容是 `map[hello:[sau sheong world] post:[456] thread:[123]]`。    

### 4.2.2 PostForm 字段

针对上面例子中出现的 hello 字段，如果一个表单域同时出现表单和 URL 中，那么最终访问字段
`r.Form["hello"]` 得出的结果中，表单中的值总是排在 URL 值的前面。   

如果一个字段同时出现在表单和 URL 中，但是用户只想获取表单中的值，可以访问 `Request.PostForm`
字段。     

同时需要注意的是 `Form` 和 `PostForm` 结构只支持 `application/x-www-form-urlencoded`编码，
所以如果是 `multipart/form-date` 编码的表单，数据是无法通过这两个结构获取到的。   

### 4.2.3 MultipartForm 字段

为了取得 multipart/form-data 编码的表单数据，我们需要用到 `Request` 结构的 `ParseMultipartForm`
方法和 `MultipartForm` 字段，而不再使用 `ParseForm` 方法和 `Form` 字段，不过 `ParseMultipartForm`
在需要时也会自行调用 `ParseForm` 方法。    

```go
r.ParseMultipartForm(1024)
fmt.Fprintln(w, r.MultipartForm)
```    

第一行代码说明了我们想要从 multipart 编码的表单里取出的字节量，第二行则打印 MultipartForm 字段。    

`&{map[hello:[sau sheong] post:[456]] map[]}`    

`MultipartForm` 字段只包含表单字段，另外 `MultipartForm` 字段的值也不再是一个映射，而是一个
包含了两个映射的结构，其中第一个映射的键为字符串，值是字符串组成的切片，而第二个映射这里暂时为空，
因为它是用来记录用户上传的文件的。     

除了上面提到的几个方法之外, Request 结构还提供了另外一些方法,它们可以让用户更容易地获取表单中
的键值对。其中, `FormValue` 方法允许直接访问与给定键相关联的值,就像访问 `Form` 字段中的键值对一
样,唯一的区别在于:因为 `FormValue` 方法在需要时会自动调用 `ParseForm` 方法或者 `ParseMultipartForm`
方法,所以用户在执行 `FormValue` 方法之前,不需要手动调用上面提到的两个语法分析方法。   

`FormValue` 方法即使在给定键拥有多个值的情况下,只会从 `Form` 结构中取出给定键的第一个值。也就是
说相当于 Form 中对应字段的切片的第一个元素把。    

除了访问的是 `PostForm` 字段而不是 `Form` 字段之外, `PostFormValue` 方法的作用跟上面介绍的
`FormValue` 方法的作用基本相同。   

这里也有一个需要注意的地方: 如果你将表单的 enctype 设置成了 multipart/form-data ,然后尝试
通过 `FormValue` 方法或者 `PostFormValue` 方法来获取键的值,那么即使这两个方法调用了
`ParseMultipartForm` 方法,你也不会得到任何结果，不过还是可以获取到 URL 中的键值对。   

原因在于 `FormValue` 和 `PostFormValue` 方法是从解析后的 `Form` 和 `PostForm` 中取值的，
那么如果表单使用 multipart/form-data 编码的时候的，解析后的数据是在 `MultipartForm` 字段中。   

**问题：**那这里是怎么决定使用 `FormValue` 或 `PostFormValue` 的时候，是使用 `ParseForm`
还是 `ParseMultipartForm` 方法来解析呢。   

### 4.2.4 文件

```go
package main

import (
	"io/ioutil"
	"log"
	"fmt"
	"net/http"
)

type FileHandler struct {}

func (f *FileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(1024)
	fileHeader := r.MultipartForm.File["upload"][0]
	file, err := fileHeader.Open()
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(w, string(data))
}

func main() {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}
	file := FileHandler{}
	http.Handle("/process", &file)
	server.ListenAndServe()
}
```   

跟 `FormValue` 方法和 `PostFormValue` 方法类似, net/http 库也提供了一个 `FormFile` 方法,
它可以快速地获取被上传的文件: `FormFile` 方法在被调用时将返回给定键的第一个值。   

## 4.3 ResponseWriter

`ResponseWriter` 是一个接口,处理器可以通过这个接口创建 HTTP 响应。`ResponseWriter` 在创建
响应时会用到 `http.response` 结构,因为该结构是一个非导出(nonexported)的结构,所以用户只能通过
`ResponseWriter` 来使用这个结构,而不能直接使用它。   

`ResponseWriter` 接口拥有以下 3 个方法：   

+ Write
+ WriteHeader
+ Header

`Write` 方法接受一个字节数组作为参数，并将数组中的字节写入 HTTP 响应的主体中。如果用户在使用
`Write` 方法执行写入操作的时候，没有为首部设置相应类型的内容类型，那么响应的内容类型将通过
检测被写入的前 512 字节决定。（牛逼啊）   

`WriteHeader` 方法的名字带有一点误导性质，它并不能用于设置响应的首部。`WriteHeader` 方法
接受一个代表 HTTP 响应状态码的整数作为参数，并将这个整数用作 HTTP 响应的返回状态码。在调用这个
方法之后,用户可以继续对 `ResponseWriter` 进行写入,但是不能对响应的首部做任何写入操作。
如果用户在调用 `Write` 方法之前没有执行过 `WriteHeader` 方法,那么程序默认会使用200 OK作
为响应的状态码。   

## 4.4 Cookie

### 4.4.1 Go 与 cookie

```go
type Cookie struct {
  Name        string
  Value       string
  Path        string
  Domain      string
  Expires     time.Time
  RawExpires  string
  MaxAge      int
  Secure      bool
  HttpOnly    bool
  Raw         string
  Unparsed    []string
}
```   

### 4.4.2 将 cookie 发送至浏览器

Cookie 结构的 `String` 方法可以返回一个经过序列化处理的 cookie，其中 Set-Cookie 响应首部
的值就是由这些序列化之后的 cookie 组成的。    

```go
import (
	"net/http"
)

func setCookie(w http.ResponseWriter, r *http.Request) {
	c1 := http.Cookie{
		Name:    "first_cookie",
		Value:	 "Go Web Programming",
		HttpOnly: true,
	}

	c2 := http.Cookie{
		Name:			"second_cookie",
		Value:		"Manning Publications Co",
		HttpOnly: true,
	}

	w.Header().Set("Set-Cookie", c1.String())
	w.Header().Add("Set-Cookie", c2.String())
}

func main() {
	server := http.Server{
		Addr:  "127.0.0.1:8080",
	}

	http.HandleFunc("/set_cookie", setCookie)

	server.ListenAndServe()
}
```    

除了 Set 方法和 Add 方法之外,Go语言还提供了一种更为快捷方便的 cookie 设置方法,那就是使用
net/http 库中的 SetCookie 方法。    

```go
func setCookie(w http.ResponseWriter, r *http.Request) {
  c1 := http.Cookie{
		Name:    "first_cookie",
		Value:	 "Go Web Programming",
		HttpOnly: true,
	}

	c2 := http.Cookie{
		Name:			"second_cookie",
		Value:		"Manning Publications Co",
		HttpOnly: true,
  }
  
  http.SetCookie(w, &c1)
  http.SetCookie(w, &c2)
}
```   

在访问 cookie 上面，可以使用 r.Header["Cookie"] 访问，或者 Request 结构提供了 `Cookie`
方法，可以获取指定名字的 cookie，如果指定的 cookie 不存在，将返回一个错误。如果想要获取所有的
cookie 还可以使用 `Cookies` 方法，Cookies 方法可以返回一个包含了所有cookie的切片,这
个切片跟访问 Header 字段时获取的切片是完全相同的。    


