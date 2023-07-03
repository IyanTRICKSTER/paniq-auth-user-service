package usecases

import (
	"context"
	"os"
	"paniq-auth-user-service/pkg/contracts"
	"paniq-auth-user-service/pkg/contracts/statusCodes"
	"paniq-auth-user-service/pkg/entities"
	"paniq-auth-user-service/pkg/requests"
	"paniq-auth-user-service/pkg/response"
	"strconv"
)

type AuthUsecase struct {
	userRepo     contracts.IUserRepository
	jwtService   contracts.IJWTService
	hashFunction contracts.IHash
}

func (a *AuthUsecase) Login(ctx context.Context, request requests.LoginUserRequest) response.AuthResponse {

	//Fetch user
	res := a.userRepo.FetchUserByEmail(ctx, request.Email)

	if res.IsFailed() {
		if res.ErrorIs(statusCodes.ModelNotFound) {
			return response.New(response.AuthResponse{},
				false,
				statusCodes.ModelNotFound,
				"email atau password salah.",
				nil)
		}

		return response.New(response.AuthResponse{},
			false,
			statusCodes.Error,
			res.GetMessage(),
			nil)
	}

	user := res.GetData().(entities.UserEntity)

	//password checking
	check, err := a.hashFunction.HashCheck(user.Password, request.Password)
	if err != nil || check == false {
		return response.New(response.AuthResponse{},
			false,
			statusCodes.Error,
			"email atau password salah.",
			nil)
	}

	//credentials valid
	//generate access token
	accessTokenlifeSpan, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_TOKEN_LIFESPAN"))
	accessToken, _ := a.jwtService.GenerateToken(user.ID, accessTokenlifeSpan, os.Getenv("JWT_ACCESS_TOKEN_SECRET"))

	refreshTokenlifeSpan, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_TOKEN_LIFESPAN"))
	refreshToken, _ := a.jwtService.GenerateToken(user.ID, refreshTokenlifeSpan, os.Getenv("JWT_REFRESH_TOKEN_SECRET"))

	return response.New(response.AuthResponse{},
		true,
		statusCodes.Success,
		"login sukses.",
		map[string]any{
			"token_type":    "Bearer",
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	)

}

func (a *AuthUsecase) Logout(ctx context.Context) response.AuthResponse {
	//TODO implement me
	panic("implement me")
}

func (a *AuthUsecase) RefreshToken(ctx context.Context, introspect contracts.IRefreshTokenMiddleware) response.AuthResponse {

	payload, err := a.jwtService.ExtractPayloadFromToken(introspect.GetRefreshToken(), os.Getenv("JWT_REFRESH_TOKEN_SECRET"))
	if err == nil {
		//generate access token
		lifeSpan, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_TOKEN_LIFESPAN"))
		accessToken, _ := a.jwtService.GenerateToken(uint(payload["user_id"].(float64)), lifeSpan, os.Getenv("JWT_ACCESS_TOKEN_SECRET"))

		return response.New(response.AuthResponse{},
			true,
			statusCodes.Success,
			"refresh token success.",
			map[string]any{
				"token_type":   "Bearer",
				"access_token": accessToken,
			},
		)
	}

	return response.New(response.AuthResponse{},
		false,
		statusCodes.ErrExtractJWTToken,
		"refresh token is invalid, "+err.Error(),
		nil,
	)

}

func (a *AuthUsecase) IntrospectToken(ctx context.Context, introspect contracts.IMiddlewareIntrospectToken) response.AuthResponse {

	payload, err := a.jwtService.ExtractPayloadFromToken(
		introspect.GetAccessToken(),
		os.Getenv("JWT_ACCESS_TOKEN_SECRET"))

	if err != nil {
		return response.New(response.AuthResponse{},
			false,
			statusCodes.ErrExtractJWTToken,
			err.Error(),
			nil)
	}

	res := a.userRepo.FetchUserACL(
		ctx,
		uint(payload["user_id"].(float64)),
		introspect.GetTargetResourceName(),
	)

	if res.IsFailed() {
		return response.New(response.AuthResponse{},
			false,
			statusCodes.Error,
			res.GetMessage(),
			nil)
	}

	user := res.GetData().(entities.UserEntity)

	permissionCode := ""
	for _, v := range user.Role.Permissions {
		permissionCode += string(v.Code)
	}

	return response.New(
		response.AuthResponse{},
		true,
		statusCodes.Success,
		"success",
		map[string]any{
			"user_id":     strconv.Itoa(int(user.ID)),
			"username":    user.Username,
			"role":        user.Role.Name,
			"permissions": permissionCode,
			"user_major":  user.Major,
		})
}

func NewAuthUsecase(
	userRepo contracts.IUserRepository,
	jwtService contracts.IJWTService,
	hashFunction contracts.IHash) contracts.IAuthUsecase {
	return &AuthUsecase{
		userRepo:     userRepo,
		jwtService:   jwtService,
		hashFunction: hashFunction,
	}
}
