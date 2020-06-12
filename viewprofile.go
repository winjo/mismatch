package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func viewprofile(w http.ResponseWriter, r *http.Request) {
	s := getSession(r)
	p := struct {
		Pipeline
		Userinfo struct {
			UserName  string
			FirstName string
			LastName  string
			Gender    string
			Birthdate string
			Birthyear string
			City      string
			State     string
			Picture   string
		}
		NotFound    bool
		QueryUserid int
	}{
		Pipeline: Pipeline{
			Title:   "View Profile",
			Session: getSession(r),
		},
	}

	if s.Userid != 0 {
		r.ParseForm()
		queryUserID, _ := strconv.Atoi(r.Form.Get("user_id"))
		p.QueryUserid = queryUserID
		var userID int
		if queryUserID != 0 {
			userID = queryUserID
		} else {
			userID = s.Userid
		}
		firstName := sql.NullString{}
		lastName := sql.NullString{}
		gender := sql.NullString{}
		birthdate := sql.NullString{}
		city := sql.NullString{}
		state := sql.NullString{}
		picture := sql.NullString{}
		err := db.
			QueryRow("SELECT username, first_name, last_name, gender, birthdate, city, state, picture FROM mismatch_user WHERE user_id = ?", userID).
			Scan(
				&p.Userinfo.UserName,
				&firstName,
				&lastName,
				&gender,
				&birthdate,
				&city,
				&state,
				&picture,
			)
		switch {
		case err == sql.ErrNoRows:
			p.NotFound = true
		case err != nil:
			panicky(err)
		default:
			p.Userinfo.FirstName = firstName.String
			p.Userinfo.LastName = lastName.String
			p.Userinfo.Gender = gender.String
			p.Userinfo.Birthdate = birthdate.String
			p.Userinfo.City = city.String
			p.Userinfo.State = state.String
			p.Userinfo.Picture = picture.String
			if p.Userinfo.Birthdate != "" {
				p.Userinfo.Birthyear = strings.Split(p.Userinfo.Birthdate, "-")[0]
			}
		}
	}
	err := tmpl.ExecuteTemplate(w, "viewprofile.html", p)
	if err != nil {
		log.Println(err)
	}
}
