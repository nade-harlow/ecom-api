package helper

import (
	"math/rand"
	"time"
)

const stringLength = 6

func GenerateRandomString() string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, stringLength)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
