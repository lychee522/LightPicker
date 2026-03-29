package handler

// @author 肖肖雨歇
// @description 认证与初始化控制器，负责账号创建与登录逻辑

import (
	"net/http"

	"picgo-lite/internal/config" // 注意替换为你的实际模块名
	"picgo-lite/internal/model"
	"picgo-lite/internal/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// InitAdmin 首次初始化管理员账号
func InitAdmin(c *gin.Context) {
	var count int64
	config.DB.Model(&model.User{}).Count(&count)
	if count > 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "好兄弟，系统已经初始化过了，别想梅开二度！"})
		return
	}

	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数不对啊，账号或密码没填全？"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	user := model.User{
		Username: req.Username,
		Password: string(hashedPassword),
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建账号失败了"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "太棒了！管理员账号创建成功，准备起飞！"})
}

// Login 管理员登录接口
func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "账号密码不能为空！"})
		return
	}

	var user model.User
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 校验密码哈希
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 密码正确，签发 JWT 通行证
	token, err := utils.GenerateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token签发失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功，好兄弟欢迎回来！",
		"token":   token,
	})
}
