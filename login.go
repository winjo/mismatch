package main

import (
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
			rows, err := db.Query("SELECT user_id, username FROM mismatch_user WHERE username = ? AND password = ?", username, sha(password))
			panicky(err)
			if ok := rows.Next(); ok {
				var userid int64
				var username string
				err := rows.Scan(&userid, &username)
				panicky(err)
				setSession(w, session{userid, username})
				w.Header().Set("Location", "/")
				w.WriteHeader(http.StatusFound)
			} else {
				p.ErrorMsg = "Sorry, you must enter a valid username and password to log in."
			}
		} else {
			p.ErrorMsg = "Sorry, you must enter your username and password to log in."
		}
	}
	tmpl.ExecuteTemplate(w, "login.html", p)
}
