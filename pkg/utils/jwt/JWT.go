package jwtUtils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"paniq-auth-user-service/pkg/contracts"
	"time"
)

type JWTService struct{}

func (J *JWTService) GenerateToken(userID uint, lifeSpan int, secretKey string) (string, error) {

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Second * time.Duration(lifeSpan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secretKey))
}

func (J *JWTService) ValidateToken(token string, secretKey string) bool {

	_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return false
	}
	return true

}

func (J *JWTService) ExtractPayloadFromToken(token string, secretKey string) (map[string]interface{}, error) {

	extractedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := extractedToken.Claims.(jwt.MapClaims)
	if ok && extractedToken.Valid {
		return claims, nil
	}
	return nil, nil
}

func New() contracts.IJWTService {
	return &JWTService{}
}
