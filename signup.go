package main

import (
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
			db := getDB()
			defer func() {
				db.Close()
			}()
			stmt, err := db.Prepare("SELECT * FROM mismatch_user WHERE username = ?")
			panicky(err)
			rows, err := stmt.Query(username)
			panicky(err)
			if ok := rows.Next(); ok {
				p.ErrorMsg = "An account already exists for this username. Please use a different address."
			} else {
				stmt, err := db.Prepare("INSERT INTO mismatch_user (username, password, join_date) VALUES (?, ?, ?)")
				panicky(err)
				stmt.Exec(username, sha(password1), time.Now().Format("2006-01-02 15:04:05"))
				p.SuccessCreated = true
			}
		} else {
			p.ErrorMsg = "You must enter all of the sign-up data, including the desired password twice."
			p.Username = username
		}
	}

	tmpl.ExecuteTemplate(w, "signup.html", p)
}
