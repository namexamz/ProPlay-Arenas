package service

import (
	"errors"
	"time"
	"user-service/internal/models"
	"user-service/internal/repository"

	"github.com/golang-jwt/jwt/v4"
)

type AuthService struct {
	jwtSecret []byte
	userRepo  repository.UserRepository
}

func NewAuthService(secret string, userRepo repository.UserRepository) *AuthService {
	return &AuthService{
		jwtSecret: []byte(secret),
		userRepo:  userRepo,
	}
}

func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	claims := models.Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
func (s *AuthService) ParseToken(tokenStr string) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
