package handler

// @author 肖肖雨歇
// @description 核心上传业务：支持哈希秒传、记录宽高比例与归属文件夹

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"picgo-lite/internal/config"
	"picgo-lite/internal/model"
	"picgo-lite/internal/utils"

	"github.com/gin-gonic/gin"
)

func UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件没传过来啊兄弟！"})
		return
	}

	// 划重点：接收前端白嫖算力算出来的宽高和文件夹ID
	albumID, _ := strconv.Atoi(c.PostForm("album_id"))
	width, _ := strconv.Atoi(c.PostForm("width"))
	height, _ := strconv.Atoi(c.PostForm("height"))

	hashStr, err := utils.CalcFileHash(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "计算文件哈希翻车了"})
		return
	}

	var existingImage model.Image
	if err := config.DB.Where("hash = ?", hashStr).First(&existingImage).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "触发秒传！", "url": "/" + existingImage.StoragePath})
		return
	}

	yearMonth := time.Now().Format("200601")
	saveDir := filepath.Join("storage", "uploads", yearMonth)
	os.MkdirAll(saveDir, os.ModePerm)

	ext := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), hashStr[:8], ext)
	savePath := filepath.Join(saveDir, newFileName)
	relativePath := "uploads/" + yearMonth + "/" + newFileName

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件写入硬盘失败了"})
		return
	}

	newImage := model.Image{
		Filename:    file.Filename,
		StoragePath: relativePath,
		Hash:        hashStr,
		Size:        file.Size,
		Width:       width,         // 记录神圣的宽度
		Height:      height,        // 记录神圣的高度
		AlbumID:     uint(albumID), // 存入对应的文件夹 (0代表根目录)
	}
	config.DB.Create(&newImage)

	c.JSON(http.StatusOK, gin.H{"success": true, "url": "/" + relativePath})
}
