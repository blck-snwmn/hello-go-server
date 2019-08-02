package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func proxyDummy(w http.ResponseWriter, r *http.Request) {
	//proxy側が実際のURLへ取得しにいかないため、
	//あくまでproxyの設定をしたURLへ通信していることを確認
	url, err := url.Parse("http://127.0.0.1:18888")
	if err != nil {
		log.Fatal(err)
		return
	}
	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(url),
		},
	}
	resp, err := client.Get("https://github.com")
	if err != nil {
		log.Fatal(err)
		return
	}
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatal(err)
		return
	}
	// /への結果が表示される
	fmt.Println(string(dump))
}

func doGetToHTTPS(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{
		Transport: &http.Transport{
			//証明書の検証を行わない
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	values := url.Values{
		"name": {"hello world"},
	}
	resp, err := client.Get("https://localhost:18888" + "?" + values.Encode())
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(w, string(body))
	log.Println(string(body))
}

func main() {
	var httpServer http.Server
	http.HandleFunc("/", handler)

	http.HandleFunc("/error", returnErrorHandler)
	http.HandleFunc("/error/writeError", writeErrorHandler)
	http.HandleFunc("/error/notFound", returnNotFoundHandler)

	http.HandleFunc("/upload", multipartHandler)

	http.HandleFunc("/useCookie", useCookieHandler)
	//注意
	// respnse body をすべて読み切らない場合、
	// keep alive を利用できない
	http.HandleFunc("/doGet", doGet)
	http.HandleFunc("/doGetWithCooke", doGetWithCookie)
	http.HandleFunc("/doPost", doPost)
	http.HandleFunc("/doPostWithFile", doPostWithText)
	http.HandleFunc("/doPostWithMultipart", doPostWithMultipart)
	http.HandleFunc("/doPut", doPut)

	http.HandleFunc("/proxyDummy", proxyDummy)

	http.HandleFunc("/doGetToHTTPS", doGetToHTTPS)

	http.HandleFunc("/websocket", webSocketHandler)
	http.HandleFunc("/websocket/send", sendMessageHandler)

	log.Println("start http listen :18888")
	httpServer.Addr = ":18888"

	//http
	log.Println(httpServer.ListenAndServe())

	//オレオレ証明書による https 接続
	// log.Println(httpServer.ListenAndServeTLS("./server.crt", "./server.key"))
}
