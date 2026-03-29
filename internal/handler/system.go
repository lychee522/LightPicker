package handler

// @author 肖肖雨歇
// @description 系统底层操作：数据库的一键备份与恢复

import (
	"net/http"
	"picgo-lite/internal/config"

	"github.com/gin-gonic/gin"
)

func BackupDB(c *gin.Context) {
	// 直接把底层的 SQLite 数据库文件发给前端下载
	c.FileAttachment("./storage/data.db", "picgo_data_backup.db")
}

func RestoreDB(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "没收到数据库文件"})
		return
	}
	// 强制覆盖当前的数据库文件
	if err := c.SaveUploadedFile(file, "./storage/data.db"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "覆盖数据库失败"})
		return
	}
	// 重新连接数据库，让新数据生效
	config.InitDB("./storage/data.db")
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "数据库恢复成功，已重新连接！"})
}
