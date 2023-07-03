package repositories

import (
	"context"
	"github.com/stretchr/testify/mock"
	"paniq-auth-user-service/pkg/contracts/apiResources"
	"paniq-auth-user-service/pkg/entities"
	"paniq-auth-user-service/pkg/response"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (u *UserRepositoryMock) FetchUserACL(context context.Context, userID uint, resource apiResources.RESOURCE) response.UserResponse {
	args := u.Called(userID, resource)
	return args.Get(0).(response.UserResponse)
}

func (u *UserRepositoryMock) FetchAllUser(context context.Context) response.UserResponse {
	args := u.Called(context)
	return args.Get(0).(response.UserResponse)
}

func (u *UserRepositoryMock) FetchUserByID(context context.Context, id uint) response.UserResponse {
	args := u.Called(id)
	return args.Get(0).(response.UserResponse)
}

func (u *UserRepositoryMock) CreateUser(context context.Context, user entities.UserEntity) response.UserResponse {
	args := u.Called(user)
	return args.Get(0).(response.UserResponse)
}

func (u *UserRepositoryMock) CreateUsers(context context.Context, users []entities.UserEntity) response.UserResponse {
	//TODO implement me
	panic("implement me")
}

func (u *UserRepositoryMock) UpdateUser(context context.Context, user entities.UserEntity) response.UserResponse {
	args := u.Called(context, user)
	return args.Get(0).(response.UserResponse)
}

func (u *UserRepositoryMock) FetchUserByEmail(ctx context.Context, email string) response.UserResponse {
	args := u.Called(email)
	return args.Get(0).(response.UserResponse)
}
