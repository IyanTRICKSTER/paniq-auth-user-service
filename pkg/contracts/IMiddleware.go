package contracts

import (
	"paniq-auth-user-service/pkg/contracts/apiResources"
	"paniq-auth-user-service/pkg/contracts/permissionCodes"
)

type IAuthenticatedRequest interface {
	GetUserID() uint
	GetUsername() string
	GetRole() string
	GetPermissionCodes() string
	HasRole(role string) bool
	HasPermission(code permissionCodes.PermissionCode) bool
	Validate() (bool, error)
}

type IMiddlewareIntrospectToken interface {
	GetTargetResourceName() apiResources.RESOURCE
	GetAccessToken() string
}

type IRefreshTokenMiddleware interface {
	GetRefreshToken() string
	Validate() (bool, error)
}
