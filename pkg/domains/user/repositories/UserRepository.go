package repositories

import (
	"context"
	"gorm.io/gorm"
	"paniq-auth-user-service/pkg/contracts"
	"paniq-auth-user-service/pkg/contracts/apiResources"
	"paniq-auth-user-service/pkg/contracts/statusCodes"
	"paniq-auth-user-service/pkg/database"
	"paniq-auth-user-service/pkg/entities"
	"paniq-auth-user-service/pkg/response"
)

type UserRepository struct {
	db database.Database
}

func (u UserRepository) FetchUserACL(context context.Context, userID uint, resource apiResources.RESOURCE) response.UserResponse {
	var user entities.UserEntity
	if err := u.db.GetConnection().WithContext(context).Preload("Role.Permissions", func(tx *gorm.DB) *gorm.DB {
		return tx.Select("id", "name", "code", "resource").Where("resource = ?", resource)
	}).First(&user, userID).Error; err != nil {
		return response.New(response.UserResponse{},
			false,
			statusCodes.Error,
			"",
			nil)
	}
	return response.New(response.UserResponse{},
		true,
		statusCodes.Success,
		"ok",
		user)
}

func (u UserRepository) FetchAllUser(ctx context.Context) response.UserResponse {
	var users []entities.UserEntity
	u.db.GetConnection().Find(&users)
	return response.New(
		response.UserResponse{},
		true,
		statusCodes.Success,
		"",
		users)
}

func (u UserRepository) FetchUserByEmail(ctx context.Context, email string) response.UserResponse {

	var user entities.UserEntity
	err := u.db.GetConnection().WithContext(ctx).Where("email = ?", email).First(&user).Error

	if err != nil {
		//handle record not found
		if err == gorm.ErrRecordNotFound {
			return response.New(
				response.UserResponse{},
				false,
				statusCodes.ModelNotFound,
				err.Error(),
				nil)
		}

		//handle other errors
		return response.New(
			response.UserResponse{},
			false,
			statusCodes.Error,
			err.Error(),
			nil)
	}

	return response.New(
		response.UserResponse{},
		true,
		statusCodes.Success,
		"ok",
		user)

}

func (u UserRepository) FetchUserByID(ctx context.Context, ID uint) response.UserResponse {

	var user entities.UserEntity
	err := u.db.GetConnection().WithContext(ctx).Where("id = ?", ID).First(&user).Error

	if err != nil {
		//handle record not found
		if err == gorm.ErrRecordNotFound {
			return response.New(
				response.UserResponse{},
				false,
				statusCodes.ModelNotFound,
				err.Error(),
				nil)
		}

		//handle other errors
		return response.New(
			response.UserResponse{},
			false,
			statusCodes.Error,
			err.Error(),
			nil)
	}

	return response.New(
		response.UserResponse{},
		true,
		statusCodes.Success,
		"ok",
		user)

}

func (u UserRepository) CreateUser(context context.Context, user entities.UserEntity) response.UserResponse {

	err := u.db.GetConnection().WithContext(context).Create(&user).Error

	if err != nil {
		//handle other errors
		return response.New(
			response.UserResponse{},
			false,
			statusCodes.Error,
			err.Error(),
			nil)
	}

	return response.New(
		response.UserResponse{},
		true,
		statusCodes.Success,
		"",
		user)
}

func (u UserRepository) CreateUsers(context context.Context, users []entities.UserEntity) response.UserResponse {

	err := u.db.GetConnection().WithContext(context).CreateInBatches(&users, 1000).Error
	if err != nil {
		//handle other errors
		return response.New(
			response.UserResponse{},
			false,
			statusCodes.Error,
			err.Error(),
			nil)
	}

	return response.New(
		response.UserResponse{},
		true,
		statusCodes.Success,
		"",
		nil)

}

func (u UserRepository) UpdateUser(context context.Context, user entities.UserEntity) response.UserResponse {

	err := u.db.GetConnection().WithContext(context).Save(&user).Error
	if err != nil {
		//handle other errors
		return response.New(
			response.UserResponse{},
			false,
			statusCodes.Error,
			err.Error(),
			nil)
	}

	return response.New(
		response.UserResponse{},
		true,
		statusCodes.Success,
		"",
		user)
}

func NewUserRepo(db database.Database) contracts.IUserRepository {
	return &UserRepository{db: db}
}
