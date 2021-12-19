package util

import (
	"golang.org/x/crypto/bcrypt"
)

func GetBcrypt(pw string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pw), Config.BCRYPT_COST)
	return string(bytes), err
}

func HashIsValid(hash, pw string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw)); err != nil {
		return false
	}
	return true
}
