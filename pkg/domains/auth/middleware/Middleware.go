package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"paniq-auth-user-service/pkg/contracts/apiResources"
	"paniq-auth-user-service/pkg/contracts/statusCodes"
	"paniq-auth-user-service/pkg/response"
	"strconv"
	"strings"
)

const (
	AuthenticatedRequest     string = "auth-req"
	IntrospectedAccessToken  string = "intros-acc-token"
	IntrospectedRefreshToken string = "intros-ref-token"
)

// HandleIntrospectTokenMiddleware GIN MIDDLEWARE
func HandleIntrospectTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		targetResource := c.Request.Header.Get("X-Target-Resource")
		if targetResource == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				response.New(response.AuthResponse{},
					false,
					statusCodes.Error,
					"invalid resource!", nil).ToMapStringInterface())
			return
		}

		var accessToken string

		if token := c.Query("access-token"); token != "" {
			accessToken = token
		} else {
			if bearerToken := c.Request.Header.Get("Authorization"); len(strings.Split(bearerToken, " ")) == 2 {
				accessToken = strings.Split(bearerToken, " ")[1]
			}
		}

		c.Set(IntrospectedAccessToken, NewIntrospectTokenMiddleware(
			accessToken,
			apiResources.RESOURCE(targetResource)))
		c.Next()
	}
}

// HandleUnauthenticatedRequestMiddleware Setiap request hasil introspect token akan divalidasi di middleware ini
func HandleUnauthenticatedRequestMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//get user id
		userID, err := strconv.Atoi(c.Request.Header.Get("X-User-ID"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":      false,
				"status_code": statusCodes.Error,
				"message":     "user id is invalid",
				"data":        nil,
			})
			c.Abort()
			return
		}
		m := NewAuthenticatedRequestMiddleware(
			uint(userID),
			c.Request.Header.Get("X-User-Name"),
			c.Request.Header.Get("X-User-Role"),
			c.Request.Header.Get("X-User-Permission"),
		)
		if status, err := m.Validate(); !status {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":      false,
				"status_code": statusCodes.Error,
				"message":     err.Error(),
				"data":        nil,
			})
			c.Abort()
			return
		}

		c.Set(AuthenticatedRequest, m)
		c.Next()
	}
}

func HandleRefreshTokenRequestMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var refreshToken string
		if token := c.Query("refresh-token"); token != "" {
			refreshToken = token
		} else {
			if bearerToken := c.Request.Header.Get("Authorization"); len(strings.Split(bearerToken, " ")) == 2 {
				refreshToken = strings.Split(bearerToken, " ")[1]
			}
		}

		m := NewRefreshTokenMiddleware(refreshToken)
		if status, err := m.Validate(); !status {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":      false,
				"status_code": statusCodes.Error,
				"message":     err.Error(),
				"data":        nil,
			})
			c.Abort()
		}

		c.Set(IntrospectedRefreshToken, m)
		c.Next()
	}
}
