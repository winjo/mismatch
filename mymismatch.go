package main

import (
	"database/sql"
	"log"
	"net/http"
)

func mymismatch(w http.ResponseWriter, r *http.Request) {
	s := getSession(r)
	p := struct {
		Pipeline
		Userinfo struct {
			UserName  string
			FirstName string
			LastName  string
			City      string
			State     string
			Picture   string
		}
		MismatchUserID int
		MismatchTopics []string
		ShouldQuestion bool
	}{
		Pipeline: Pipeline{
			Title:   "My Mismatch",
			Session: getSession(r),
		},
		MismatchUserID: -1,
	}

	if s.Userid != 0 {
		db := getDB()
		defer db.Close()
		rows, err := db.Query("SELECT * FROM mismatch_response WHERE user_id = ?", s.Userid)
		panicky(err)
		if rows.Next() {
			rows, err := db.Query(`
				SELECT mr.response_id, mr.topic_id, mr.response, mt.name AS topic_name 
				FROM mismatch_response AS mr 
				INNER JOIN mismatch_topic AS mt USING (topic_id) 
				WHERE mr.user_id = ?
			`, s.Userid)
			panicky(err)
			type userResponse struct {
				responseID int
				topicID    int
				response   int8
				topicName  string
			}
			userResponses := []userResponse{}
			for rows.Next() {
				item := userResponse{}
				response := sql.NullInt32{} // response 可能为 NULL
				err = rows.Scan(
					&item.responseID,
					&item.topicID,
					&response,
					&item.topicName,
				)
				panicky(err)
				if !response.Valid {
					item.response = 0
				} else {
					item.response = int8(response.Int32)
				}
				userResponses = append(userResponses, item)
			}
			mismathcScore, mismathcUserID, mismatchTopics := 0, -1, []string{}
			rows, err = db.Query("SELECT user_id FROM mismatch_user WHERE user_id != ?", s.Userid)
			panicky(err)
			for rows.Next() {
				var userID int
				err := rows.Scan(&userID)
				panicky(err)
				rows, err := db.Query("SELECT response_id, topic_id, response FROM mismatch_response WHERE user_id = ?", userID)
				type mismatchResponse struct {
					responseID int
					topicID    int
					response   int8
				}
				mismatchResponses := []mismatchResponse{}
				for rows.Next() {
					item := mismatchResponse{}
					response := sql.NullInt32{} // response 可能为 NULL
					err := rows.Scan(
						&item.responseID,
						&item.topicID,
						&response,
					)
					panicky(err)
					if !response.Valid {
						item.response = 0
					} else {
						item.response = int8(response.Int32)
					}
					mismatchResponses = append(mismatchResponses, item)
				}
				if len(mismatchResponses) == 0 {
					continue
				}
				score, topics := 0, []string{}
				for i := range userResponses {
					if userResponses[i].response+mismatchResponses[i].response == 3 {
						score++
						topics = append(topics, userResponses[i].topicName)
					}
				}
				if score > mismathcScore {
					mismathcScore = score
					mismathcUserID = userID
					mismatchTopics = topics[:]
				}
			}
			if mismathcUserID != -1 {
				db.
					QueryRow("SELECT username, first_name, last_name, city, state, picture FROM mismatch_user WHERE user_id = ?", mismathcUserID).
					Scan(
						&p.Userinfo.UserName,
						&p.Userinfo.FirstName,
						&p.Userinfo.LastName,
						&p.Userinfo.City,
						&p.Userinfo.State,
						&p.Userinfo.Picture,
					)
				p.MismatchUserID = mismathcUserID
				p.MismatchTopics = mismatchTopics
			}
		} else {
			p.ShouldQuestion = true
		}
	}

	err := tmpl.ExecuteTemplate(w, "mymismatch.html", p)
	if err != nil {
		log.Println(err)
	}
}
