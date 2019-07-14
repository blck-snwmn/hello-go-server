package main

import "net/http"

func writeErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

func returnErrorHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "error!!", http.StatusBadRequest)
}

func returnNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}
