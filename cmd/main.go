package main

// @author 肖肖雨歇
// @description 极简图床完全体：彻底消灭 301 重定向死循环，完美合体 SPA！

import (
	"io/fs"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"picgo-lite/internal/config"
	"picgo-lite/internal/handler"
	"picgo-lite/internal/middleware"
	"picgo-lite/web" // 🌟 引入前端静态文件包

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 初始化数据库与配置
	config.InitDB("./storage/data.db")
	handler.LoadSettings()

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
	}

	// 4. 提取前端打包文件
	dist, err := fs.Sub(web.DistFS, "dist")
	if err != nil {
		log.Fatalf("前端打包文件提取失败: %v", err)
	}

	// 🌟 5. 终极 SPA 挂载方案：暴力读写，彻底绕开重定向大坑！
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// 如果是错误的接口或图片地址，直接报 404
		if strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/uploads/") {
			c.JSON(404, gin.H{"error": "接口或文件不存在"})
			return
		}

		// 处理路径
		filePath := strings.TrimPrefix(path, "/")
		if filePath == "" {
			filePath = "index.html"
		}

		// 检查该文件是否存在于前端资源中
		file, err := dist.Open(filePath)
		if err != nil {
			// 文件找不到（比如 Vue 的内部路由 /settings），统统交给 index.html 渲染
			filePath = "index.html"
		} else {
			file.Close() // 找到了就关掉，仅用于检查是否存在
		}

		// 🌟 核心：直接读取字节流并返回，绝不重定向！
		data, err := fs.ReadFile(dist, filePath)
		if err != nil {
			c.String(500, "前端资源已损坏或未打包")
			return
		}

		// 智能识别文件类型 (Content-Type)
		contentType := mime.TypeByExtension(filepath.Ext(filePath))
		if contentType == "" {
			contentType = http.DetectContentType(data)
		}

		// 啪！直接把前端页面糊到浏览器上！
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
