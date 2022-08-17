package middleware

import (
	"errors"

	"go-api/core/response"
	"go-api/service"

	"github.com/gin-gonic/gin"
)

func RbacCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := c.MustGet("claims").(*service.CustomClaims)
		role_id := claims.RoleID
		path := c.FullPath()
		method := c.Request.Method
		checked := service.NewRbacService().CheckPrivilege(role_id, path, method)
		if !checked {
			response.ResponseUnauthorized(c, "AuthError", errors.New("NO PRIVILEGE"))
			return
		}
		c.Next()
	}
}
