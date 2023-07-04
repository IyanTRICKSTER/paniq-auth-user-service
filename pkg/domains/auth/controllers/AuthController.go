package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"paniq-auth-user-service/pkg/contracts"
	"paniq-auth-user-service/pkg/domains/auth/middleware"
	"paniq-auth-user-service/pkg/requests"
)

type AuthController struct {
	router      *gin.Engine
	authUsecase contracts.IAuthUsecase
}

func RunAuthController(router *gin.Engine, usecase contracts.IAuthUsecase) {
	ac := AuthController{
		router:      router,
		authUsecase: usecase,
	}

	router.Use(middleware.HandleCORS)
	router.POST("/api/auth/login", ac.Login)
	router.GET("/api/auth/introspect", middleware.HandleIntrospectTokenMiddleware(), ac.IntrospectToken)
	router.GET("/api/auth/refresh", middleware.HandleRefreshTokenRequestMiddleware(), ac.RefreshToken)
}

func (c AuthController) Login(ctx *gin.Context) {

	var req requests.LoginUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res := c.authUsecase.Login(ctx, req)
	if res.IsFailed() {
		ctx.JSON(http.StatusBadRequest, res.ToMapStringInterface())
		return
	}
	ctx.JSON(http.StatusOK, res.ToMapStringInterface())
}

func (c AuthController) IntrospectToken(ctx *gin.Context) {

	if introspectedJWTToken, exists := ctx.Get(middleware.IntrospectedAccessToken); exists {
		res := c.authUsecase.IntrospectToken(ctx, introspectedJWTToken.(contracts.IMiddlewareIntrospectToken))
		if res.IsFailed() {
			ctx.JSON(http.StatusInternalServerError, res.ToMapStringInterface())
			return
		}

		data := res.GetData().(map[string]any)

		//Set Custom Response Header. These headers are very important
		//to enable API GATEWAY forward authenticated request towards UPSTREAM SERVER
		ctx.Header("X-User-Id", data["user_id"].(string))
		ctx.Header("X-User-Role", data["role"].(string))
		ctx.Header("X-User-Permission", data["permissions"].(string))
		ctx.Header("X-User-Name", data["username"].(string))
		ctx.Header("X-User-Major", data["user_major"].(string))

		ctx.JSON(http.StatusOK, res.ToMapStringInterface())
		return
	}
}

func (c AuthController) RefreshToken(ctx *gin.Context) {
	if refreshTokenReq, exists := ctx.Get(middleware.IntrospectedRefreshToken); exists {
		res := c.authUsecase.RefreshToken(context.Background(), refreshTokenReq.(contracts.IRefreshTokenMiddleware))
		if res.IsFailed() {
			ctx.JSON(http.StatusBadRequest, res.ToMapStringInterface())
			return
		}
		ctx.JSON(http.StatusOK, res.ToMapStringInterface())
		return
	}
}
