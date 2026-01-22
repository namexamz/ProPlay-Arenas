package models

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}
