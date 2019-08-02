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
