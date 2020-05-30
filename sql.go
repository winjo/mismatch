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

func getDB() *sql.DB {
	db, err := sql.Open("mysql", "root:abc123@/mismatch?charset=utf8mb4")
	panicky(err)
	return db
}
