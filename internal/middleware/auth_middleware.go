package middleware

import (
	"strings"

	"github.com/AlexSamarskii/delivery_service/internal/entity"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct{}

func NewAuthMiddleware() AuthMiddleware {
	return AuthMiddleware{}
}

func (am AuthMiddleware) Auth() gin.HandlerFunc {
	return gin.BasicAuth(gin.Accounts{
		"pragmatic": "reviews",
	})
}

func extractToken(authorizationStr string) (string, error) {
	if authorizationStr == "" {
		return "", entity.ErrMissingAuthorizationHeader
	}

	if !strings.HasPrefix(authorizationStr, "Bearer ") {
		return "", entity.ErrInvalidAuthFormat
	}

	token := strings.TrimPrefix(authorizationStr, "Bearer ")
	if token == "" {
		return "", entity.ErrEmptyToken
	}

	return token, nil
}
