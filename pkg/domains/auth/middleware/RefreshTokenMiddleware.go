package middleware

import (
	"errors"
	"paniq-auth-user-service/pkg/contracts"
)

type RefreshTokenMiddleware struct {
	refreshToken string
}

func (r RefreshTokenMiddleware) GetRefreshToken() string {
	return r.refreshToken
}

func (r RefreshTokenMiddleware) Validate() (bool, error) {
	if r.refreshToken == "" {
		return false, errors.New("access token is not provided")
	}
	return true, nil
}

func NewRefreshTokenMiddleware(token string) contracts.IRefreshTokenMiddleware {
	return &RefreshTokenMiddleware{refreshToken: token}
}
