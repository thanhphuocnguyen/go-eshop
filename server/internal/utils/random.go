package utils

import (
	"math/rand"
	"time"
)

var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func RandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func GenerateSKU() string {
	return RandomString(10)
}
