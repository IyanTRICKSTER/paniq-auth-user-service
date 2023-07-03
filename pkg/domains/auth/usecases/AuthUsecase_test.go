package usecases

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"paniq-auth-user-service/pkg/contracts"
	"paniq-auth-user-service/pkg/contracts/apiResources"
	"paniq-auth-user-service/pkg/contracts/permissionCodes"
	"paniq-auth-user-service/pkg/contracts/statusCodes"
	"paniq-auth-user-service/pkg/domains/auth/middleware"
	"paniq-auth-user-service/pkg/domains/user/repositories"
	"paniq-auth-user-service/pkg/entities"
	"paniq-auth-user-service/pkg/requests"
	"paniq-auth-user-service/pkg/response"
	bcryptUtils "paniq-auth-user-service/pkg/utils/bcrypt"
	jwtUtils "paniq-auth-user-service/pkg/utils/jwt"
	"testing"
)

var authUsecase contracts.IAuthUsecase
var userRepoMock repositories.UserRepositoryMock
var jwtSvc contracts.IJWTService
var accessTokenSecret string

func init() {

	//Load .env file
	if err := godotenv.Load("../../../../.env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	userRepoMock = repositories.UserRepositoryMock{}
	jwtSvc = jwtUtils.New()
	authUsecase = NewAuthUsecase(&userRepoMock, jwtSvc, bcryptUtils.NewHashFunction())
	accessTokenSecret = os.Getenv("JWT_ACCESS_TOKEN_SECRET")

}

func TestAuthUsecase_Login(t *testing.T) {

	t.Run("Login user not found", func(t *testing.T) {

		//Mock Response
		res := response.New(response.UserResponse{}, false, statusCodes.ModelNotFound, "", nil)
		userRepoMock.On("FetchUserByEmail", "iyanpratama2002@gmail.com").Return(res)

		//Do Login
		req := requests.LoginUserRequest{
			Email:    "iyanpratama2002@gmail.com",
			Password: "iyan12345",
		}
		loginRes := authUsecase.Login(context.Background(), req)
		assert.False(t, loginRes.GetStatus())
		assert.True(t, loginRes.ErrorIs(statusCodes.ModelNotFound))
		assert.Nil(t, loginRes.GetData())
		log.Println(loginRes.GetMessage())

	})

	t.Run("Login error", func(t *testing.T) {

		//Mock Response
		res := response.New(response.UserResponse{}, false, statusCodes.Error, "", nil)
		userRepoMock.On("FetchUserByEmail", "iyan@gmail.com").Return(res)

		//Do Login
		req := requests.LoginUserRequest{
			Email:    "iyan@gmail.com",
			Password: "iyan12345",
		}
		loginRes := authUsecase.Login(context.Background(), req)
		assert.False(t, loginRes.GetStatus())
		assert.True(t, loginRes.ErrorIs(statusCodes.Error))
		assert.Nil(t, loginRes.GetData())
		log.Println(loginRes.GetMessage())

	})

	t.Run("Login user with wrong password", func(t *testing.T) {

		//Mock Response
		res := response.New(response.UserResponse{},
			true,
			statusCodes.Success,
			"",
			entities.UserEntity{
				ID:       1,
				Password: "$2a$10$sGigrufsGHOTKti.h6nZZ.LfInGt2pYyTr0GPaIUArrqKJhux2ul."},
		)
		userRepoMock.On("FetchUserByEmail", "iyanpratama2003@gmail.com").Return(res)

		//Do Login
		req := requests.LoginUserRequest{
			Email:    "iyanpratama2003@gmail.com",
			Password: "iyan123",
		}
		loginRes := authUsecase.Login(context.Background(), req)
		assert.False(t, loginRes.GetStatus())
		assert.True(t, loginRes.ErrorIs(statusCodes.Error))
		assert.Nil(t, loginRes.GetData())
		log.Println(loginRes.GetMessage())

	})

	t.Run("Login user success", func(t *testing.T) {

		//Mock Response
		res := response.New(response.UserResponse{},
			true,
			statusCodes.Success,
			"",
			entities.UserEntity{
				ID:       1,
				Password: "$2a$10$sGigrufsGHOTKti.h6nZZ.LfInGt2pYyTr0GPaIUArrqKJhux2ul."},
		)
		userRepoMock.On("FetchUserByEmail", "iyanpratama2004@gmail.com").Return(res)

		//Do Login
		req := requests.LoginUserRequest{
			Email:    "iyanpratama2004@gmail.com",
			Password: "iyan12345",
		}
		loginRes := authUsecase.Login(context.Background(), req)
		assert.True(t, loginRes.GetStatus())
		assert.NotNil(t, loginRes.GetData())
		log.Println(loginRes.GetMessage())
		log.Println(loginRes.GetData())

	})

}

func TestAuthUsecase_IntrospectToken(t *testing.T) {

	t.Run("Introspect token failed invalid jwt token", func(t *testing.T) {

		//Mock Response
		res := response.New(response.UserResponse{}, false, statusCodes.ModelNotFound, "", nil)
		userRepoMock.On("FetchUserByEmail", "iyanpratama2002@gmail.com").Return(res)

		intrRes := authUsecase.IntrospectToken(
			context.Background(),
			middleware.NewIntrospectTokenMiddleware("invalid token", apiResources.USER))

		assert.False(t, intrRes.GetStatus())
		assert.Equal(t, statusCodes.ErrExtractJWTToken, intrRes.GetStatusCode())
		assert.Nil(t, intrRes.GetData())

		//log.Println(intrRes.GetMessage())

	})

	t.Run("Introspect token failed invalid context canceled", func(t *testing.T) {

		//Mock Response
		res := response.New(response.UserResponse{}, false, statusCodes.Error, "", nil)
		userRepoMock.On("FetchUserACL", uint(2), apiResources.USER).Return(res)

		//generate jwt token
		token, err := jwtSvc.GenerateToken(2, 7200, accessTokenSecret)
		assert.Nil(t, err)

		introRes := authUsecase.IntrospectToken(
			context.Background(),
			middleware.NewIntrospectTokenMiddleware(token, apiResources.USER))

		assert.False(t, introRes.GetStatus())
		assert.Equal(t, statusCodes.Error, introRes.GetStatusCode())
		assert.Nil(t, introRes.GetData())
	})

	t.Run("Introspect token success", func(t *testing.T) {

		//Mock Response
		user := entities.UserEntity{
			ID: 3,
			Role: entities.RoleEntity{
				Name: "admin",
				Permissions: []entities.PermissionEntity{
					{
						Code: permissionCodes.LIST,
					},
				},
			},
			Username: "iyan wkwk",
			Major:    "cs digri",
		}
		res := response.New(response.UserResponse{}, true, statusCodes.Success, "", user)
		userRepoMock.On("FetchUserACL", uint(3), apiResources.USER).Return(res)

		//generate jwt token
		token, err := jwtSvc.GenerateToken(3, 7200, accessTokenSecret)
		assert.Nil(t, err)

		introRes := authUsecase.IntrospectToken(
			context.Background(),
			middleware.NewIntrospectTokenMiddleware(token, apiResources.USER))

		assert.True(t, introRes.GetStatus())
		assert.Equal(t, statusCodes.Success, introRes.GetStatusCode())
		assert.NotNil(t, introRes.GetData())
		//log.Println(introRes.GetData())
	})
}

func TestRefreshToken(t *testing.T) {

	t.Run("refresh token success", func(t *testing.T) {

		refreshToken, err := jwtSvc.GenerateToken(uint(10), 7200, os.Getenv("JWT_REFRESH_TOKEN_SECRET"))
		assert.Nil(t, err)

		res := authUsecase.RefreshToken(context.Background(), middleware.NewRefreshTokenMiddleware(refreshToken))
		accessToken := res.GetData().(map[string]interface{})
		tokenStatus := jwtSvc.ValidateToken(accessToken["access_token"].(string), os.Getenv("JWT_ACCESS_TOKEN_SECRET"))
		assert.True(t, tokenStatus)
		assert.True(t, res.GetStatus())
	})

	t.Run("refresh token success", func(t *testing.T) {

		res := authUsecase.RefreshToken(context.Background(), middleware.NewRefreshTokenMiddleware("wkwk"))

		assert.False(t, res.GetStatus())
		assert.Equal(t, statusCodes.ErrExtractJWTToken, res.GetStatusCode())
		assert.Nil(t, res.GetData())
	})
}
