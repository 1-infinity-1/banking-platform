package auth

import "golang.org/x/crypto/bcrypt"

func checkPassword(hash string, password string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)

	return err == nil
}
