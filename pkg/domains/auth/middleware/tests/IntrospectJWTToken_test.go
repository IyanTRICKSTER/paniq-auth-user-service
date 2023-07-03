package middleware

import (
	"github.com/stretchr/testify/assert"
	"paniq-auth-user-service/pkg/contracts/apiResources"
	"paniq-auth-user-service/pkg/domains/auth/middleware"
	"testing"
)

func TestIntrospectJWTToken(t *testing.T) {
	i := middleware.NewIntrospectTokenMiddleware("token", apiResources.USER)
	assert.Equal(t, "token", i.GetAccessToken())
	assert.Equal(t, apiResources.USER, i.GetTargetResourceName())
}
