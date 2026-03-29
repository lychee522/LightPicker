package handler

// @author 肖肖雨歇
// @description 系统底层配置：保存防盗链白名单，并动态刷新内存

import (
	"net/http"
	"strings"

	"picgo-lite/internal/config"
	"picgo-lite/internal/middleware"

	"github.com/gin-gonic/gin"
)

// Setting 极简配置表，直接在这里定义，轻量到底
type Setting struct {
	ID    uint   `gorm:"primaryKey"`
	Key   string `gorm:"uniqueIndex"`
	Value string
}

func GetWhitelist(c *gin.Context) {
	var s Setting
	config.DB.Where("key = ?", "whitelist").First(&s)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": s.Value})
}

func SaveWhitelist(c *gin.Context) {
	var req struct {
		Value string `json:"value"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		return
	}

	var s Setting
	if err := config.DB.Where("key = ?", "whitelist").First(&s).Error; err != nil {
		config.DB.Create(&Setting{Key: "whitelist", Value: req.Value})
	} else {
		config.DB.Model(&s).Update("value", req.Value)
	}

	// 🌟 最牛逼的一步：保存到数据库后，瞬间刷新内存里的保安名单！
	if req.Value == "" {
		middleware.Whitelist = []string{}
	} else {
		middleware.Whitelist = strings.Split(req.Value, ",")
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "防盗链神盾已生效！"})
}

// LoadSettings 初始化时被 main.go 调用，自动建表并加载缓存
func LoadSettings() {
	config.DB.AutoMigrate(&Setting{})
	var s Setting
	if err := config.DB.Where("key = ?", "whitelist").First(&s).Error; err == nil && s.Value != "" {
		middleware.Whitelist = strings.Split(s.Value, ",")
	}
}
