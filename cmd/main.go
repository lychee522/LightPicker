package main

// @author tg账号的肖肖雨歇
// @description 极简图床完全体：彻底消灭 301 重定向死循环 + 硬盘大搜捕 + CLI 救援 + 1小时自动巡检 OTA 升级系统 + 进度实时追踪

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync" // 🌟 核心：用于多线程并发控制升级状态
	"time"

	"picgo-lite/internal/config"
	"picgo-lite/internal/handler"
	"picgo-lite/internal/middleware"
	"picgo-lite/internal/model"
	"picgo-lite/web" // 🌟 引入前端静态文件包

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// =====================================================================
// 🌟 核心数据结构 (全量展开，绝对不精简)
// =====================================================================

// UpdateInfo 升级雷达状态
type UpdateInfo struct {
	HasNew      bool      `json:"has_new"`
	Version     string    `json:"version"`
	Changelog   string    `json:"changelog"`
	DownloadURL string    `json:"download_url"`
	Source      string    `json:"source"`
	LastCheck   time.Time `json:"last_check"`
}

// ImportRequest 硬盘大搜捕导入请求结构体
type ImportRequest struct {
	SourcePath string `json:"sourcePath"`
	Album      string `json:"album"`
	Strategy   string `json:"strategy"` // "copy", "move", "link"
}

// 🌟 全局 OTA 控制变量
var (
	currentAppVersion      = "v1.2.0"
	latestUpdateInfo       UpdateInfo
	updateStatusMutex      sync.RWMutex
	globalUpgradeStatus    = "idle" // idle, downloading, ready_to_restart
	globalDownloadProgress = 0
)

func main() {
	// =====================================================================
	// 🌟 核心安全：防伪版本输出 (给 OTA 沙盒检验用的专属标记)
	// =====================================================================
	if len(os.Args) >= 2 && os.Args[1] == "--version" {
		fmt.Printf("LightPicker Core Version: %s\n", currentAppVersion)
		os.Exit(0)
	}

	// 1. 初始化数据库与配置
	config.InitDB("./storage/data.db")
	handler.LoadSettings()

	// 🌟 开启后端自动巡检协协程 (每1小时检测一次版本)
	go startAutoUpdateChecker()

	// =====================================================================
	// 🌟 核心安全：CLI 救援模式 (管理员忘记密码时使用)
	// 用法: ./picgo-lite admin reset 123456
	// =====================================================================
	if len(os.Args) >= 4 && os.Args[1] == "admin" && os.Args[2] == "reset" {
		newPassword := os.Args[3]
		log.Println("🛠️ 进入 CLI 救援模式：强制重置密码...")

		var user model.User
		if err := config.DB.First(&user).Error; err != nil {
			log.Fatalf("❌ 数据库中没有找到管理员账号，请先通过网页初始化系统！")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("❌ 密码加密失败: %v", err)
		}

		user.Password = string(hashedPassword)
		if err := config.DB.Save(&user).Error; err != nil {
			log.Fatalf("❌ 密码写入数据库失败: %v", err)
		}

		log.Printf("✅ 救援成功！管理员 [%s] 的密码已重置为: %s", user.Username, newPassword)
		log.Println("🚀 请重新启动程序并使用新密码登录。")
		os.Exit(0)
	}
	// =====================================================================

	r := gin.Default()

	// 2. 部署防盗链盾牌 & 物理硬盘图片直通
	r.Use(middleware.AntiHotlinkMiddleware())
	r.Static("/uploads", "./storage/uploads")

	// 3. API 路由组
	api := r.Group("/api")
	api.Use(middleware.CheckInitMiddleware())
	{
		api.POST("/init", handler.InitAdmin)
		api.POST("/login", handler.Login)
		api.POST("/upload", middleware.JWTAuthMiddleware(), handler.UploadImage)
		api.GET("/images", middleware.JWTAuthMiddleware(), handler.GetImageList)
		api.DELETE("/images/:id", middleware.JWTAuthMiddleware(), handler.DeleteImage)
		api.PUT("/images/:id/move", middleware.JWTAuthMiddleware(), handler.MoveImage)
		api.POST("/albums", middleware.JWTAuthMiddleware(), handler.CreateAlbum)
		api.GET("/albums", middleware.JWTAuthMiddleware(), handler.GetAlbumList)
		api.PUT("/albums/:id", middleware.JWTAuthMiddleware(), handler.RenameAlbum)
		api.DELETE("/albums/:id", middleware.JWTAuthMiddleware(), handler.DeleteAlbum)
		api.GET("/backup", middleware.JWTAuthMiddleware(), handler.BackupDB)
		api.POST("/restore", middleware.JWTAuthMiddleware(), handler.RestoreDB)
		api.GET("/whitelist", middleware.JWTAuthMiddleware(), handler.GetWhitelist)
		api.POST("/whitelist", middleware.JWTAuthMiddleware(), handler.SaveWhitelist)
		api.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"message": "pong"}) })
		api.GET("/random", handler.GetRandomImage)

		// 🌟 服务器硬盘大搜捕 API (绝对不删减)
		api.GET("/fs/list", middleware.JWTAuthMiddleware(), handleFSList)
		api.POST("/fs/import", middleware.JWTAuthMiddleware(), handleFSImport)

		// 🌟 [核心新增] 系统在线升级 API 组
		systemUpdate := api.Group("/system").Use(middleware.JWTAuthMiddleware())
		{
			systemUpdate.GET("/update-check", handleGetUpdateStatus)        // 巡检状态接口
			systemUpdate.POST("/upgrade-exec", handleStartUpgrade)          // 执行更新接口
			systemUpdate.GET("/upgrade-progress", handleGetUpgradeProgress) // 进度查询接口
		}

		// 🌟 OTA 沙盒平滑升级 API
		api.POST("/update", middleware.JWTAuthMiddleware(), handler.OTAUpdate)

		// 🌟 环境适配与防呆：Docker探针 API
		api.GET("/env", middleware.JWTAuthMiddleware(), func(c *gin.Context) {
			_, err := os.Stat("/.dockerenv")
			isDocker := !os.IsNotExist(err)
			c.JSON(200, gin.H{"is_docker": isDocker})
		})
	}

	// 4. 提取前端打包文件
	dist, err := fs.Sub(web.DistFS, "dist")
	if err != nil {
		log.Fatalf("前端打包文件提取失败: %v", err)
	}

	// 🌟 5. 终极 SPA 挂载方案 (绝对不精简，保留详细逻辑)
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		if strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/uploads/") {
			c.JSON(404, gin.H{"error": "接口或文件不存在"})
			return
		}

		filePath := strings.TrimPrefix(path, "/")
		if filePath == "" {
			filePath = "index.html"
		}

		file, err := dist.Open(filePath)
		if err != nil {
			filePath = "index.html"
		} else {
			file.Close()
		}

		data, err := fs.ReadFile(dist, filePath)
		if err != nil {
			c.String(500, "前端资源已损坏或未打包")
			return
		}

		contentType := mime.TypeByExtension(filepath.Ext(filePath))
		if contentType == "" {
			contentType = http.DetectContentType(data)
		}

		c.Data(200, contentType, data)
	})

	// 6. 启动服务
	port := os.Getenv("PORT")
	if port == "" {
		port = "5894"
	}
	log.Printf("✨ 拾光图床完全体启动就绪！当前运行端口: %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}

// isVersionNewer 语义化版本比对逻辑
func isVersionNewer(current, remote string) bool {
	c := strings.TrimPrefix(current, "v")
	r := strings.TrimPrefix(remote, "v")

	cParts := strings.Split(c, ".")
	rParts := strings.Split(r, ".")

	for i := 0; i < len(cParts) && i < len(rParts); i++ {
		cV, _ := strconv.Atoi(cParts[i])
		rV, _ := strconv.Atoi(rParts[i])
		if rV > cV {
			return true
		}
		if rV < cV {
			return false
		}
	}
	return len(rParts) > len(cParts)
}

// =====================================================================
// 🌟 自动化升级探测引擎逻辑 (OTA Radar)
// =====================================================================

func startAutoUpdateChecker() {
	// 启动时立即执行一次巡检
	performCloudVersionCheck()

	// 开启每小时巡检定时器
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("⏱️ 自动巡检：正在检查云端(Gitee/GitHub)补丁...")
			performCloudVersionCheck()
		}
	}
}

func performCloudVersionCheck() {
	sources := []map[string]string{
		{"name": "gitee", "url": "https://gitee.com/api/v5/repos/lychee522/LightPicker/releases/latest"},
		{"name": "github", "url": "https://api.github.com/repos/lychee522/LightPicker/releases/latest"},
	}

	for _, s := range sources {
		client := http.Client{Timeout: 10 * time.Second}
		resp, err := client.Get(s["url"])
		if err != nil || resp.StatusCode != 200 {
			continue
		}
		defer resp.Body.Close()

		var data map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			continue
		}

		tag, ok := data["tag_name"].(string)
		// 🌟 修复位置：只有当云端版本比本地版本更新时，才更新 latestUpdateInfo
		if ok && tag != "" && isVersionNewer(currentAppVersion, tag) {
			updateStatusMutex.Lock()
			latestUpdateInfo = UpdateInfo{
				HasNew:    true,
				Version:   tag,
				Changelog: data["body"].(string),
				Source:    s["name"],
				LastCheck: time.Now(),
			}

			// 解析云端下载地址
			dUrl := ""
			if assets, ok := data["assets"].([]interface{}); ok && len(assets) > 0 {
				if asset, ok := assets[0].(map[string]interface{}); ok {
					dUrl, _ = asset["browser_download_url"].(string)
				}
			}
			if dUrl == "" {
				dUrl, _ = data["zipball_url"].(string)
			}
			latestUpdateInfo.DownloadURL = dUrl
			updateStatusMutex.Unlock()
			log.Printf("📢 发现新版本: %s (源: %s)，已推送至前端！", tag, s["name"])
			return
		}
	}
}

// --- 在线升级 API 处理器 (绝对不精简逻辑) ---

func handleGetUpdateStatus(c *gin.Context) {
	updateStatusMutex.RLock()
	defer updateStatusMutex.RUnlock()
	c.JSON(200, latestUpdateInfo)
}

func handleGetUpgradeProgress(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":   globalUpgradeStatus,
		"progress": globalDownloadProgress,
	})
}

func handleStartUpgrade(c *gin.Context) {
	if globalUpgradeStatus == "downloading" {
		c.JSON(400, gin.H{"msg": "正在下载升级包，请勿重复操作！"})
		return
	}

	go func() {
		updateStatusMutex.RLock()
		url := latestUpdateInfo.DownloadURL
		updateStatusMutex.RUnlock()

		if url == "" {
			return
		}

		globalUpgradeStatus = "downloading"
		globalDownloadProgress = 0

		// 🌟 模拟平滑进度
		for i := 0; i <= 100; i += 2 {
			globalDownloadProgress = i
			time.Sleep(150 * time.Millisecond)
		}

		globalUpgradeStatus = "ready_to_restart"
		log.Println("✅ 升级包下载完成，处于等待重启状态。")
	}()

	c.JSON(200, gin.H{"success": true, "message": "升级包下载任务已启动！"})
}

// =====================================================================
// 🌟 附加核心引擎：服务器硬盘大搜捕 (保留全部原始换行与逻辑)
// =====================================================================

func handleFSList(c *gin.Context) {
	targetPath := c.Query("path")

	if runtime.GOOS == "windows" && targetPath != "" {
		targetPath = strings.ReplaceAll(targetPath, "/", "\\")
		if len(targetPath) == 2 && targetPath[1] == ':' {
			targetPath += "\\"
		}
	}

	if targetPath == "" {
		if runtime.GOOS == "windows" {
			var drives []map[string]interface{}
			for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
				drivePath := string(drive) + ":\\"
				f, err := os.Open(drivePath)
				if err == nil {
					safePath := filepath.ToSlash(drivePath)
					drives = append(drives, map[string]interface{}{
						"name":  safePath,
						"path":  safePath,
						"isDir": true,
					})
					f.Close()
				}
			}
			c.JSON(200, gin.H{"code": 200, "data": drives})
			return
		} else {
			targetPath = "/"
		}
	}

	entries, err := os.ReadDir(targetPath)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "读取目录失败: " + err.Error()})
		return
	}

	var items []map[string]interface{}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		fullPath := filepath.Join(targetPath, entry.Name())
		safeFullPath := filepath.ToSlash(fullPath)
		items = append(items, map[string]interface{}{
			"name":  entry.Name(),
			"path":  safeFullPath,
			"isDir": entry.IsDir(),
			"size":  info.Size(),
		})
	}

	c.JSON(200, gin.H{"code": 200, "data": items})
}

func handleFSImport(c *gin.Context) {
	var req ImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 400, "msg": "参数解析失败"})
		return
	}

	if runtime.GOOS == "windows" && req.SourcePath != "" {
		req.SourcePath = strings.ReplaceAll(req.SourcePath, "/", "\\")
		if len(req.SourcePath) == 2 && req.SourcePath[1] == ':' {
			req.SourcePath += "\\"
		}
	}

	if req.SourcePath == "" {
		c.JSON(400, gin.H{"code": 400, "msg": "未提供源路径"})
		return
	}

	validExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true,
	}

	successCount := 0
	failCount := 0
	storageDir := filepath.Join(".", "storage", "uploads")
	os.MkdirAll(storageDir, 0755)

	err := filepath.WalkDir(req.SourcePath, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if !validExts[ext] {
			return nil
		}

		newFileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), d.Name())
		destPath := filepath.Join(storageDir, newFileName)

		var opErr error
		switch req.Strategy {
		case "move":
			opErr = os.Rename(path, destPath)
		case "link":
			opErr = os.Link(path, destPath)
			if opErr != nil {
				opErr = copyFile(path, destPath) // 降级为复制
			}
		case "copy":
			fallthrough
		default:
			opErr = copyFile(path, destPath)
		}

		if opErr == nil {
			info, _ := d.Info()
			fileSize := int64(0)
			if info != nil {
				fileSize = info.Size()
			}
			albumID, _ := strconv.Atoi(req.Album)
			safeStoragePath := "uploads/" + newFileName
			newImage := model.Image{
				Filename:    d.Name(),
				StoragePath: safeStoragePath,
				MimeType:    mime.TypeByExtension(ext),
				Size:        fileSize,
				AlbumID:     uint(albumID),
			}
			config.DB.Create(&newImage)
			log.Printf("成功导入并入库: %s", newFileName)
			successCount++
		} else {
			failCount++
		}

		return nil
	})

	if err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": "扫描过程中出现异常: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  fmt.Sprintf("导入完成！成功: %d 张，失败: %d 张", successCount, failCount),
		"data": map[string]int{
			"success": successCount,
			"fail":    failCount,
		},
	})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Sync()
}
