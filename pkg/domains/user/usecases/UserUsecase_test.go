package usecases

import (
	"context"
	"errors"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"log"
	"os"
	"paniq-auth-user-service/pkg/contracts"
	"paniq-auth-user-service/pkg/contracts/statusCodes"
	"paniq-auth-user-service/pkg/domains/auth/middleware"
	"paniq-auth-user-service/pkg/domains/notification"
	"paniq-auth-user-service/pkg/domains/user/repositories"
	"paniq-auth-user-service/pkg/entities"
	"paniq-auth-user-service/pkg/requests"
	"paniq-auth-user-service/pkg/response"
	bcryptUtils "paniq-auth-user-service/pkg/utils/bcrypt"
	jwtUtils "paniq-auth-user-service/pkg/utils/jwt"
	"testing"
	"time"
)

var userRepoMock repositories.UserRepositoryMock
var userUsecase contracts.IUserUsecase
var hashFunctionMock bcryptUtils.HashFunctionMock
var jwtSvc contracts.IJWTService
var jwtSvcMock jwtUtils.JWTServiceMock
var notificationSvc contracts.INotificationService

func init() {

	//Load .env file
	if err := godotenv.Load("../../../../.env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	userRepoMock = repositories.UserRepositoryMock{}
	hashFunctionMock = bcryptUtils.HashFunctionMock{}
	jwtSvc = jwtUtils.New()
	jwtSvcMock = jwtUtils.JWTServiceMock{}
	notificationSvc = notification.NewUsecase(
		os.Getenv("SMTP_HOST"),
		os.Getenv("SMTP_PORT"),
		os.Getenv("SMTP_USERNAME"),
		os.Getenv("SMTP_PASSWORD"),
	)
	userUsecase = NewUserUsecase(&userRepoMock, &hashFunctionMock, &jwtSvcMock, notificationSvc)
}

func TestUserUsecase_FetchAllUser(t *testing.T) {

	t.Run("Fetch failed, doesn't have permission", func(t *testing.T) {

		authCtx := context.WithValue(
			context.Background(),
			middleware.AuthenticatedRequest,
			middleware.NewAuthenticatedRequestMiddleware(10, "iyan", "member", "ud"))

		res := userUsecase.FetchAllUser(authCtx)
		assert.False(t, res.GetStatus())
		assert.Equal(t, statusCodes.ErrForbidden, res.GetStatusCode())
		assert.Nil(t, res.GetData())

	})

	t.Run("Fetch success, has the permission", func(t *testing.T) {
		authCtx := context.WithValue(
			context.Background(),
			middleware.AuthenticatedRequest,
			middleware.NewAuthenticatedRequestMiddleware(10, "iyan", "member", "l"))

		//Mock Fetch All repo
		userRepoMock.On("FetchAllUser", authCtx).Return(
			response.New(
				response.UserResponse{},
				true, statusCodes.Success, "", []entities.UserEntity{
					{
						ID:         1,
						RoleID:     0,
						Role:       entities.RoleEntity{},
						Username:   "",
						Password:   "",
						Email:      "",
						Avatar:     "",
						NIM:        nil,
						NIP:        nil,
						Major:      "",
						ResetToken: "",
						CreatedAt:  time.Time{},
						UpdatedAt:  time.Time{},
						DeletedAt:  gorm.DeletedAt{},
					},
				}),
		)

		res := userUsecase.FetchAllUser(authCtx)
		assert.True(t, res.GetStatus())
		assert.Equal(t, statusCodes.Success, res.GetStatusCode())
		assert.NotNil(t, res.GetData())
	})
}

func TestUserUsecase_FetchUserByID(t *testing.T) {

	t.Run("Fetch success, has the permission", func(t *testing.T) {
		userID := uint(121)
		authCtx := context.WithValue(
			context.Background(),
			middleware.AuthenticatedRequest,
			middleware.NewAuthenticatedRequestMiddleware(10, "iyan", "member", "l"))

		//Mock Fetch All repo
		userRepoMock.On("FetchUserByID", userID).Return(
			response.New(
				response.UserResponse{},
				true, statusCodes.Success, "", entities.UserEntity{
					ID:         userID,
					RoleID:     0,
					Role:       entities.RoleEntity{},
					Username:   "",
					Password:   "",
					Email:      "",
					Avatar:     "",
					NIM:        nil,
					NIP:        nil,
					Major:      "",
					ResetToken: "",
					CreatedAt:  time.Time{},
					UpdatedAt:  time.Time{},
					DeletedAt:  gorm.DeletedAt{},
				}),
		)

		res := userUsecase.FetchUserByID(authCtx, userID)
		assert.True(t, res.GetStatus())
		assert.Equal(t, statusCodes.Success, res.GetStatusCode())
		assert.NotNil(t, res.GetData())
	})

}

func TestUserUsecase_CreateAccount(t *testing.T) {

	t.Run("create failed, doesn't have permission", func(t *testing.T) {

		authCtx := context.WithValue(
			context.Background(),
			middleware.AuthenticatedRequest,
			middleware.NewAuthenticatedRequestMiddleware(10, "iyan", "member", "lud"))

		req := requests.RegisterUserRequest{
			Username: "",
			Password: "",
			Email:    "",
			NIM:      "",
			NIP:      "",
			Major:    "",
		}

		res := userUsecase.CreateAccount(authCtx, req)
		assert.False(t, res.GetStatus())
		assert.Equal(t, statusCodes.ErrForbidden, res.GetStatusCode())
	})

	t.Run("create failed email already used", func(t *testing.T) {

		authCtx := context.WithValue(
			context.Background(),
			middleware.AuthenticatedRequest,
			middleware.NewAuthenticatedRequestMiddleware(10, "iyan", "admin", "c"))

		req := requests.RegisterUserRequest{
			Username: "",
			Password: "",
			Email:    "iyan@gmail.com",
			NIM:      "",
			NIP:      "",
			Major:    "",
		}

		//Mock repo create user
		userRes := response.New(response.UserResponse{}, true, statusCodes.Success, "", nil)
		userRepoMock.On("FetchUserByEmail", req.Email).Return(userRes)

		res := userUsecase.CreateAccount(authCtx, req)
		assert.Nil(t, res.GetData())
		assert.False(t, res.GetStatus())
		assert.Equal(t, statusCodes.ErrDuplicatedModel, res.GetStatusCode())
	})

	t.Run("create failed", func(t *testing.T) {

		authCtx := context.WithValue(
			context.Background(),
			middleware.AuthenticatedRequest,
			middleware.NewAuthenticatedRequestMiddleware(10, "iyan", "admin", "c"))

		req := requests.RegisterUserRequest{
			Username: "iyan",
			Password: "payload",
			Email:    "iyanputra@gmail.com",
			NIM:      "",
			NIP:      "",
			Major:    "cs digri",
		}

		//Mock Hash Function
		hashFunctionMock.On("Hash", "payload").Return("payload")

		//Mock repo FetchUserByEmail
		userRes := response.New(response.UserResponse{}, false, statusCodes.Error, "", nil)
		userRepoMock.On("FetchUserByEmail", req.Email).Return(userRes)

		//Mock repo create user
		user := entities.UserEntity{
			RoleID:     uint(2), //add default role 'member'
			Username:   req.Username,
			Password:   hashFunctionMock.Hash(req.Password),
			Email:      req.Email,
			Avatar:     "https://picsum.photos/30/30", // add random image
			NIM:        &req.NIM,
			NIP:        &req.NIP,
			Major:      req.Major,
			ResetToken: "",
			CreatedAt:  time.Time{},
			UpdatedAt:  time.Time{},
		}

		userRes2 := response.New(response.UserResponse{}, false, statusCodes.Error, "", nil)
		userRepoMock.On("CreateUser", user).Return(userRes2)

		res := userUsecase.CreateAccount(authCtx, req)
		assert.Nil(t, res.GetData())
		assert.False(t, res.GetStatus())
		assert.Equal(t, statusCodes.Error, res.GetStatusCode())
	})

	t.Run("create success", func(t *testing.T) {

		authCtx := context.WithValue(
			context.Background(),
			middleware.AuthenticatedRequest,
			middleware.NewAuthenticatedRequestMiddleware(10, "iyan", "admin", "c"))

		req := requests.RegisterUserRequest{
			Username: "iyan",
			Password: "payload",
			Email:    "iyanpratama@gmail.com",
			NIM:      "",
			NIP:      "",
			Major:    "cs digri",
		}

		//Mock Hash Function
		hashFunctionMock.On("Hash", "payload").Return("payload")

		//Mock repo FetchUserByEmail
		userRes := response.New(response.UserResponse{}, false, statusCodes.Error, "", nil)
		userRepoMock.On("FetchUserByEmail", req.Email).Return(userRes)

		//Mock repo create user
		user := entities.UserEntity{
			RoleID:     uint(2), //add default role 'member'
			Username:   req.Username,
			Password:   hashFunctionMock.Hash(req.Password),
			Email:      req.Email,
			Avatar:     "https://picsum.photos/30/30", // add random image
			NIM:        &req.NIM,
			NIP:        &req.NIP,
			Major:      req.Major,
			ResetToken: "",
			CreatedAt:  time.Time{},
			UpdatedAt:  time.Time{},
		}

		userRes2 := response.New(response.UserResponse{}, true, statusCodes.Success, "", user)
		userRepoMock.On("CreateUser", user).Return(userRes2)

		res := userUsecase.CreateAccount(authCtx, req)
		assert.Nil(t, res.GetData())
		assert.True(t, res.GetStatus())
		assert.Equal(t, statusCodes.Success, res.GetStatusCode())
	})
}

func TestUserUsecase_UpdateAccount(t *testing.T) {

	t.Run("update failed, user doesn't have permission", func(t *testing.T) {

		authCtx := context.WithValue(
			context.Background(),
			middleware.AuthenticatedRequest,
			middleware.NewAuthenticatedRequestMiddleware(10, "iyan", "member", ""))

		req := requests.UpdateUserRequest{
			Avatar: "https://picsum.photos/40/40",
			NIM:    "11210910000004",
			NIP:    "11210910000044",
			Major:  "Physics",
		}

		res := userUsecase.UpdateAccount(authCtx, 10, req)
		assert.False(t, res.GetStatus())
		assert.Equal(t, statusCodes.ErrForbidden, res.GetStatusCode())
		assert.Nil(t, res.GetData())
	})

	t.Run("update failed user is not exist", func(t *testing.T) {

		var userID uint
		userID = 10

		authCtx := context.WithValue(
			context.Background(),
			middleware.AuthenticatedRequest,
			middleware.NewAuthenticatedRequestMiddleware(userID, "iyan", "member", "u"))

		req := requests.UpdateUserRequest{
			Avatar: "https://picsum.photos/40/40",
			NIM:    "11210910000004",
			NIP:    "11210910000044",
			Major:  "Physics",
		}

		//Mock Repo fetch user
		userRepoMock.On("FetchUserByID", userID).
			Return(response.New(response.UserResponse{}, false, statusCodes.ModelNotFound, "", nil))

		res := userUsecase.UpdateAccount(authCtx, userID, req)
		assert.False(t, res.GetStatus())
		assert.Equal(t, statusCodes.ModelNotFound, res.GetStatusCode())
		assert.Nil(t, res.GetData())
	})

	t.Run("update failed user id is not match", func(t *testing.T) {

		var userID uint
		userID = 11

		authCtx := context.WithValue(
			context.Background(),
			middleware.AuthenticatedRequest,
			middleware.NewAuthenticatedRequestMiddleware(userID, "iyan", "member", "u"))

		req := requests.UpdateUserRequest{
			Avatar: "https://picsum.photos/40/40",
			NIM:    "11210910000004",
			NIP:    "11210910000044",
			Major:  "Physics",
		}

		//Mock Repo fetch user
		userRepoMock.On("FetchUserByID", userID).
			Return(response.New(response.UserResponse{}, true,
				statusCodes.ModelNotFound,
				"",
				entities.UserEntity{ID: 12}))

		res := userUsecase.UpdateAccount(authCtx, userID, req)
		assert.False(t, res.GetStatus())
		assert.Equal(t, statusCodes.ErrForbidden, res.GetStatusCode())
		assert.Nil(t, res.GetData())
	})

	t.Run("update failed repo error", func(t *testing.T) {

		var userID uint
		userID = 12

		authCtx := context.WithValue(
			context.Background(),
			middleware.AuthenticatedRequest,
			middleware.NewAuthenticatedRequestMiddleware(userID, "iyan", "member", "u"))

		req := requests.UpdateUserRequest{
			Avatar: "https://picsum.photos/40/40",
			NIM:    "11210910000004",
			NIP:    "11210910000044",
			Major:  "Physics",
		}

		//Mock Repo fetch user
		userRepoMock.On("FetchUserByID", userID).
			Return(response.New(response.UserResponse{}, true,
				statusCodes.Success,
				"",
				entities.UserEntity{ID: userID}))

		//Mock repo Update user
		var user entities.UserEntity
		user.ID = userID
		user.Avatar = req.Avatar
		user.NIM = &req.NIM
		user.NIP = &req.NIP
		user.Major = req.Major

		userRepoMock.On("UpdateUser", authCtx, user).
			Return(response.New(response.UserResponse{}, false,
				statusCodes.Error,
				"",
				nil))

		// do update
		res := userUsecase.UpdateAccount(authCtx, userID, req)
		log.Println(res.GetMessage())
		assert.False(t, res.GetStatus())
		assert.Equal(t, statusCodes.Error, res.GetStatusCode())
		assert.Nil(t, res.GetData())
	})

	t.Run("update success", func(t *testing.T) {

		var userID uint
		userID = 13

		authCtx := context.WithValue(
			context.Background(),
			middleware.AuthenticatedRequest,
			middleware.NewAuthenticatedRequestMiddleware(userID, "iyan", "member", "u"))

		req := requests.UpdateUserRequest{
			Avatar: "https://picsum.photos/40/40",
			NIM:    "11210910000004",
			NIP:    "11210910000044",
			Major:  "Physics",
		}

		//Mock Repo fetch user
		userRepoMock.On("FetchUserByID", userID).
			Return(response.New(response.UserResponse{}, true,
				statusCodes.Success,
				"",
				entities.UserEntity{ID: userID}))

		//Mock repo Update user
		var user entities.UserEntity
		user.ID = userID
		user.Avatar = req.Avatar
		user.NIM = &req.NIM
		user.NIP = &req.NIP
		user.Major = req.Major

		userRepoMock.On("UpdateUser", authCtx, user).
			Return(response.New(response.UserResponse{}, true,
				statusCodes.Success,
				"",
				user))

		// do update
		res := userUsecase.UpdateAccount(authCtx, userID, req)
		log.Println(res.GetMessage())
		assert.True(t, res.GetStatus())
		assert.Equal(t, statusCodes.Success, res.GetStatusCode())
		assert.NotNil(t, res.GetData())
	})
}

func TestUserUsecase_ChangePassword(t *testing.T) {

	t.Run("fetch user by email fail model not found", func(t *testing.T) {

		ctx := context.Background()
		req := requests.ChangeUserPasswordRequest{Email: "iyanpratama2001@gmail.com"}

		userRepoMock.On("FetchUserByEmail", req.Email).Return(
			response.New(response.UserResponse{}, false,
				statusCodes.ModelNotFound,
				"",
				nil),
		)

		res := userUsecase.ChangePassword(ctx, req)
		assert.False(t, res.GetStatus())
		assert.Equal(t, statusCodes.ModelNotFound, res.GetStatusCode())
		assert.Nil(t, res.GetData())
	})

	t.Run("fetch user by email fail, other error", func(t *testing.T) {

		ctx := context.Background()
		req := requests.ChangeUserPasswordRequest{Email: "iyanpratama2002@gmail.com"}

		userRepoMock.On("FetchUserByEmail", req.Email).Return(
			response.New(response.UserResponse{}, false,
				statusCodes.Error,
				"",
				nil),
		)

		res := userUsecase.ChangePassword(ctx, req)
		assert.False(t, res.GetStatus())
		assert.Equal(t, statusCodes.Error, res.GetStatusCode())
		assert.Nil(t, res.GetData())
	})

	t.Run("failed convert lifespan token", func(t *testing.T) {
		ctx := context.Background()
		req := requests.ChangeUserPasswordRequest{Email: "iyanpratama2003@gmail.com"}

		userRepoMock.On("FetchUserByEmail", req.Email).Return(
			response.New(response.UserResponse{}, true,
				statusCodes.Success,
				"",
				entities.UserEntity{}),
		)

		os.Setenv("JWT_PASS_RESET_TOKEN_LIFESPAN", "abc") // this will create error when
		// strconv.Atoi try to convert the "abc"

		res := userUsecase.ChangePassword(ctx, req)
		assert.False(t, res.GetStatus())
		assert.Equal(t, statusCodes.Error, res.GetStatusCode())
		assert.Nil(t, res.GetData())
	})

	t.Run("generate token error", func(t *testing.T) {
		userID := uint(10)
		ctx := context.Background()
		req := requests.ChangeUserPasswordRequest{Email: "iyanpratama2004@gmail.com"}

		userRepoMock.On("FetchUserByEmail", req.Email).Return(
			response.New(response.UserResponse{}, true,
				statusCodes.Success,
				"",
				entities.UserEntity{ID: userID}),
		)

		os.Setenv("JWT_PASS_RESET_TOKEN_LIFESPAN", "3600") // this will create error when
		//strconv.Atoi try to convert the "abc"

		jwtSvcMock.On(
			"GenerateToken",
			userID,
			3600,
			os.Getenv("JWT_PASS_RESET_TOKEN_SECRET")).Return("", errors.New("generate token error"))

		res := userUsecase.ChangePassword(ctx, req)
		assert.False(t, res.GetStatus())
		assert.Equal(t, statusCodes.Error, res.GetStatusCode())
		assert.Nil(t, res.GetData())

	})

	t.Run("generate token success", func(t *testing.T) {
		userID := uint(11)
		ctx := context.Background()
		req := requests.ChangeUserPasswordRequest{Email: "iyanpratama2005@gmail.com"}
		user := entities.UserEntity{ID: userID}

		userRepoMock.On("FetchUserByEmail", req.Email).Return(
			response.New(response.UserResponse{}, true,
				statusCodes.Success,
				"",
				user),
		)

		os.Setenv("JWT_PASS_RESET_TOKEN_LIFESPAN", "3600") // this will create error when
		//strconv.Atoi try to convert the "abc"

		jwtSvcMock.On(
			"GenerateToken",
			userID,
			3600,
			os.Getenv("JWT_PASS_RESET_TOKEN_SECRET")).Return("ini token", nil)

		user.ResetToken = "ini token"
		userRepoMock.On("UpdateUser", ctx, user).Return(
			response.New(response.UserResponse{}, true,
				statusCodes.Success,
				"",
				nil),
		)

		res := userUsecase.ChangePassword(ctx, req)
		assert.True(t, res.GetStatus())
		assert.Equal(t, statusCodes.Success, res.GetStatusCode())
		assert.Nil(t, res.GetData())

	})
}

func TestUserUsecase_ResetPassword(t *testing.T) {

	t.Run("validate password fail, password confirmation is different", func(t *testing.T) {
		ctx := context.Background()
		req := requests.ResetUserPasswordRequest{
			ResetToken: "ini token",
			Password:   "password1",
			CPassword:  "password2",
		}
		res := userUsecase.ResetPassword(ctx, req)
		assert.False(t, res.GetStatus())
		assert.Equal(t, statusCodes.ErrResetPasswordNotMatch, res.GetStatusCode())
		assert.Nil(t, res.GetData())
	})

	t.Run("extract token failed", func(t *testing.T) {
		ctx := context.Background()
		req := requests.ResetUserPasswordRequest{
			ResetToken: "ini token",
			Password:   "password1",
			CPassword:  "password1",
		}

		//Mock JWT service
		jwtSvcMock.On("ExtractPayloadFromToken", req.ResetToken, os.Getenv("JWT_PASS_RESET_TOKEN_SECRET")).
			Return(map[string]interface{}{}, errors.New("extract token failed"))

		res := userUsecase.ResetPassword(ctx, req)
		assert.False(t, res.GetStatus())
		assert.Equal(t, statusCodes.ErrExtractJWTToken, res.GetStatusCode())
		assert.Nil(t, res.GetData())
	})

	t.Run("fetch user failed model not found.", func(t *testing.T) {
		userID := uint(125)
		ctx := context.Background()
		req := requests.ResetUserPasswordRequest{
			ResetToken: "ini token 1",
			Password:   "password1",
			CPassword:  "password1",
		}

		//Mock JWT service
		jwtSvcMock.On("ExtractPayloadFromToken", req.ResetToken, os.Getenv("JWT_PASS_RESET_TOKEN_SECRET")).
			Return(map[string]interface{}{"user_id": float64(userID)}, nil)

		//Mock Fetch User Repo
		userRepoMock.On("FetchUserByID", userID).
			Return(
				response.New(response.UserResponse{}, false,
					statusCodes.ModelNotFound,
					"",
					nil),
			)

		res := userUsecase.ResetPassword(ctx, req)
		assert.False(t, res.GetStatus())
		assert.Equal(t, statusCodes.ModelNotFound, res.GetStatusCode())
		assert.Nil(t, res.GetData())
	})

	t.Run("fetch user failed, other error", func(t *testing.T) {
		userID := uint(133)
		ctx := context.Background()
		req := requests.ResetUserPasswordRequest{
			ResetToken: "ini token 2",
			Password:   "password1",
			CPassword:  "password1",
		}

		//Mock JWT service
		jwtSvcMock.On("ExtractPayloadFromToken", req.ResetToken, os.Getenv("JWT_PASS_RESET_TOKEN_SECRET")).
			Return(map[string]interface{}{"user_id": float64(userID)}, nil)

		//Mock Fetch User Repo
		userRepoMock.On("FetchUserByID", userID).
			Return(
				response.New(response.UserResponse{}, false,
					statusCodes.Error,
					"",
					nil),
			)

		res := userUsecase.ResetPassword(ctx, req)
		assert.False(t, res.GetStatus())
		assert.Equal(t, statusCodes.Error, res.GetStatusCode())
		assert.Nil(t, res.GetData())
	})

	t.Run("save new updated password failed", func(t *testing.T) {
		userID := uint(140)
		user := entities.UserEntity{ID: userID}
		ctx := context.Background()
		req := requests.ResetUserPasswordRequest{
			ResetToken: "ini token 112",
			Password:   "password_10",
			CPassword:  "password_10",
		}

		//Mock JWT service
		jwtSvcMock.On("ExtractPayloadFromToken", req.ResetToken, os.Getenv("JWT_PASS_RESET_TOKEN_SECRET")).
			Return(map[string]interface{}{"user_id": float64(userID)}, nil)

		//Mock Fetch User Repo
		userRepoMock.On("FetchUserByID", userID).
			Return(
				response.New(response.UserResponse{}, true,
					statusCodes.Success,
					"",
					user),
			)

		//Mock Hash function
		hashedPW := "hashedwkwk11"
		hashFunctionMock.On("Hash", req.Password).Return(hashedPW)

		user.Password = hashedPW
		//Mock Fetch User Repo
		userRepoMock.On("UpdateUser", ctx, user).
			Return(
				response.New(response.UserResponse{}, false,
					statusCodes.Error,
					"",
					nil),
			)

		res := userUsecase.ResetPassword(ctx, req)
		assert.False(t, res.GetStatus())
		assert.Equal(t, statusCodes.Error, res.GetStatusCode())
		assert.Nil(t, res.GetData())
	})

	t.Run("save new updated password success.", func(t *testing.T) {
		userID := uint(100)
		user := entities.UserEntity{ID: userID}
		ctx := context.Background()
		req := requests.ResetUserPasswordRequest{
			ResetToken: "ini token 22",
			Password:   "password1",
			CPassword:  "password1",
		}

		//Mock JWT service
		jwtSvcMock.On("ExtractPayloadFromToken", req.ResetToken, os.Getenv("JWT_PASS_RESET_TOKEN_SECRET")).
			Return(map[string]interface{}{"user_id": float64(userID)}, nil)

		//Mock Fetch User Repo
		userRepoMock.On("FetchUserByID", userID).
			Return(
				response.New(response.UserResponse{}, true,
					statusCodes.Success,
					"",
					user),
			)

		//Mock Hash function
		hashedPW := "hashedwkwk22"
		hashFunctionMock.On("Hash", req.Password).Return(hashedPW)

		user.Password = hashedPW
		//Mock Fetch User Repo
		userRepoMock.On("UpdateUser", ctx, user).
			Return(
				response.New(response.UserResponse{}, true,
					statusCodes.Success,
					"",
					nil),
			)

		res := userUsecase.ResetPassword(ctx, req)
		assert.True(t, res.GetStatus())
		assert.Equal(t, statusCodes.Success, res.GetStatusCode())
		assert.Nil(t, res.GetData())
	})

}
