package domain

import "github.com/golang-jwt/jwt/v5"

type TokenClaims struct {
	Name   string `json:"name"`
	UserID string `json:"userID"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}
