package handler

// @author tg账号的肖肖雨歇
// @description 认证与初始化控制器，负责账号创建与登录逻辑 (新增内存级 IP 防爆破护盾)

import (
	"net/http"
	"sync"
	"time"

	"picgo-lite/internal/config"
	"picgo-lite/internal/model"
	"picgo-lite/internal/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// --- 🌟 核心安全：内存级 IP 防爆破拦截器 ---
type LoginAttempt struct {
	Count       int
	LastAttempt time.Time
}

var (
	failedAttempts = make(map[string]*LoginAttempt)
	attemptMutex   sync.RWMutex
	MaxAttempts    = 10               // 最大连续错误次数
	LockoutDur     = 30 * time.Minute // 封禁时长
)

// checkAndRecordIP 检查并记录 IP 错误次数
func checkAndRecordIP(ip string, isSuccess bool) (bool, time.Duration) {
	attemptMutex.Lock()
	defer attemptMutex.Unlock()

	attempt, exists := failedAttempts[ip]

	// 如果登录成功，清除该 IP 的错误记录
	if isSuccess {
		if exists {
			delete(failedAttempts, ip)
		}
		return true, 0
	}

	// 登录失败逻辑
	now := time.Now()
	if !exists {
		failedAttempts[ip] = &LoginAttempt{Count: 1, LastAttempt: now}
		return true, 0
	}

	// 如果仍在封禁期内
	if attempt.Count >= MaxAttempts && now.Sub(attempt.LastAttempt) < LockoutDur {
		remain := LockoutDur - now.Sub(attempt.LastAttempt)
		return false, remain
	}

	// 如果过了封禁期，重新开始计数
	if now.Sub(attempt.LastAttempt) >= LockoutDur {
		attempt.Count = 1
		attempt.LastAttempt = now
		return true, 0
	}

	// 累加错误次数
	attempt.Count++
	attempt.LastAttempt = now

	// 刚刚达到阈值，触发封禁
	if attempt.Count >= MaxAttempts {
		return false, LockoutDur
	}

	return true, 0
}

// ----------------------------------------

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
	// 获取请求 IP
	clientIP := c.ClientIP()

	// 🌟 防爆破护盾：前置拦截
	allowed, remain := checkAndRecordIP(clientIP, false) // 先预判一下，不增加错误次数
	if !allowed {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": "密码错误次数过多，IP已被封禁，请在 " + remain.Round(time.Second).String() + " 后重试，或重启程序解除封禁。",
		})
		return
	}

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
		checkAndRecordIP(clientIP, false) // 记录错误
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 校验密码哈希
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		checkAndRecordIP(clientIP, false) // 记录错误
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 🌟 登录成功，清除该 IP 的错误计数器
	checkAndRecordIP(clientIP, true)

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
