package main

import (
	"database/sql"
	"log"
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {
	p := struct {
		Pipeline
		SuccessCreated bool
		Username       string
	}{
		Pipeline: Pipeline{
			Title:   "Log In",
			Session: getSession(r),
		},
	}

	if r.Method == "POST" {
		r.ParseForm()
		username, password := r.Form.Get("username"), r.Form.Get("password")
		p.Username = username
		if username != "" && password != "" {
			db := getDB()
			defer db.Close()
			var userID int
			var userName string
			err := db.
				QueryRow("SELECT user_id, username FROM mismatch_user WHERE username = ? AND password = ?", username, sha(password)).
				Scan(&userID, &userName)
			switch {
			case err == sql.ErrNoRows:
				p.ErrorMsg = "Sorry, you must enter a valid username and password to log in."
			case err != nil:
				panicky(err)
			default:
				setSession(w, session{userID, userName})
				w.Header().Set("Location", "/")
				w.WriteHeader(http.StatusFound)
			}
		} else {
			p.ErrorMsg = "Sorry, you must enter your username and password to log in."
		}
	}
	err := tmpl.ExecuteTemplate(w, "login.html", p)
	if err != nil {
		log.Println(err)
	}
}
