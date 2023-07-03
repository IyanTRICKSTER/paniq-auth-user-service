package jwtUtils

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestJWT(t *testing.T) {

	jwtService := New()

	userID := uint(10)
	secretKey := "88992727"

	t.Run("test generate token success", func(t *testing.T) {
		token, err := jwtService.GenerateToken(userID, 7200, secretKey)
		assert.Nil(t, err)
		assert.NotNil(t, token)
	})

	t.Run("validate token success", func(t *testing.T) {

		token, err := jwtService.GenerateToken(userID, 7200, secretKey)
		assert.Nil(t, err)

		status := jwtService.ValidateToken(token, secretKey)
		assert.True(t, status)

	})

	t.Run("validate token failed", func(t *testing.T) {

		token, err := jwtService.GenerateToken(userID, 7200, secretKey)
		assert.Nil(t, err)

		status := jwtService.ValidateToken(token+"wkwk", secretKey)
		assert.False(t, status)

	})

	t.Run("validate expired token failed", func(t *testing.T) {

		status := jwtService.ValidateToken("wkwk", secretKey)
		assert.False(t, status)

	})

	t.Run("extract token success", func(t *testing.T) {

		token, err := jwtService.GenerateToken(userID, 7200, secretKey)
		assert.Nil(t, err)

		payload, err := jwtService.ExtractPayloadFromToken(token, secretKey)
		log.Println(token)
		assert.Nil(t, err)
		assert.NotNil(t, payload)
		assert.Equal(t, userID, uint(payload["user_id"].(float64)))

	})

	t.Run("extract token failed", func(t *testing.T) {

		_, err := jwtService.GenerateToken(userID, 7200, secretKey)
		assert.Nil(t, err)

		payload, err := jwtService.ExtractPayloadFromToken("wkwk", "wkwk")
		assert.NotNil(t, err)
		assert.Nil(t, payload)
		//assert.Equal(t, userID, uint(payload["user_id"].(float64)))

	})
}
