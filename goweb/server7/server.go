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