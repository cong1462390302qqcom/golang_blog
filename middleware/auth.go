package middleware

import (
	"github.com/gin-gonic/gin"
	"golang_blog/utils"
	"strings"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}
		//验证token
		// 检查Bearer前缀
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			utils.Unauthorized(c, "Bearer token is required")
			c.Abort()
			return
		}
		//解析token
		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			utils.Unauthorized(c, "Invalid token")
			c.Abort()
			return
		}
		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.UserName)
		c.Next()
	}
}
