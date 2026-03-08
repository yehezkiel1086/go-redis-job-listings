package util

import "golang.org/x/crypto/bcrypt"

func HashPassword(pass []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
}

func ComparePassword(hash, password []byte) error {
	return bcrypt.CompareHashAndPassword(hash, password)
}
