package main

import (
	"log"
	"net/http"
	"time"
)

func signup(w http.ResponseWriter, r *http.Request) {
	p := struct {
		Pipeline
		SuccessCreated bool
		Username       string
	}{
		Pipeline: Pipeline{
			Title:   "Sign Up",
			Session: getSession(r),
		},
	}

	if r.Method == "POST" {
		r.ParseForm()
		username := r.Form.Get("username")
		password1 := r.Form.Get("password1")
		password2 := r.Form.Get("password2")
		if username != "" && password1 != "" && password2 != "" && password1 == password2 {
			rows, err := db.Query("SELECT * FROM mismatch_user WHERE username = ?", username)
			panicky(err)
			if ok := rows.Next(); ok {
				p.ErrorMsg = "An account already exists for this username. Please use a different address."
			} else {
				_, err := db.Exec(
					"INSERT INTO mismatch_user (username, password, join_date) VALUES (?, ?, ?)",
					username, sha(password1), time.Now().Format("2006-01-02 15:04:05"),
				)
				panicky(err)
				p.SuccessCreated = true
			}
		} else {
			p.ErrorMsg = "You must enter all of the sign-up data, including the desired password twice."
			p.Username = username
		}
	}

	err := tmpl.ExecuteTemplate(w, "signup.html", p)
	if err != nil {
		log.Println(err)
	}
}
