# Package os

## Overview

os 包提供了一组用来访问操作系统功能的接口，这些接口是平台无关的。os 包的接口是被设计为跨平台的。    

## 常量

`OpenFile` 的标志，这些标志是包装了操作系统对应的标志产生的。并非每个操作系统对实现了所有的标志。   

```go
const (
  O_RDONLY int = syscall.O_RDONLY    // 只读形式打开
  O_WRONLY int = syscall.O_WRONLY    // 只写形式打开
  O_RDWR   int = syscall.O_RDWR      // 读写形式开始
  O_APPEND int = syscall.O_APPEND    // 写入时采用追加数据的形式
  O_CREATE int = syscall.O_CREAT     // 如果打开时文件不存在，创建一个新文件
  O_EXCL   int = syscall.O_EXCL      // 与 O_CREATE 一起使用，文件必须是不存在的
  O_SYNC   int = syscall.O_SYNC      // 以同步 I/O 的形式打开
  O_TRUNC  int = syscall.O_TRUNC     // 如果可能，当打开文件时截取文件
)
```    

路径分隔符：   

```go
const (
  PathSeparator     = '/'   // 系统特定的路径分隔符
  PathListSeparator = ':'   // 系统特定的路径列表分隔符
)
```    

`DevNull` 是系统 null 设备的名字，在 Unix 类系统上，即 /dev/null，在 Windows 上就是 NUL.    

`const DevNull = "/dev/null"`    

## 变量

一些系统调用错误的别名：   

```go
var (
    ErrInvalid    = errors.New("invalid argument") // methods on File will return this error when the receiver is nil
    ErrPermission = errors.New("permission denied")
    ErrExist      = errors.New("file already exists")
    ErrNotExist   = errors.New("file does not exist")
    ErrClosed     = errors.New("file already closed")
)
```    

`Stdin, Stdout, Stderr` 是打开的文件的，这些文件指向了标准输入、输出、错误的文件描述符。   

```go
var (
    Stdin  = NewFile(uintptr(syscall.Stdin), "/dev/stdin")
    Stdout = NewFile(uintptr(syscall.Stdout), "/dev/stdout")
    Stderr = NewFile(uintptr(syscall.Stderr), "/dev/stderr")
)
```    

Args 保存了命令行参数，第一个元素是被调用程序的名字：   

`var Args []string`    

## 函数和类型

### func Chdir

`func Chdir(dir string) error`    

修改 CWD，如果出错了，错误会是一个 `*PathError` 类型的。    

### fun Chmod

`func Chmod(name string, mode FileMode) error`   

修改指定文件的模式。如果文件是一个软连接文件，修改连接目标的模式，错误为 `*PathError`。   

不同的系统上有不同的模式 bit。   

Unix 上包括权限位, ModeSetuid, ModeSetgid, ModeSticky。   

Windows 上看不懂。    

### func Chown

`func Chown(name string, uid, gid int) error`    

如果文件是软链接，修改链接目标的 uid 和 gid，错误为 `*PathError` 类型。   

### func Chtimes

`func Chtimes(name string, atime time.Time, mtime time.Time) error`    

修改访问和修改时间。错误为 `*PathError` 类型。   

### func Clearenv

`func Clearenv()`    

删除所有环境变量。    

### func Environ

`func Environ() []string`    

以 `key=value` 的字符串形式返回所有环境变量的副本。    

### func Executable

`func Executable() (string, error)`    

返回启动当前进程的可执行文件的路径名。无法保证路径仍指向正确的可执行文件。如果使用符号链接启动进程，则根据操作系统，结果可能是符号链接或它指向的路径。如果需要稳定的结果，path/filepath.EvalSymlinks可能会有所帮助。    

可能返回的是一个类似这样的字符串 `/tmp/go-build648073082/command-line-arguments/_obj/exe/os`。    

返回的是绝对路径。    

### func Exit

`func Exit(code int)`    

程序会立刻终止，dederred 函数不会继续运行。   

### func Expand

`func Expand(s string, mapping func(string) string) string`    

根据映射函数将字符串中的 $var 或 ${var} 替换掉。个人感觉这里字符串应该是一个包含类似环境变量
变量名的字符串，然后映射函数会将其中的变量名映射为变量值，然后返回整个替换后的字符串。   

### func ExpandEnv

`func ExpandEnv(s string) string`    

根据当前环境变量替换值，如果引用了未定义的变量，则映射为空字符串。`os.ExpandEnv(s)` 等价于
`os.Expand(s, os.GetEnv)`。     

`fmt.Println(os.ExpandEnv("$USER lives in ${HOME}"))`     

注意 `ExpandEnv` 函数看样子是只能进行环境变量的替换，但是 `Expand` 函数却可以进行普通的替换：   

```go
s2 := "$user lives in $home, ${hh}"
fmt.Println(os.Expand(s2, mapping))    // dengbo111 lives in 37-5-402, hhhhhh

func mapping(origin string) string {
  switch origin {
  case "user":
    return "dengbo111"
  case "home":
    return "37-5-402"
  default:
    return "hhhhhh"
  }
}
```   

### func Getegid

`func Getegid() int`    

返回调用者的组 id。好像是这个意思，在 Windows 上，返回 -1.    

### func Getenv()

`func Getenv(key string) string`    

如果对应的环境变量不存在，返回空字符串。如果想要区分一个空字符串还是未设置的变量，使用 `LookupEnv`。   

### func Geteuid

`func Geteuid() int`   

返回调用者的 uid。在 Windows 上返回 -1.    

e 代表 effective，应该是意指当前激活的那个 id，比如用户是可以有多个用户组的，但 Getegid 可能
只返回当前用户的激活组。    

### func Getgroups

`func Getgroups() ([]int, error)`    

返回调用者的所有所属组 id。在 Windows 上返回 syscall.EWINDOWS。    

### func Getpagesize

`func Getpagesize() int`    

返回系统内存页的大小。好像是以 B 为单位的。    

### func Getpid()

`func Getpid() int`   

### func Getuid()

`func Getuid() int`    

Windows 上返回 -1.   

### func Getwd()

`func Getwd() (dir string, err error)`   

以绝对路径的形式返回工作目录，如果可以通过多个路径访问当前工作目录（例如软链接），则随机返回其中
的一个。话说什么情况会出错呢。    

### func Hostname

`func Hostname() (name string, err error)`    

返回由内核决定的主机名。   

### func IsExist

`func IsExist(err error) bool`   

IsExist返回一个布尔值，指示是否已知错误报告文件或目录已存在。 ErrExist满足它以及一些系统调用错误。    

好像是这个意思，这个布尔值用来只是我们之前一些文件访问函数返回的错误，是否表明了那个文件或目录是
已经存在的。    

### func IsNotExist

`func IsNotExist(err error) bool`    

例子：   

```go
filename := "a-nonexistent-file"

if _, err := os.Stat(filename); os.IsNotExist(err) {
  fmt.Println("file does not exist")
}
```   

### func IsPathSeparator

`func IsPathSeparator(c uint8) bool`   

c 是否是目录分隔字符。    

### func IsPermission

`func IsPermission(err error) bool`    

是否错误是权限问题产生的。    

### func Lchown

`func Lchown(name string, uid, gid int) error`    

如果文件是软链接，修改链接自身。错误为 `*PathError` 类型。   

在 Windows 上总返回用 `*PathError` 包裹后的 syscall.EWINDOWS 错误。    

### func Link

`func Link(oldname, newname string) error`   

创建硬链接文件，错误为 `*LinkError`。    

### func LookupEnv

`func LookupEnv(key string) (string, bool)`    

如果环境变量存在，返回其值及布尔值 true，否则返回空字符串及布尔值 false。    

### func Mkdir

`func Mkdir(name string, perm FileMode) error`   

错误为 `*PathError` 类型。   

### func MkdirAll

`func MkdirAll(path string, perm FileMode) error`   

递归创建。    

### func NewSyscallError

`func NewSyscallError(syscall string, err error) error`    

这个好像类似 `errors.New` 用来生成错误的，使用给定的系统调用名和错误细节，产生一个 `SyscallError`
并返回。   

### func Readlink

`func Readlink(name string) (string, error)`    

应该的意思是返回软链接文件的链接目的地的路径吧。错误为 `*PathError` 类型。   

### func Remove

`func Remove(name string) error`   

删除文件或目录，错误为 `*PathError` 类型。   

### func RemoveAll

`func RemoveAll(path string) error`   

移除指定的路径及其包含的子路径。    

### func Rename

`func Rename(oldpath, newpath string) error`    

如果 newpath 已存在且不是一个目录，直接替换。错误为 `*LinkError` 类型。   

### func SameFile

`func SameFile(fi1, fi2 FileInfo) bool`   

### func Setenv

`func Setenv(key, value string) error`   

### func Symlink

`func Symlink(oldname, newname string) error`   

错误为 `*LinkError`。   

### func TempDir

`func TempDir() string`   

返回当前系统上存放临时文件的默认目录。   

在 Unix 上，返回 `$TMPDIR` 或者 `/tmp`，在 Windows 上，使用 `GetTempPath`，返回
`%TMP%, %TEMP, %USERPROFILE%` 中的第一个非空值，或者二是 Windows 目录。    

### func Truncate

`func Truncate(name string, size int64) error`    

修改文件的大小，如果是软链接文件，修改其目标的大小。错误为 `*PathError`。   

### func Unsetenv

`func Unsetenv(key string) error`   

## type File

File 代表了一个打开文件的描述符。    

### func Create

`func Create(name string) (*File, error)`    

使用 0666 的模式创建一个文件。如果文件存在的话就截断 truncate（什么意思）。如果成功的话，返回
其文件描述符，模式为 O_RDWR，注意这里是文件描述符的模式。错误为 `*PathError`。   

### func NewFile

`func NewFile(fd uintptr, name string) *File`    

返回一个新文件，其文件描述符为 fd，如果 fd 不是一个有效的描述符的话，返回 nil。   

### func Open

`func Open(name string) (*File, error)`    

以读取模式打开文件。即默认的文件描述符模式为 O_RDONLY。错误为 `*PathError`。    

### func OpenFile

`func OpenFile(name string, flag int, perm FileMode) (*File)`   

```go
f, err := os.OpenFile("notes.txt", os.O_RDWR | os.O_CREATE, 0755)
if err != nil {
  log.Fatal(err)
}
if err := f.Close(); err != nil {
  log.Fatal(err)
}
```    

### func Pipe

`func Pipe() (r *File, w *File, err error)`    

应该是用管道链接两个文件，从 r 读取内容写到 w 中。    

### func (*File) Chdir

`func (f *File) Chdir() error`

将 CWD 修改为指定的文件，因为 File 必须是一个目录。   

### func (*File) Chmod

`func (f *File) Chmod(mode FileMode) error`    

### func (*File) Chown

`func (f *File) Chown(uid, gid int) error`   

### func (*File) Close

`func (f *File) Close() error`    

### func (*File) Fd  

`func (f *File) Fd() uintptr`    

返回引用了被打开文件的整型的 Unix 文件描述符。   

### func (*File) Name

`func (f *File) Name() string`   

返回文件名。   

### func (*File) Read

`func (f *File) Read(b []byte) (n int, err error)`    

从文件中最多读取 len(b) 字节。在文件末尾，返回 0, io.EOF.   

### func (*File) ReadAt

`func (f *File) ReadAt(b []byte, off int64) (n int, err error)`   

从 off 字节偏移处读取至多 len(b) 字节。ReadAt 在 n &lt; len(b) 的情况下返回一个非 nil 的
错误。在文件末尾，错误是 io.EOF。    

### func (*File) Readdir

`func (f *File) Readdir(n int) ([]FileInfo, error)`    

返回目录中至多 n 条 FileInfo 值。这 FileInfo 好像是排序过的，如果我们后面再调用这个函数，可能
会返回后面的几条数据。   

如果到达了目录结尾，错误就是 io.EOF。     

```go
	file, err := os.Open(".")
	if err != nil {
		fmt.Println(err)
	}

	infos, err := file.Readdir(2)
	if err != nil {
		fmt.Println(err)
	}
	for _, f := range infos {
		fmt.Println(f.Name())
  }
```    

### func (*File) Readdirnames

`func (f *File) Readdirnames(n int) (names []string, err error)`   

好像于 Readdir 类似，但是直接返回了名字。    

### func (*File) Seek

`func (f *File) Seek(offset int64, whence int) (ret int64, err error)`    

Seek 操作将下次的读写偏移位置设置为 offset 处。具体 offset 的解释方式取决于 whence，0 代表
相对于文件起始位置偏移，1 代表相对当前位置偏移，2 代表相对文件末尾。    

### func (*File) Stat

`func (f *File) Stat() (FileInfo, error)`    

### func (*File) Sync

`func (f *File) Sync() error`    

提交当前文件内容到存储设备中。通常来说，意味着刷新系统内存中的数据到磁盘上。   

### func (*File) Truncate

`func (f *File) Truncate(size int64) error`    

改变文件的大小。   

### func (*File) Write

`func (f *File) Write(b []byte) (n int, err error)`    

写入 len(b) 到文件中。    

### func (*File) WriteAt

`func (f *File) WriteAt(b []byte, off int64) (n int, err error)`   

### func (*File) WriteString

`func (f *File) WriteString(s string) (n int, err error)`   

类似与 `Write` 方法，不过是直接把字符串写入而不是 byte 切片。   

## type FileInfo



