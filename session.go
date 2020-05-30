package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"
)

const sessionID = "SESSION_ID"
const key = "go-go-go-fun-fun"

type session struct {
	Userid   int64
	Username string
}

func getSession(r *http.Request) session {
	cookie, _ := r.Cookie(sessionID)
	if cookie == nil {
		return session{}
	}
	return string2Session(cookie.Value)
}

func setSession(w http.ResponseWriter, s session) {
	str := session2String(s)
	cookie := http.Cookie{Name: sessionID, Value: str, Expires: time.Now().AddDate(0, 0, 1)}
	http.SetCookie(w, &cookie)
}

func string2Session(str string) (s session) {
	s = session{}
	k := []byte(key)
	ciphertext, err := base64.RawStdEncoding.DecodeString(str)
	if err != nil {
		return
	}
	block, err := aes.NewCipher(k)
	if err != nil {
		return
	}
	if len(ciphertext) < aes.BlockSize {
		return
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	if len(ciphertext)%aes.BlockSize != 0 {
		return
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)
	ciphertext = unpadding(ciphertext, aes.BlockSize)
	_ = json.Unmarshal(ciphertext, &s)
	return
}

func session2String(s session) string {
	text, err := json.Marshal(s)
	k := []byte(key)
	if err != nil {
		return ""
	}
	text = padding(text, aes.BlockSize)
	block, err := aes.NewCipher(k)
	if err != nil {
		return ""
	}
	ciphertext := make([]byte, aes.BlockSize+len(text))
	iv := ciphertext[:aes.BlockSize]
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], text)
	return base64.RawStdEncoding.EncodeToString(ciphertext)
}

func padding(ciphertext []byte, blockSize int) []byte {
	num := blockSize - len(ciphertext)%blockSize
	pad := bytes.Repeat([]byte{byte(num)}, num)
	return append(ciphertext, pad...)
}

func unpadding(text []byte, blockSize int) []byte {
	n := len(text)
	num := int(text[n-1])
	if num <= blockSize {
		return text[:(n - num)]
	}
	return text
}
