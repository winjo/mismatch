package main

import (
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var tmpl *template.Template

// Pipeline common form template.
type Pipeline struct {
	Session  session
	Title    string
	ErrorMsg string
}

func main() {
	tmpl = template.Must(template.New("").ParseGlob("view/*.html"))

	http.HandleFunc("/images/", recovery(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	}))
	http.HandleFunc("/", recovery(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			index(w, r)
		} else {
			http.NotFound(w, r)
		}
	}))
	http.HandleFunc("/login", recovery(login))
	http.HandleFunc("/logout", recovery(logout))
	http.HandleFunc("/signup", recovery(signup))
	http.HandleFunc("/style.css", recovery(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "view/style.css")
	}))

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal(nil)
	}
}
