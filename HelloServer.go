package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
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

func writeErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

func returnErrorHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "error!!", http.StatusBadRequest)
}

func returnNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func main() {
	var httpServer http.Server
	http.HandleFunc("/", handler)

	http.HandleFunc("/error", returnErrorHandler)
	http.HandleFunc("/error/writeError", writeErrorHandler)
	http.HandleFunc("/error/notFound", returnNotFoundHandler)

	http.HandleFunc("/upload", multipartHandler)

	http.HandleFunc("/useCookie", useCookieHandler)

	log.Println("start http listen :18888")
	httpServer.Addr = ":18888"
	log.Println(httpServer.ListenAndServe())
}
