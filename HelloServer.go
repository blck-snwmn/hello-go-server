package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

func handler(w http.ResponseWriter, r *http.Request) {
	dump, _ := httputil.DumpRequest(r, true)
	r.ParseForm() // PostForm で値を取得するために必要
	fmt.Println(string(dump))

	fmt.Println("Query", r.URL.Query().Encode())
	fmt.Println("Post Values", r.PostForm.Encode())

	fmt.Fprintf(w, "<html><body>hello world</body></html>")
}

func main() {
	var httpServer http.Server
	http.HandleFunc("/", handler)
	log.Println("start http listen :18888")
	httpServer.Addr = ":18888"
	log.Println(httpServer.ListenAndServe())
}
