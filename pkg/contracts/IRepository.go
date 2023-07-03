package contracts

import (
	"context"
	"paniq-auth-user-service/pkg/contracts/apiResources"
	"paniq-auth-user-service/pkg/entities"
	"paniq-auth-user-service/pkg/response"
)

type IUserRepository interface {
	FetchAllUser(context context.Context) response.UserResponse
	FetchUserByEmail(context context.Context, email string) response.UserResponse
	FetchUserByID(context context.Context, id uint) response.UserResponse
	FetchUserACL(context context.Context, userID uint, resource apiResources.RESOURCE) response.UserResponse
	CreateUser(context context.Context, user entities.UserEntity) response.UserResponse
	CreateUsers(context context.Context, users []entities.UserEntity) response.UserResponse
	UpdateUser(context context.Context, user entities.UserEntity) response.UserResponse
}
