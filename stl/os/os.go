package main

import (
	"os"
	"fmt"
)

func main() {

	// Environ
	environ := os.Environ()
	fmt.Println("-------top 10 environment variables begin-------")
	for i:= 0; i < 10; i++ {
		fmt.Println(environ[i])
	}
	fmt.Println("-------top 10 environment variables end-------")

	executable, err := os.Executable()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(executable)

	s := "$USER lives in ${HOME}"
	fmt.Println(os.Expand(s, os.Getenv))
	fmt.Println(os.ExpandEnv(s))

	s2 := "$user lives in $home, ${hh}"
	fmt.Println(os.Expand(s2, mapping))

	fmt.Printf("The env var USER is: %s\n", os.Getenv("USER"))

	fmt.Println("Current user effective group id is:", os.Getegid())
	fmt.Println("Current user effective user id is:", os.Geteuid())

	groups, _ := os.Getgroups()
	fmt.Println("Current user all groups id is:", groups)
	fmt.Println("Current system memory page size is(unit: B):", os.Getpagesize())

	// 就是当前工作目录，但是奇怪为什么会有错误情况的出现
	wd, _ := os.Getwd()
	fmt.Println("Current working directory:", wd)

	hostname, _ := os.Hostname()
	fmt.Println("Current hostname is:", hostname)

	fileNotExist := "11111.md"
	if _, err := os.Stat(fileNotExist); os.IsNotExist(err) {
		fmt.Println("file does not exist")
	}

	separatorStr := "/\\"
	separator := separatorStr[0]
	winSeparator := separatorStr[1]
	fmt.Println("/ is path separetor:", os.IsPathSeparator(separator))
	fmt.Println("\\ is Linux path separator", os.IsPathSeparator(winSeparator))

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

	file, err = os.Open("../../..")
	if err != nil {
		fmt.Println(err)
	}

	names, err := file.Readdirnames(5)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(names)
}

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