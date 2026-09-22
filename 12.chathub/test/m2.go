package main

import (
	"fmt"
	"net/http"
)

type ChatHandler struct {
}

func (ChatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
	fmt.Println(r.Method)
	fmt.Println(r.URL.Path)
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/hello/", ChatHandler{})
	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()

}
