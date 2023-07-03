package response

import (
	"github.com/stretchr/testify/assert"
	"paniq-auth-user-service/pkg/contracts/statusCodes"
	"paniq-auth-user-service/pkg/entities"
	"testing"
)

func TestResponse(t *testing.T) {
	t.Run("create user response", func(t *testing.T) {
		response := New(UserResponse{}, true, statusCodes.Error, "ok", entities.UserEntity{})
		assert.True(t, response.GetStatus())
		assert.Equal(t, statusCodes.Error, response.GetStatusCode())
		assert.Equal(t, "ok", response.GetMessage())
		response.SetMessage("ko")
		assert.Equal(t, "ko", response.GetMessage())
		assert.True(t, response.ErrorIs(statusCodes.Error))
		assert.False(t, response.ErrorIs(statusCodes.Success))
		assert.False(t, response.IsFailed())
		assert.IsType(t, entities.UserEntity{}, response.GetData())
		assert.IsType(t, UserResponse{}, response)
		assert.Equal(t, response.ToMapStringInterface()["data"], response.GetData())
	})

	t.Run("create auth response", func(t *testing.T) {
		response := New(AuthResponse{}, true, statusCodes.Error, "ok", entities.UserEntity{})
		assert.True(t, response.GetStatus())
		assert.Equal(t, statusCodes.Error, response.GetStatusCode())
		assert.Equal(t, "ok", response.GetMessage())
		assert.IsType(t, entities.UserEntity{}, response.GetData())
		assert.IsType(t, AuthResponse{}, response)
	})

	t.Run("create auth response with status false", func(t *testing.T) {
		response := New(AuthResponse{}, false, statusCodes.Error, "ok", entities.UserEntity{})
		assert.False(t, response.GetStatus())
		assert.True(t, response.IsFailed())
		assert.Equal(t, statusCodes.Error, response.GetStatusCode())
		assert.Equal(t, "ok", response.GetMessage())
		assert.IsType(t, entities.UserEntity{}, response.GetData())
		assert.IsType(t, AuthResponse{}, response)
	})

}
