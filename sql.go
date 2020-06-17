package main

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"os"
)

func sha(s1 string) string {
	a := sha1.Sum([]byte(s1))
	return hex.EncodeToString(a[:])
}

var db *sql.DB

func openDB() {
	var err error
	dns := os.Getenv("MYSQL_DNS")
	if dns == "" {
		dns = "root:abc123@/mismatch?charset=utf8mb4"
	}
	db, err = sql.Open("mysql", dns)
	panicky(err)
}
