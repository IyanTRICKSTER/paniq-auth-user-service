package middleware

import (
	"github.com/stretchr/testify/assert"
	"paniq-auth-user-service/pkg/contracts/permissionCodes"
	"paniq-auth-user-service/pkg/domains/auth/middleware"
	"testing"
)

func TestAuthenticatedMiddleware(t *testing.T) {
	m := middleware.NewAuthenticatedRequestMiddleware(
		10,
		"iyan",
		"admin",
		"ld")

	assert.Equal(t, uint(10), m.GetUserID())
	assert.Equal(t, "iyan", m.GetUsername())
	assert.Equal(t, "admin", m.GetRole())
	assert.Equal(t, "ld", m.GetPermissionCodes())
	assert.True(t, m.HasRole("admin"))
	assert.False(t, m.HasRole("user"))
	assert.True(t, m.HasPermission(permissionCodes.LIST))
	assert.True(t, m.HasPermission(permissionCodes.DELETE))

	status, err := m.Validate()
	assert.True(t, status)
	assert.Nil(t, err)

}

func TestAuthenticatedMiddlewareFail(t *testing.T) {
	m := middleware.NewAuthenticatedRequestMiddleware(
		0,
		"iyan",
		"admin",
		"ld")

	assert.Equal(t, uint(0), m.GetUserID())
	assert.Equal(t, "iyan", m.GetUsername())
	assert.Equal(t, "admin", m.GetRole())
	assert.Equal(t, "ld", m.GetPermissionCodes())
	assert.True(t, m.HasRole("admin"))
	assert.False(t, m.HasRole("user"))
	assert.True(t, m.HasPermission(permissionCodes.LIST))
	assert.True(t, m.HasPermission(permissionCodes.DELETE))

	status, err := m.Validate()
	assert.False(t, status)
	assert.NotNil(t, err)

}
