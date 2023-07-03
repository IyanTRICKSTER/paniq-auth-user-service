package bcryptUtils

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestHash(t *testing.T) {

	hashFunction := NewHashFunction()

	t.Run("Test Hash Success", func(t *testing.T) {
		payload := "MyPassword"
		hashed := hashFunction.Hash(payload)

		check, err := hashFunction.HashCheck(hashed, payload)
		assert.True(t, check)
		assert.Nil(t, err)
	})

	t.Run("Test Check Hash Failed", func(t *testing.T) {

		payload := "MyPassword"
		hashed := hashFunction.Hash(payload)

		check, err := hashFunction.HashCheck(hashed, "myPassword")
		assert.False(t, check)
		assert.ErrorIs(t, bcrypt.ErrMismatchedHashAndPassword, err)
		assert.NotNil(t, err)

	})
}
