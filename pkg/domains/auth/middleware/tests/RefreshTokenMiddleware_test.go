package middleware

import (
	"github.com/stretchr/testify/assert"
	"paniq-auth-user-service/pkg/domains/auth/middleware"
	"testing"
)

func TestRefreshTokenMiddleware(t *testing.T) {

	t.Run("token exists", func(t *testing.T) {
		m := middleware.NewRefreshTokenMiddleware("token")
		assert.Equal(t, "token", m.GetRefreshToken())
		status, err := m.Validate()
		assert.True(t, status)
		assert.Nil(t, err)
	})

	t.Run("token not exists", func(t *testing.T) {
		m := middleware.NewRefreshTokenMiddleware("")
		assert.Equal(t, "", m.GetRefreshToken())
		status, err := m.Validate()
		assert.False(t, status)
		assert.NotNil(t, err)
	})
}
