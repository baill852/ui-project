package auth

import (
	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

type authUsecase struct {
	key []byte
}

func NewAuthUsecase() AuthUsecase {
	return &authUsecase{
		key: []byte(viper.GetString("secretKey")),
	}
}

func (u *authUsecase) GenerateToken(account string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"account": account,
	})
	tokenString, err := token.SignedString(u.key)
	return tokenString, err
}

func (u *authUsecase) ValidateToken(tokenString string) (jwt.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return u.key, nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims, err
}
