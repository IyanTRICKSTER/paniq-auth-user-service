package usecases

import (
	"context"
	"mime/multipart"
	"os"
	"paniq-auth-user-service/pkg/contracts"
	"paniq-auth-user-service/pkg/contracts/permissionCodes"
	"paniq-auth-user-service/pkg/contracts/statusCodes"
	"paniq-auth-user-service/pkg/domains/auth/middleware"
	"paniq-auth-user-service/pkg/entities"
	"paniq-auth-user-service/pkg/requests"
	"paniq-auth-user-service/pkg/response"
	"strconv"
	"time"
)

type UserUsecase struct {
	userRepo            contracts.IUserRepository
	hashFunction        contracts.IHash
	jwtService          contracts.IJWTService
	notificationService contracts.INotificationService
}

func (u *UserUsecase) FetchAllUser(ctx context.Context) response.UserResponse {
	if !u.checkUserHasPermission(
		ctx.Value(middleware.AuthenticatedRequest).(contracts.IAuthenticatedRequest), permissionCodes.LIST) {
		return response.New(
			response.UserResponse{},
			false,
			statusCodes.ErrForbidden,
			"forbidden",
			nil)
	}
	return u.userRepo.FetchAllUser(ctx)
}

func (u *UserUsecase) FetchUserByID(ctx context.Context, userID uint) response.UserResponse {
	return u.userRepo.FetchUserByID(ctx, userID)
}

func (u *UserUsecase) CreateAccount(ctx context.Context, request requests.RegisterUserRequest) response.UserResponse {
	if !u.checkUserHasPermission(
		ctx.Value(middleware.AuthenticatedRequest).(contracts.IAuthenticatedRequest), permissionCodes.CREATE) {
		return response.New(
			response.UserResponse{},
			false,
			statusCodes.ErrForbidden,
			"forbidden",
			nil)
	}

	res := u.userRepo.FetchUserByEmail(ctx, request.Email)
	if res.GetStatus() {
		return response.New(
			response.UserResponse{},
			false,
			statusCodes.ErrDuplicatedModel,
			"credential sudah digunakan, pastikan email, nim, dan nip belum digunakan.",
			nil)
	}

	user := entities.UserEntity{
		RoleID:     uint(2), //add default role 'member'
		Username:   request.Username,
		Password:   u.hashFunction.Hash(request.Password),
		Email:      request.Email,
		Avatar:     "https://picsum.photos/30/30", // add random image
		NIM:        &request.NIM,
		NIP:        &request.NIP,
		Major:      request.Major,
		ResetToken: "",
		CreatedAt:  time.Time{},
		UpdatedAt:  time.Time{},
	}

	res = u.userRepo.CreateUser(ctx, user)
	if res.IsFailed() {
		return res
	}

	return response.New(
		response.UserResponse{},
		true,
		statusCodes.Success,
		"success.",
		nil)

}

func (u *UserUsecase) UpdateAccount(ctx context.Context, userID uint, request requests.UpdateUserRequest) response.UserResponse {
	if !u.checkUserHasPermission(
		ctx.Value(middleware.AuthenticatedRequest).(contracts.IAuthenticatedRequest), permissionCodes.UPDATE) {
		return response.New(
			response.UserResponse{},
			false,
			statusCodes.ErrForbidden,
			"forbidden, you are not authorized to access this resource",
			nil)
	}

	res := u.userRepo.FetchUserByID(ctx, userID)
	if res.IsFailed() {
		return res
	}

	user := res.GetData().(entities.UserEntity)
	if !u.isOwner(userID, user) {
		return response.New(
			response.UserResponse{},
			false,
			statusCodes.ErrForbidden,
			"forbidden",
			nil)
	}

	//sementara update avatar dengan url,
	//nanti harus dibuat storage service nya
	user.Avatar = request.Avatar
	user.NIM = &request.NIM
	user.NIP = &request.NIP
	user.Major = request.Major

	res = u.userRepo.UpdateUser(ctx, user)
	if res.IsFailed() {
		return res
	}

	res.SetMessage("update user success")
	return res
}

func (u *UserUsecase) ChangePassword(ctx context.Context, request requests.ChangeUserPasswordRequest) response.UserResponse {

	res := u.userRepo.FetchUserByEmail(ctx, request.Email)
	if res.IsFailed() {
		if res.ErrorIs(statusCodes.ModelNotFound) {
			res.SetMessage("email tidak terdaftar.")
			return res
		}
		return res
	}

	user := res.GetData().(entities.UserEntity)

	//Generate Reset Token using JWT
	tokenLifespan, err := strconv.Atoi(os.Getenv("JWT_PASS_RESET_TOKEN_LIFESPAN"))
	if err != nil {
		return response.New(response.UserResponse{},
			false, statusCodes.Error, err.Error(), nil)
	}

	token, err := u.jwtService.GenerateToken(
		user.ID,
		tokenLifespan,
		os.Getenv("JWT_PASS_RESET_TOKEN_SECRET"),
	)
	if err != nil {
		return response.New(response.UserResponse{},
			false, statusCodes.Error, err.Error(), nil)
	}

	//Notify token to user
	go u.notificationService.NotifyWithEmail(
		"paniq@corp.mail",
		request.Email,
		"Reset Password",
		"Berikut adalah reset token password anda "+token,
	)

	user.ResetToken = token
	res = u.userRepo.UpdateUser(ctx, user)
	return res
}

func (u *UserUsecase) ResetPassword(ctx context.Context, request requests.ResetUserPasswordRequest) response.UserResponse {

	status, err := request.ValidatePassword()
	if err != nil || status == false {
		return response.New(response.UserResponse{},
			false,
			statusCodes.ErrResetPasswordNotMatch,
			err.Error(),
			nil)
	}

	payload, err := u.jwtService.ExtractPayloadFromToken(
		request.ResetToken, os.Getenv("JWT_PASS_RESET_TOKEN_SECRET"))
	if err != nil {
		return response.New(response.UserResponse{},
			false,
			statusCodes.ErrExtractJWTToken,
			err.Error(),
			nil)
	}

	res := u.userRepo.FetchUserByID(ctx, uint(payload["user_id"].(float64)))
	if res.IsFailed() {
		if res.ErrorIs(statusCodes.ModelNotFound) {
			res.SetMessage("user tidak dapat ditemukan, password gagal diperbarui.")
			return res
		}
		return res
	}

	user := res.GetData().(entities.UserEntity)
	user.Password = u.hashFunction.Hash(request.Password)
	user.ResetToken = ""

	res = u.userRepo.UpdateUser(ctx, user)
	if res.IsFailed() {
		res.SetMessage("update password gagal.")
		return res
	}
	res.SetMessage("update password berhasil.")
	return res
}

func (u *UserUsecase) CreateAccounts(ctx context.Context, request requests.RegisterUserCSVRequest) response.UserResponse {
	//TODO implement me
	panic("implement me")
}

func (u *UserUsecase) DisableAccount(ctx context.Context, userID uint) response.UserResponse {
	//TODO implement me
	panic("implement me")
}

func (u *UserUsecase) readCSVFile(data multipart.FileHeader) {
	//TODO implement me
	panic("implement me")
}

func (u *UserUsecase) validateEachCSVRow(request requests.RegisterUserCSVRequest) bool {
	//TODO implement me
	panic("implement me")
}

func (u *UserUsecase) checkUserHasPermission(auth contracts.IAuthenticatedRequest, permission permissionCodes.PermissionCode) bool {
	return auth.HasPermission(permission)
}

func (u UserUsecase) isOwner(userID uint, user entities.UserEntity) bool {
	if userID != user.ID {
		return false
	}
	return true
}

func NewUserUsecase(
	userRepo contracts.IUserRepository,
	hashFunction contracts.IHash,
	jwtService contracts.IJWTService,
	notificationService contracts.INotificationService,
) contracts.IUserUsecase {
	return &UserUsecase{
		userRepo:            userRepo,
		hashFunction:        hashFunction,
		jwtService:          jwtService,
		notificationService: notificationService,
	}
}
