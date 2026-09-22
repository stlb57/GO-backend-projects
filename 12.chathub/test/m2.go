package main

import "net/http"

type ChatHandler struct {
}

func (ChatHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello"))
}

func main() {
	server := http.Server{
		Addr:    ":8080",
		Handler: ChatHandler{},
	}
	server.ListenAndServe()
}
