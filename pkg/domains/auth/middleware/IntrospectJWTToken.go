package middleware

import (
	"paniq-auth-user-service/pkg/contracts"
	"paniq-auth-user-service/pkg/contracts/apiResources"
)

type IntrospectTokenMiddleware struct {
	targetResource apiResources.RESOURCE
	accessToken    string
}

func (i IntrospectTokenMiddleware) GetTargetResourceName() apiResources.RESOURCE {
	return i.targetResource
}

func (i IntrospectTokenMiddleware) GetAccessToken() string {
	return i.accessToken
}

func NewIntrospectTokenMiddleware(accessToken string, resource apiResources.RESOURCE) contracts.IMiddlewareIntrospectToken {
	return IntrospectTokenMiddleware{accessToken: accessToken, targetResource: resource}
}
