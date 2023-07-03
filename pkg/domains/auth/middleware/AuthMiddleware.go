package middleware

import (
	"errors"
	"paniq-auth-user-service/pkg/contracts"
	"paniq-auth-user-service/pkg/contracts/permissionCodes"
	"strings"
)

type AuthenticatedRequestMiddleware struct {
	userID          uint
	username        string
	role            string
	permissionCodes string
}

func (a *AuthenticatedRequestMiddleware) GetUserID() uint {
	return a.userID
}

func (a *AuthenticatedRequestMiddleware) GetUsername() string {
	return a.username
}

func (a *AuthenticatedRequestMiddleware) GetRole() string {
	return a.role
}

func (a *AuthenticatedRequestMiddleware) GetPermissionCodes() string {
	return a.permissionCodes
}

func (a *AuthenticatedRequestMiddleware) HasRole(role string) bool {
	return strings.Contains(a.role, role)
}

func (a *AuthenticatedRequestMiddleware) HasPermission(code permissionCodes.PermissionCode) bool {
	return strings.Contains(a.permissionCodes, string(code))
}

func (a *AuthenticatedRequestMiddleware) Validate() (bool, error) {
	if a.userID == 0 || a.username == "" || a.role == "" {
		return false, errors.New("invalid request headers")
	}
	return true, nil
}

func NewAuthenticatedRequestMiddleware(
	userID uint,
	username string,
	role string,
	permissionCodes string,
) contracts.IAuthenticatedRequest {
	return &AuthenticatedRequestMiddleware{
		userID:          userID,
		username:        username,
		role:            role,
		permissionCodes: permissionCodes,
	}
}
