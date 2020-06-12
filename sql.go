package main

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
)

func sha(s1 string) string {
	a := sha1.Sum([]byte(s1))
	return hex.EncodeToString(a[:])
}

var db *sql.DB

func openDB() {
	var err error
	db, err = sql.Open("mysql", "root:abc123@/mismatch?charset=utf8mb4")
	panicky(err)
}
