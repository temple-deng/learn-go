package main

import (
	"fmt"
	"net/http"
)

type HeaderHandler struct {}

func (h *HeaderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	header := r.Header
	fmt.Fprintln(w, header)
}

type BodyHandler struct {}

func (b *BodyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	len := r.ContentLength
	body := make([]byte,  len)
	r.Body.Read(body)
	fmt.Fprintln(w, string(body))
	defer r.Body.Close()
}

type ProcessHandler struct {}

func (p *ProcessHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Fprintln(w, r.Form)
}

func main() {
	header := HeaderHandler{}
	body := BodyHandler{}
	process := ProcessHandler{}

	mux := http.NewServeMux()
	server := http.Server{
		Addr: "127.0.0.1:8080",
		Handler: mux,
	}

	mux.Handle("/headers", &header)
	mux.Handle("/body", &body)
	mux.Handle("/process", &process)

	server.ListenAndServe()
}