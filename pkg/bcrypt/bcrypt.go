package bcrypt

import (
	"enigmanations/cats-social/pkg/env"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bcryptSalt, err := env.GetEnvInt("BCRYPT_SALT")
	if err != nil {
		return "", err
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptSalt)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
