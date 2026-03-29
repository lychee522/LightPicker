package middleware

// @author 肖肖雨歇
// @description 终极神盾防盗链：大道至简 Contains 匹配 + Origin 备用检查

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var Whitelist []string

func AntiHotlinkMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 只管 /uploads/ 目录下的实体图片
		if !strings.HasPrefix(c.Request.URL.Path, "/uploads/") {
			c.Next()
			return
		}

		// 2. 如果没配白名单，天下大同，全部放行！
		if len(Whitelist) == 0 {
			c.Next()
			return
		}

		// 3. 尝试查户口：先摸 Referer 口袋，没有再摸 Origin 口袋
		referer := c.Request.Header.Get("Referer")
		if referer == "" {
			referer = c.Request.Header.Get("Origin")
		}

		if referer == "" {
			c.Next()
			return
		}

		// 4. 绝对信任本地开发环境

		// 4. 绝对信任本地开发环境
		if strings.Contains(referer, "localhost") || strings.Contains(referer, "127.0.0.1") {
			c.Next()
			return
		}

		// 5. 绝对信任图床服务器自己的真实域名
		reqHost := strings.Split(c.Request.Host, ":")[0]
		if strings.Contains(referer, reqHost) {
			c.Next()
			return
		}

		// 6. 大道至简：暴力 Contains 匹配！只要来源里包含你写的白名单字符串，全部放行！
		allowed := false
		for _, domain := range Whitelist {
			d := strings.TrimSpace(domain)
			if d == "" {
				continue
			}
			if strings.Contains(referer, d) {
				allowed = true
				break
			}
		}

		// 7. 踢出局
		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "盗链拦截：您的域名未在白名单中！"})
			return
		}

		c.Next()
	}
}
