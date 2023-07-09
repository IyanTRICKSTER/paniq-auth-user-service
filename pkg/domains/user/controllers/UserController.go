package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
	"paniq-auth-user-service/pkg/contracts"
	"paniq-auth-user-service/pkg/contracts/statusCodes"
	"paniq-auth-user-service/pkg/domains/auth/middleware"
	"paniq-auth-user-service/pkg/requests"
	"paniq-auth-user-service/pkg/response"
	"strconv"
)

type UserController struct {
	router      *gin.Engine
	userUsecase contracts.IUserUsecase
}

func RunUserController(router *gin.Engine, usecase contracts.IUserUsecase) {
	ac := UserController{
		router:      router,
		userUsecase: usecase,
	}

	apiGroup := router.Group("api/user")
	apiGroup.GET("list", middleware.HandleUnauthenticatedRequestMiddleware(), ac.FetchAllUser)
	apiGroup.GET("current", middleware.HandleUnauthenticatedRequestMiddleware(), ac.Me)
	apiGroup.GET(":id", middleware.HandleUnauthenticatedRequestMiddleware(), ac.FetchUserByID)
	apiGroup.POST("register", middleware.HandleUnauthenticatedRequestMiddleware(), ac.CreateAccount)
	apiGroup.PATCH("update/:id", middleware.HandleUnauthenticatedRequestMiddleware(), ac.UpdateAccount)
	apiGroup.POST("change-password", ac.ChangePassword)
	apiGroup.POST("reset-password", ac.ResetPassword)
}

func extractValidationErr(err error) ([]string, bool) {
	if errs, ok := err.(validator.ValidationErrors); ok {
		var msg []string
		for _, e := range errs {
			fieldName := e.Field()
			errorMsg := e.Tag() + " validation failed for " + fieldName
			msg = append(msg, errorMsg)
		}

		return msg, true
	}
	return []string{}, false
}

func (u UserController) FetchAllUser(c *gin.Context) {
	res := u.userUsecase.FetchAllUser(c)
	if res.IsFailed() {
		if res.ErrorIs(statusCodes.ErrForbidden) {
			c.JSON(http.StatusForbidden, res.ToMapStringInterface())
			return
		}
		c.JSON(http.StatusInternalServerError, res.ToMapStringInterface())
		return
	}

	c.JSON(http.StatusOK, res.ToMapStringInterface())
	return
}

func (u UserController) FetchUserByID(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			response.New(response.UserResponse{},
				false,
				statusCodes.Error,
				err.Error(),
				nil),
		)
		return
	}

	res := u.userUsecase.FetchUserByID(c, uint(userID))
	if res.IsFailed() {
		if res.ErrorIs(statusCodes.ModelNotFound) {
			c.JSON(http.StatusNotFound, res.ToMapStringInterface())
			return
		}
		c.JSON(http.StatusInternalServerError, res.ToMapStringInterface())
		return
	}

	c.JSON(http.StatusOK, res.ToMapStringInterface())
	return
}

func (u UserController) Me(c *gin.Context) {
	auth := c.Value(middleware.AuthenticatedRequest).(contracts.IAuthenticatedRequest)
	res := u.userUsecase.FetchUserByID(c, auth.GetUserID())
	if res.IsFailed() {
		if res.ErrorIs(statusCodes.ModelNotFound) {
			c.JSON(http.StatusNotFound, res.ToMapStringInterface())
			return
		}
		c.JSON(http.StatusInternalServerError, res.ToMapStringInterface())
		return
	}

	c.JSON(http.StatusOK, res.ToMapStringInterface())
	return
}

func (u UserController) CreateAccount(c *gin.Context) {

	var req requests.RegisterUserRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		if msg, ok := extractValidationErr(err); ok {
			c.JSON(
				http.StatusBadRequest,
				response.New(response.UserResponse{},
					false,
					statusCodes.Error,
					"request validation failed.",
					gin.H{"errors": msg},
				).ToMapStringInterface(),
			)
			return
		}
		c.JSON(
			http.StatusInternalServerError,
			response.New(response.UserResponse{},
				false,
				statusCodes.Error,
				err.Error(),
				nil),
		)
		return
	}

	res := u.userUsecase.CreateAccount(c, req)
	if res.IsFailed() {
		if res.ErrorIs(statusCodes.ErrForbidden) {
			c.JSON(http.StatusForbidden, res.ToMapStringInterface())
			return
		}
		if res.ErrorIs(statusCodes.ErrDuplicatedModel) {
			c.JSON(http.StatusConflict, res.ToMapStringInterface())
			return
		}
		c.JSON(http.StatusInternalServerError, res.ToMapStringInterface())
		return
	}

	c.JSON(http.StatusOK, res.ToMapStringInterface())
	return
}

func (u UserController) UpdateAccount(c *gin.Context) {
	var req requests.UpdateUserRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		if msg, ok := extractValidationErr(err); ok {
			c.JSON(
				http.StatusBadRequest,
				response.New(response.UserResponse{},
					false,
					statusCodes.Error,
					"failed process your request.",
					gin.H{"errors": msg},
				).ToMapStringInterface(),
			)
			return
		}
		c.JSON(
			http.StatusInternalServerError,
			response.New(response.UserResponse{},
				false,
				statusCodes.Error,
				err.Error(),
				nil).ToMapStringInterface(),
		)
		return
	}

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			response.New(response.UserResponse{},
				false,
				statusCodes.Error,
				err.Error(),
				nil),
		)
		return
	}

	res := u.userUsecase.UpdateAccount(c, uint(userID), req)
	if res.IsFailed() {
		if res.ErrorIs(statusCodes.ErrForbidden) {
			c.JSON(http.StatusForbidden, res.ToMapStringInterface())
			return
		}
		if res.ErrorIs(statusCodes.ModelNotFound) {
			res.SetMessage("record not found for user id " + c.Param("id"))
			c.JSON(http.StatusNotFound, res.ToMapStringInterface())
			return
		}
		c.JSON(http.StatusInternalServerError, res.ToMapStringInterface())
		return
	}

	c.JSON(http.StatusOK, res.ToMapStringInterface())
	return
}

func (u UserController) ChangePassword(c *gin.Context) {

	var req requests.ChangeUserPasswordRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		if msg, ok := extractValidationErr(err); ok {
			c.JSON(
				http.StatusBadRequest,
				response.New(response.UserResponse{},
					false,
					statusCodes.Error,
					"request validation failed.",
					gin.H{"errors": msg},
				).ToMapStringInterface(),
			)
			return
		}
		c.JSON(
			http.StatusInternalServerError,
			response.New(response.UserResponse{},
				false,
				statusCodes.Error,
				err.Error(),
				nil),
		)
		return
	}

	res := u.userUsecase.ChangePassword(c, req)
	if res.IsFailed() {
		if res.ErrorIs(statusCodes.ModelNotFound) {
			c.JSON(http.StatusNotFound, res.ToMapStringInterface())
			return
		}
		c.JSON(http.StatusInternalServerError, res.ToMapStringInterface())
		return
	}

	c.JSON(http.StatusOK, res.ToMapStringInterface())
	return
}

func (u UserController) ResetPassword(c *gin.Context) {

	var req requests.ResetUserPasswordRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		if msg, ok := extractValidationErr(err); ok {
			c.JSON(
				http.StatusBadRequest,
				response.New(response.UserResponse{},
					false,
					statusCodes.Error,
					"request validation failed.",
					gin.H{"errors": msg},
				).ToMapStringInterface(),
			)
			return
		}
		c.JSON(
			http.StatusInternalServerError,
			response.New(response.UserResponse{},
				false,
				statusCodes.Error,
				err.Error(),
				nil),
		)
		return
	}

	res := u.userUsecase.ResetPassword(c, req)
	if res.IsFailed() {
		if res.ErrorIs(statusCodes.ErrResetPasswordNotMatch) {
			c.JSON(http.StatusBadRequest, res.ToMapStringInterface())
			return
		}
		if res.ErrorIs(statusCodes.ErrExtractJWTToken) {
			c.JSON(http.StatusUnprocessableEntity, res.ToMapStringInterface())
			return
		}
		if res.ErrorIs(statusCodes.ModelNotFound) {
			c.JSON(http.StatusNotFound, res.ToMapStringInterface())
			return
		}
		c.JSON(http.StatusInternalServerError, res.ToMapStringInterface())
		return
	}

	c.JSON(http.StatusOK, res.ToMapStringInterface())
	return
}
