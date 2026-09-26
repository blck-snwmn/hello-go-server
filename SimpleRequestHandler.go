package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/textproto"
	"net/url"
	"os"
	"strings"
)

func doGet(w http.ResponseWriter, r *http.Request) {
	values := url.Values{
		"name": {"hello world"},
	}
	resp, err := http.Get("http://127.0.0.1:18888" + "?" + values.Encode())
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	_, _ = fmt.Fprintln(w, string(body))
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
	for range 10 {
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
	if err := writer.WriteField("name", "bob"); err != nil {
		panic(err)
	}
	if err := writer.WriteField("greeting", "hello world"); err != nil {
		panic(err)
	}

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
	defer func() {
		_ = file.Close()
	}()
	if _, err := io.Copy(fileWriter, file); err != nil {
		panic(err)
	}
	if err := writer.Close(); err != nil {
		panic(err)
	}

	resp, err := http.Post("http://127.0.0.1:18888/upload", writer.FormDataContentType(), &buffer)
	if err != nil {
		panic(err)
	}
	log.Println(resp.StatusCode)
	log.Println(resp.Status)
}

func doPut(w http.ResponseWriter, r *http.Request) {
	//NewRequestの第三引数の渡し方は、http.PostForm等参照
	values := url.Values{"greeting": {"put values"}}

	client := &http.Client{}

	request, err := http.NewRequest(
		"PUT",
		"http://127.0.0.1:18888",
		strings.NewReader(values.Encode()),
	)
	//ParseForm の対象にするため
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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
