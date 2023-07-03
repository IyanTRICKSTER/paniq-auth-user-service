package requests

import (
	"github.com/stretchr/testify/assert"
	"paniq-auth-user-service/pkg/requests"
	"testing"
)

func TestResetPasswordConfirmation(t *testing.T) {

	t.Run("Confirmation password is similar", func(t *testing.T) {
		req := requests.ResetUserPasswordRequest{
			Password:  "1",
			CPassword: "1",
		}

		status, err := req.ValidatePassword()
		assert.True(t, status)
		assert.Nil(t, err)
	})

	t.Run("Confirmation password is not similar", func(t *testing.T) {
		req := requests.ResetUserPasswordRequest{
			Password:  "1",
			CPassword: "0",
		}

		status, err := req.ValidatePassword()
		assert.False(t, status)
		assert.NotNil(t, err)
	})
}
