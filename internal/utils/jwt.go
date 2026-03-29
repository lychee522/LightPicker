package utils

// @author 肖肖雨歇
// @description JWT 令牌签发与验证工具类

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 极其机密的私钥，好兄弟可以自己随便改一串脸滚键盘的字符
var jwtSecret = []byte("picgo-lite-xiaoxiao-secret-key-666")

// GenerateToken 签发通行证，有效期 7 天
func GenerateToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(7 * 24 * time.Hour).Unix(), // 7天后过期
		"iat":      time.Now().Unix(),                         // 签发时间
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
