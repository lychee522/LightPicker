package middleware

// @author 肖肖雨歇
// @description 全局拦截器：检测系统是否已初始化管理员账号

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"picgo-lite/internal/config" // 注意替换为你的实际模块名
	"picgo-lite/internal/model"
)

// CheckInitMiddleware 检查系统是否已初始化的 Gin 中间件
func CheckInitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var count int64
		// 查询 User 表里是否有账号
		err := config.DB.Model(&model.User{}).Count(&count).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库查询异常"})
			c.Abort()
			return
		}

		// 划重点：如果请求的是 /api/init (初始化接口本身)，必须放行，不然就死循环了！
		if c.Request.URL.Path == "/api/init" {
			c.Next()
			return
		}

		// 如果数据库里没有账号，拦截请求并提示前端跳转到初始化页面
		if count == 0 {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "系统尚未初始化，请先创建管理员账号！",
				"action":  "redirect_to_init",
			})
			c.Abort()
			return
		}

		// 已经初始化过，继续执行后续路由逻辑
		c.Next()
	}
}