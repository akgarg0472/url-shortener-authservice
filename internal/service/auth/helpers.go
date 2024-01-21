package auth_service

import "golang.org/x/crypto/bcrypt"

// function to validate provided password against the encrypted password stored in DB
func verifyPassword(rawPassword string, encryptedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(rawPassword)) == nil
}
