# How to Write Go Code

## Code organization

+ Go 程序员通常将所有的 Go 代码保存在一个工作区 workspace 中
+ 一个工作区内可能包含多个版本控制系统的代码版本，这里的意思应该是工作区中不是会有很多的代码包吗，
很明显我们每个包可能都有一个独立的仓库管理
+ 每个仓库包含一到多个包
+ 每个包在一个目录中包含一个到多个Go 源代码文件
+ 包目录的路径决定了它被导入的路径 import path     

### Workspace

工作区是一个目录层次结构，其根目录包含三个目录：     

+ `src` 包含 Go 源代码文件
+ `pkg` 包含包对象
+ `bin` 包含可执行的命令      

go 工具构建源代码包并将生成的二进制文件安装到 pkg 和 bin 目录.      

`src` 下的子目录通常包含多个版本的仓库，用来追踪一个或多个源代码的变更。     

下面是一个工作区的例子：    

```
bin/
    hello           # command executale
    outyet          # command executable
pkg/
  linux_amd64/
      github.com/golang/example/
          stringutil.a      # package object
src/
    github.com/golang/example/
      .git/             # Git repository metadata
      hello/
        hello.go        # command source
      outyet/
        main.go         # command source
        main_test.go    # test source
      stringutil/
        reverse.go      # package source
        reverse_test.go # test source
    golang.org/x/image/
      .git/
      bmp/
        reader.go       # package source
        writer.go       # package source
```    

上面的工作区中包含两个仓库 example 和 image。example 仓库包含两个命令 hello 和
outyet 和一个库 stringutil。image 仓库包含 bmp 包。    

命令和库由不同类型的源代码包构建而成。     

### The GOPATH environment variable

`GOPATH` 环境变量指定了我们工作区的位置。它通常是在我们用户主目录中的一个 go
目录，所以在 Unix 下就是 $HOME/go，Windows 上是 %USERPROFILE%\go，通常
是 C:\Users\YourName\go。     

命令 `go env GOPATH` 打印出当前有效的 GOPATH，如果没有设置这个环境变量的话，就打印出默认的
位置。     

为了方便，通常会把工作区下的 bin 子目录添加到路径中：    

`$ export PATH=$PATH:$(go env GOPATH)/bin`     

### Import paths

导入路径是一个字符串用来唯一标识一个包。一个包的导入路径对应于它的工作区中的位置
或者是在远程仓库上的位置。    

标准库的中包的导入路径是一个短小的导入路径例如 "fmt", "net/http"。对于我们自己
的包来说，你必须选择一个不太可能与将来添加到标准库或其他外部库冲突的基本路径。    

如果我们的代码是在一些其他的地方，那么我们应该将代码的根目录作为基本路径。     

### Your first program

假设我们的代码是在如下路径中：   

`$ mkdir -p $GOPATH/src/github.com/user/hello`    

新建一个 hello.go 文件。然后使用 go 工具构建并安装这个程序：    

`$ go install github.com/user/hello`     

go 工具会在 GOPATH 指定的工作区中寻找 github.com/user/hello 包。     

如果我们在包中运行这个命令，也可以省略包路径：    

```shell
$ cd $GOPATH/src/github.com/user/hello
$ go install
```     

这条命令会构建 hello 命令，并生成一个二进制文件。他会将二进制文件安装到工作区的
bin 命令下。     

### Your first library   

`$ mkdir $GOPATH/src/github.com/user/stringutil`     

创建一个 reverse.go 文件：    

```go
package stringutil

func Reverse (s string) string {
  r := []rune(s)
  for i, j := 0, len(r) - 1; i < j; i, j = i+1, j-1 {
    r[i], r[j] = r[j], r[i]
  }
  return string(r)
}
```     

使用 `go build` 命令测试包能否成功编译：    

`$ go build github.com/user/stringutil`     

这条命令不会生成输出文件 output file。如果要生成输出文件，必须使用 `go install` 命令，
install 命令会将包对象 package object 放置到工作区的 pkg 目录中。那这条指令只是测试一下
包能否编译通过咯？     

修改 hellp.go 文件：    

```go
package main

import (
  "fmt"
  "github.com/user/stringutil"
)

func main() {
  fmt.Println(stringutil.Reverse("!oG olleH"))
}
```   

无论何时 go 工具安装了一个包或二进制文件，它都会安装其依赖。所以当我们安装hello
时，stringutil 也会自动安装。     

可执行命令的包名必须是 main。    

## Testing

Go 包含一个由 `go test` 命令和 testing 包组成的轻量级测试框架。    

测试文件是以 `_test.go` 结构，包含签名为 `func (t *testing.T)` 的名字为
TestXXX 的函数，测试框架会运行每个函数。如果函数调用了一个失败函数例如 `t.Error`
或者 `t.Fail`，就认为测试没有通过。    

为 stringutil 添加一个测试文件 reverse_test.go：    

```go
package stringutil

import "testing"

func TestReverse(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"Hello, world", "dlrow ,olleH"},
		{"Hello, 世界", "界世 ,olleH"},
		{"", ""},
	}

	for _, c := range cases {
		got := Reverse(c.in)
		if got != c.want {
			t.Errorf("Reverse(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}
```    

使用 `go test` 命令运行测试：    

```shell
$ go test github.com/user/stringutil
ok     github.com/user/stringutil  0.001s
```    

## Remote packages

导入路径可以描述如何使用诸如Git或Mercurial等版本控制系统来获取软件包源代码。go工具使用此属性
自动从远程存储库中获取软件包。上面的例子在 github 上有代码，如果我们导入一个包的仓库 URL 地址，
`go get` 命令会自动获取文件并 **构建** 和 **安装**。    

```shell
$ go get github.com/golang/example/hello
$ $GOPATH/bin/hello
Hello, Go examples!
```     

如果指定的包不在工作区中，`go get`会把包放到工作区中，如果在工作区中，则会跳过 fetch
操作，只运行 install 操作。    


