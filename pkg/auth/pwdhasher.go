package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPwd(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return string(hashedPassword), fmt.Errorf("error hashing password: %v", err)
	}
	return string(hashedPassword), nil
}

func ComparePwd(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
