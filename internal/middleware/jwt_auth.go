package middleware

// @author 肖肖雨歇
// @description JWT 鉴权拦截器：同时支持 Header 和 URL 参数 (专治文件下载)

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := ""

		// 1. 优先从 Header 取 (常规 API 请求)
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = authHeader[7:]
		} else {
			// 2. 如果 Header 没有，尝试从 URL 参数取 (专为备份下载设计)
			tokenString = c.Query("token")
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "站住！没带通行证(Token)不许传图或下载！"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("picgo-lite-xiaoxiao-secret-key-666"), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "通行证是假的或者过期了，重新登录试试！"})
			c.Abort()
			return
		}

		c.Next()
	}
}
