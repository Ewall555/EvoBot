package passwd

import (
	"golang.org/x/crypto/bcrypt"
)

func Encode(rawPassword string) string {
	if len(rawPassword) == 0 {
		return ""
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hash)
}

func Verify(rawPassword, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(rawPassword))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false
		}
		return false
	}
	return true
}
