package utils

import (
	"fmt"
	"time"

	"github.com/Auxesia23/toko-online/internal/env"
	"github.com/Auxesia23/toko-online/internal/models"
	"github.com/golang-jwt/jwt"
)

var jwtSecret []byte

func InitJwt() {
	secret := env.GetString("SECRET_KEY", "")
	jwtSecret = []byte(secret)
}

func GenerateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":      user.ID,
		"is_superuser": user.Superuser,
		"email":        user.Email,
		"exp":          time.Now().Add(time.Hour * 24).Unix(),
		"iat":          time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
