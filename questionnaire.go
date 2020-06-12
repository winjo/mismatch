package main

import (
	"database/sql"
	"log"
	"net/http"
)

func questionnaire(w http.ResponseWriter, r *http.Request) {
	type responseItem struct {
		ResponseID   int
		TopicID      int
		Response     int8
		TopicName    string
		CategoryName string
	}
	p := struct {
		Pipeline
		Response      []responseItem
		SuccessSubmit bool
	}{
		Pipeline: Pipeline{
			Title:   "Questionnaire",
			Session: getSession(r),
		},
	}

	s := getSession(r)
	if s.Userid != 0 {
		rows, err := db.Query("SELECT * FROM mismatch_response WHERE user_id = ?", s.Userid)
		if ok := rows.Next(); !ok {
			rows, err := db.Query("SELECT topic_id FROM mismatch_topic ORDER BY category_id, topic_id")
			panicky(err)
			stmt, err := db.Prepare("INSERT INTO mismatch_response (user_id, topic_id) VALUES (?, ?)")
			panicky(err)
			for rows.Next() {
				var topicID string
				err := rows.Scan(&topicID)
				panicky(err)
				_, err = stmt.Exec(s.Userid, topicID)
				panicky(err)
			}
		}
		if r.Method == "POST" {
			r.ParseForm()
			stmt, err := db.Prepare("UPDATE mismatch_response SET response = ? WHERE response_id = ?")
			panicky(err)
			for key, value := range r.Form {
				if key != "submit" && len(value) != 0 {
					_, err := stmt.Exec(value[0], key)
					panicky(err)
				}
			}
			p.SuccessSubmit = true
		}
		responseRows, err := db.Query(`
			SELECT mr.response_id, mr.topic_id, mr.response, mt.name AS topic_name, mc.name AS category_name 
			FROM mismatch_response AS mr 
			INNER JOIN mismatch_topic AS mt USING (topic_id) 
			INNER JOIN mismatch_category AS mc USING (category_id) 
			WHERE mr.user_id = ?
		`, s.Userid)
		panicky(err)
		for responseRows.Next() {
			item := responseItem{}
			response := sql.NullInt32{} // response 可能为 NULL
			err := responseRows.Scan(
				&item.ResponseID,
				&item.TopicID,
				&response,
				&item.TopicName,
				&item.CategoryName,
			)
			panicky(err)
			if !response.Valid {
				item.Response = 0
			} else {
				item.Response = int8(response.Int32)
			}
			p.Response = append(p.Response, item)
		}
	}

	err := tmpl.ExecuteTemplate(w, "questionnaire.html", p)
	if err != nil {
		log.Println(err)
	}
}
