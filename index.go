package main

import (
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	type dataItem struct {
		UserID    int64
		FirstName string
		Picture   string
	}
	p := struct {
		Pipeline
		Data []dataItem
	}{
		Pipeline: Pipeline{
			Title:   "Where opposites attract!",
			Session: getSession(r),
		},
	}

	if r.URL.Path == "/" {
		db := getDB()
		rows, err := db.Query("SELECT user_id, first_name, picture FROM mismatch_user WHERE first_name IS NOT NULL ORDER BY join_date DESC LIMIT 5")
		panicky(err)
		for rows.Next() {
			item := dataItem{}
			rows.Scan(&item.UserID, &item.FirstName, &item.Picture)
			p.Data = append(p.Data, item)
		}
		tmpl.ExecuteTemplate(w, "index.html", p)
	}
}
