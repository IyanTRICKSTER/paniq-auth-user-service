package contracts

import (
	"context"
	"paniq-auth-user-service/pkg/requests"
	"paniq-auth-user-service/pkg/response"
)

type IUserUsecase interface {
	FetchAllUser(ctx context.Context) response.UserResponse
	FetchUserByID(ctx context.Context, userID uint) response.UserResponse
	CreateAccount(ctx context.Context, request requests.RegisterUserRequest) response.UserResponse
	CreateAccounts(ctx context.Context, request requests.RegisterUserCSVRequest) response.UserResponse
	DisableAccount(ctx context.Context, userID uint) response.UserResponse
	UpdateAccount(ctx context.Context, userID uint, request requests.UpdateUserRequest) response.UserResponse
	ChangePassword(ctx context.Context, request requests.ChangeUserPasswordRequest) response.UserResponse
	ResetPassword(ctx context.Context, request requests.ResetUserPasswordRequest) response.UserResponse
	//checkUserHasPermission(resource apiResources.RESOURCE, permission permissionCodes.PermissionCode) bool
	//readCSVFile(data multipart.FileHeader)
	//validateEachCSVRow(request requests.RegisterUserCSVRequest) bool
}

type IAuthUsecase interface {
	Login(ctx context.Context, request requests.LoginUserRequest) response.AuthResponse
	Logout(ctx context.Context) response.AuthResponse
	RefreshToken(ctx context.Context, introspect IRefreshTokenMiddleware) response.AuthResponse
	IntrospectToken(ctx context.Context, introspect IMiddlewareIntrospectToken) response.AuthResponse
	//extractTargetResource(header http.Header)
}

type INotificationService interface {
	NotifyWithEmail(from string, to string, subject string, message string)
}
