package handler

// @author 肖肖雨歇
// @description 图库大厅与随机图 API：新增文件夹过滤与横竖屏盲盒支持

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"picgo-lite/internal/config"
	"picgo-lite/internal/model"

	"github.com/gin-gonic/gin"
)

// GetImageList 获取图库列表 (支持按文件夹过滤)
func GetImageList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "24"))
	albumID := c.Query("album_id")
	offset := (page - 1) * pageSize

	var images []model.Image
	var total int64

	query := config.DB.Model(&model.Image{})
	// 如果传了具体的文件夹ID，就只查这个文件夹的图
	if albumID != "" && albumID != "0" {
		query = query.Where("album_id = ?", albumID)
	}

	query.Count(&total)
	query.Order("created_at desc").Offset(offset).Limit(pageSize).Find(&images)

	c.JSON(http.StatusOK, gin.H{
		"success": true, "total": total, "page": page, "size": pageSize, "data": images,
	})
}

func DeleteImage(c *gin.Context) {
	id := c.Param("id")
	var img model.Image
	if err := config.DB.First(&img, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "这图已经灰飞烟灭了"})
		return
	}
	var count int64
	config.DB.Model(&model.Image{}).Where("hash = ?", img.Hash).Count(&count)
	if count <= 1 {
		os.Remove(filepath.Join("storage", img.StoragePath))
	}
	config.DB.Delete(&img)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "删除成功，硬盘空间+1！"})
}

// GetRandomImage 对外公开的随机图盲盒 API (高阶：支持横竖屏与文件夹双重筛选)
func GetRandomImage(c *gin.Context) {
	ori := c.Query("ori")
	albumID := c.Query("album_id") // 🌟 新增：获取指定的文件夹ID
	var img model.Image

	query := config.DB.Model(&model.Image{})

	// 🌟 新增：如果指定了文件夹，就只在这个文件夹里随机抽
	if albumID != "" && albumID != "0" {
		query = query.Where("album_id = ?", albumID)
	}

	if ori == "landscape" {
		query = query.Where("width >= height AND width > 0")
	} else if ori == "portrait" {
		query = query.Where("height > width AND height > 0")
	}

	if err := query.Order("RANDOM()").First(&img).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "图库里没有符合条件的图啊兄弟！"})
		return
	}
	// 返回图片直链
	c.Redirect(http.StatusFound, "/"+img.StoragePath)
}

// MoveImage 移动图片到其他文件夹
func MoveImage(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		AlbumID string `json:"album_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		return
	}

	config.DB.Model(&model.Image{}).Where("id = ?", id).Update("album_id", req.AlbumID)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "图片已移动！"})
}
