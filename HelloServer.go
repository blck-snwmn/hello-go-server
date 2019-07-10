package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<html><body>hello world</body></html>")
}

func main() {
	var httpServer http.Server
	http.HandleFunc("/", handler)
	log.Println("start http listen :18888")
	httpServer.Addr = ":18888"
	log.Println(httpServer.ListenAndServe())
}
