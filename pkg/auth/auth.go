package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func BcryptHash(password string) string {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		panic(err)
	}

	return string(hashedBytes)
}

func BcryptCheck(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
