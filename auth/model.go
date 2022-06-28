package auth

import "github.com/golang-jwt/jwt"

type AuthUsecase interface {
	GenerateToken(account string) (string, error)
	ValidateToken(tokenString string) (jwt.Claims, error)
}
