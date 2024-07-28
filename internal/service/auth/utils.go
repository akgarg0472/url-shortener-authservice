package auth_service

import (
	"github.com/akgarg0472/urlshortener-auth-service/internal/entity"
	"strings"

	"golang.org/x/crypto/bcrypt"

	models "github.com/akgarg0472/urlshortener-auth-service/model"
	"github.com/google/uuid"
)

// function to validate provided password against the encrypted password stored in DB
func verifyPassword(rawPassword string, encryptedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(rawPassword)) == nil
}

func createUserEntity(request models.SignupRequest) *entity.User {
	return &entity.User{
		Id:       strings.ReplaceAll(uuid.New().String(), "-", ""),
		Email:    &request.Email,
		Password: &request.Password,
		Name:     request.Name,
		Scopes:   "user",
	}
}
