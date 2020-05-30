package main

import (
	"log"
	"net/http"
)

func recovery(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		f(w, r)
	}
}

func panicky(err error) {
	if err != nil {
		panic(err)
	}
}
