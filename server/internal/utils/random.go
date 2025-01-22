package utils

import (
	crypto "crypto/rand"
	"math/big"
	mathRand "math/rand"
	"time"
)

const n = 10

var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func init() {
	mathRand.New(mathRand.NewSource(time.Now().UnixNano()))
}

func RandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[mathRand.Intn(len(letter))]
	}
	return string(b)
}

func GenerateSKU() string {
	return RandomString(10)
}

func GenerateRandomDiscountCode() (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := crypto.Int(crypto.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}
