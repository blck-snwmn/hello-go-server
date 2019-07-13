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

	query := r.URL.Query()
	fmt.Println("Query", query)
	fmt.Println("Post Values", r.PostForm)

	name := query.Get("name")

	fmt.Fprintf(w, "<html><body>hello world<br>%s</body></html>", name)
}

func main() {
	var httpServer http.Server
	http.HandleFunc("/", handler)
	log.Println("start http listen :18888")
	httpServer.Addr = ":18888"
	log.Println(httpServer.ListenAndServe())
}
