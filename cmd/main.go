package main

// @author 肖肖雨歇
// @description 极简图床完全体：防盗链盾牌就位！

import (
	"log"
	"os"

	"picgo-lite/internal/config"
	"picgo-lite/internal/handler"
	"picgo-lite/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 初始化数据库
	config.InitDB("./storage/data.db")

	// 2. 加载系统配置，唤醒防盗链保安
	handler.LoadSettings()

	r := gin.Default()

	// 🌟 3. 部署防盗链盾牌！(必须放在静态文件服务之前)
	r.Use(middleware.AntiHotlinkMiddleware())
	// 物理硬盘图片直通路由
	r.Static("/uploads", "./storage/uploads")

	// 4. API 路由组
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

		// 🌟 新增：防盗链配置路由
		api.GET("/whitelist", middleware.JWTAuthMiddleware(), handler.GetWhitelist)
		api.POST("/whitelist", middleware.JWTAuthMiddleware(), handler.SaveWhitelist)

		api.GET("/ping", func(c *gin.Context) { c.JSON(200, gin.H{"message": "pong"}) })
		api.GET("/random", handler.GetRandomImage)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "5894"
	}
	log.Printf("✨ 拾光图床完全体启动就绪！当前运行端口: %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
