package main

import (
	"net/http"
	"time"
)

func logout(w http.ResponseWriter, r *http.Request) {
	s := getSession(r)
	if s.Userid != 0 {
		cookie := http.Cookie{Name: sessionID, Expires: time.Now().AddDate(0, 0, -1)}
		http.SetCookie(w, &cookie)
	}
	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusFound)
}
