package main

import (
	"bytes"
	"database/sql"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func editprofile(w http.ResponseWriter, r *http.Request) {
	s := getSession(r)
	p := struct {
		Pipeline
		NotFound bool
		Userinfo struct {
			FirstName  string
			LastName   string
			Gender     string
			Birthdate  string
			City       string
			State      string
			OldPicture string
		}
		SuccessSubmit bool
	}{
		Pipeline: Pipeline{
			Title:   "Edit Profile",
			Session: getSession(r),
		},
	}

	if s.Userid != 0 {
		if r.Method == "POST" {
			picFile, picHandler, err := r.FormFile("new_picture")
			p.Userinfo.FirstName = strings.Trim(r.Form.Get("firstname"), " ")
			p.Userinfo.LastName = strings.Trim(r.Form.Get("lastname"), " ")
			p.Userinfo.Gender = strings.Trim(r.Form.Get("gender"), " ")
			p.Userinfo.Birthdate = strings.Trim(r.Form.Get("birthdate"), " ")
			p.Userinfo.City = strings.Trim(r.Form.Get("city"), " ")
			p.Userinfo.State = strings.Trim(r.Form.Get("state"), " ")
			p.Userinfo.OldPicture = strings.Trim(r.Form.Get("old_picture"), " ")
			var newPicture string
			if err == nil {
				defer picFile.Close()
				contentType := picHandler.Header.Get("Content-Type")
				b, _ := ioutil.ReadAll(picFile)
				imageReader := bytes.NewReader(b)
				config, _, _ := image.DecodeConfig(imageReader)
				size := picHandler.Size

				if (contentType == "image/gif" || contentType == "image/jpeg" || contentType == "image/pjpeg" || contentType == "image/png") && size > 0 && size <= 32<<10 && config.Width <= 120 && config.Height <= 120 {
					f, err := os.OpenFile("images/"+picHandler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
					if err == nil {
						defer f.Close()
						newPicture = picHandler.Filename
						fileReader := bytes.NewReader(b)
						io.Copy(f, fileReader)
						if p.Userinfo.OldPicture != newPicture {
							os.Remove("images/" + p.Userinfo.OldPicture)
						}
					} else {
						p.ErrorMsg = "Sorry, there was a problem uploading your picture."
					}
				} else {
					p.ErrorMsg = "Your picture must be a GIF, JPEG, or PNG image file no greater than 32 KB and 120x120 pixels in size."
				}
			}
			if p.ErrorMsg == "" {
				userinfo := p.Userinfo
				if userinfo.FirstName != "" && userinfo.LastName != "" && userinfo.Gender != "" && userinfo.Birthdate != "" && userinfo.City != "" && userinfo.State != "" {
					if newPicture != "" {
						_, err := db.Exec(
							"UPDATE mismatch_user SET first_name = ?, last_name = ?, gender = ?, birthdate = ?, city = ?, state = ?, picture = ? WHERE user_id = ?",
							userinfo.FirstName, userinfo.LastName, userinfo.Gender, userinfo.Birthdate, userinfo.City, userinfo.State, newPicture, s.Userid,
						)
						panicky(err)
					} else {
						_, err := db.Exec(
							"UPDATE mismatch_user SET first_name = ?, last_name = ?, gender = ?, birthdate = ?, city = ?, state = ? WHERE user_id = ?",
							userinfo.FirstName, userinfo.LastName, userinfo.Gender, userinfo.Birthdate, userinfo.City, userinfo.State, s.Userid,
						)
						panicky(err)
					}
					p.SuccessSubmit = true
				} else {
					p.ErrorMsg = "You must enter all of the profile data (the picture is optional)."
				}
			}
		} else {
			rows, err := db.Query("SELECT first_name, last_name, gender, birthdate, city, state, picture FROM mismatch_user WHERE user_id = ?", s.Userid)
			panicky(err)
			if rows.Next() {
				firstName := sql.NullString{}
				lastName := sql.NullString{}
				gender := sql.NullString{}
				birthdate := sql.NullString{}
				city := sql.NullString{}
				state := sql.NullString{}
				oldPicture := sql.NullString{}
				err := rows.Scan(
					&firstName,
					&lastName,
					&gender,
					&birthdate,
					&city,
					&state,
					&oldPicture,
				)
				panicky(err)
				p.Userinfo.FirstName = firstName.String
				p.Userinfo.LastName = lastName.String
				p.Userinfo.Gender = gender.String
				p.Userinfo.Birthdate = birthdate.String
				p.Userinfo.City = city.String
				p.Userinfo.State = state.String
				p.Userinfo.OldPicture = oldPicture.String
			} else {
				p.NotFound = true
			}
		}
	}

	err := tmpl.ExecuteTemplate(w, "editprofile.html", p)
	if err != nil {
		log.Println(err)
	}
}
