package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/textproto"
	"net/url"
	"os"
	"strconv"
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

func multipartHandler(w http.ResponseWriter, r *http.Request) {
	dump, _ := httputil.DumpRequest(r, true)
	fmt.Println(string(dump))

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
	}
	file, _, err := r.FormFile("attachment-file")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	savePath := "D:/hoge.txt"
	saveFile, err := os.Create(savePath)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer saveFile.Close()

	_, err = io.Copy(saveFile, file)
	if err != nil {
		log.Fatal(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func useCookieHandler(w http.ResponseWriter, r *http.Request) {
	dump, _ := httputil.DumpRequest(r, true)
	fmt.Println(string(dump))

	storeCookie := &http.Cookie{
		Name:  "COUNT",
		Value: "1",
	}
	cookie, _ := r.Cookie("COUNT")

	if cookie != nil {
		if i, err := strconv.Atoi(cookie.Value); err == nil {
			storeCookie.Value = strconv.Itoa(i + 1)
		}
	}

	http.SetCookie(w, storeCookie)

	//COUNT key exsit in Cookie
	if cookie != nil {
		fmt.Fprintf(w, "<html><body>cookie content is %s</body></html>", cookie.Value)
	} else {
		fmt.Fprintln(w, "<html><body>no cookie</body></html>")
	}
}

func doGet(w http.ResponseWriter, r *http.Request) {
	values := url.Values{
		"name": {"hello world"},
	}
	resp, err := http.Get("http://127.0.0.1:18888" + "?" + values.Encode())
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

func doGetWithCookie(w http.ResponseWriter, r *http.Request) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	client := http.Client{
		Jar: jar,
	}
	for index := 0; index < 10; index++ {
		resp, err := client.Get("http://127.0.0.1:18888/useCookie")
		if err != nil {
			log.Fatal(err)
			return
		}
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			log.Fatal(err)
			return
		}
		//カウントが進むことを確認できる
		fmt.Println(string(dump))
	}
}

func doPost(w http.ResponseWriter, r *http.Request) {
	values := url.Values{
		"name": {"hello world post"},
	}
	resp, err := http.PostForm("http://127.0.0.1:18888", values)
	if err != nil {
		panic(err)
	}
	log.Println(resp.StatusCode)
	log.Println(resp.Status)
}

func doPostWithText(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("D:/test.txt")
	if err != nil {
		panic(err)
	}
	resp, err := http.Post("http://127.0.0.1:18888/", "text/plan", file)
	if err != nil {
		panic(err)
	}
	log.Println(resp.StatusCode)
	log.Println(resp.Status)
}

func doPostWithMultipart(w http.ResponseWriter, r *http.Request) {
	var buffer bytes.Buffer
	//boundary も決まる
	writer := multipart.NewWriter(&buffer)
	writer.WriteField("name", "bob")
	writer.WriteField("greeting", "hello world")

	//application/octet-stream になる
	// fileWriter, err := writer.CreateFormFile("attachment-file", "D:/test.txt")
	part := make(textproto.MIMEHeader)
	part.Set("Content-Type", "text/plain")
	part.Set("Content-Disposition", `form-data; name="attachment-file"; filename="test.txt"`)
	fileWriter, err := writer.CreatePart(part)
	if err != nil {
		panic(err)
	}
	file, err := os.Open("D:/test.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	io.Copy(fileWriter, file)
	writer.Close()

	resp, err := http.Post("http://127.0.0.1:18888/upload", writer.FormDataContentType(), &buffer)
	if err != nil {
		panic(err)
	}
	log.Println(resp.StatusCode)
	log.Println(resp.Status)
}

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

func doPut(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}
	request, err := http.NewRequest("PUT", "http://127.0.0.1:18888", nil)
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(dump))
}

func main() {
	var httpServer http.Server
	http.HandleFunc("/", handler)

	http.HandleFunc("/error", returnErrorHandler)
	http.HandleFunc("/error/writeError", writeErrorHandler)
	http.HandleFunc("/error/notFound", returnNotFoundHandler)

	http.HandleFunc("/upload", multipartHandler)

	http.HandleFunc("/useCookie", useCookieHandler)

	http.HandleFunc("/doGet", doGet)
	http.HandleFunc("/doGetWithCooke", doGetWithCookie)
	http.HandleFunc("/doPost", doPost)
	http.HandleFunc("/doPostWithFile", doPostWithText)
	http.HandleFunc("/doPostWithMultipart", doPostWithMultipart)
	http.HandleFunc("/doPut", doPut)

	http.HandleFunc("/proxyDummy", proxyDummy)

	log.Println("start http listen :18888")
	httpServer.Addr = ":18888"
	log.Println(httpServer.ListenAndServe())
}
