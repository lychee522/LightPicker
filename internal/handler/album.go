package handler

// @author 肖肖雨歇

import (
	"picgo-lite/internal/config"
	"picgo-lite/internal/model"

	"github.com/gin-gonic/gin"
)

func CreateAlbum(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		return
	}
	album := model.Album{Name: req.Name}
	config.DB.Create(&album)
	c.JSON(200, gin.H{"success": true, "message": "文件夹创建成功！"})
}

func GetAlbumList(c *gin.Context) {
	var albums []model.Album
	config.DB.Order("created_at desc").Find(&albums)
	c.JSON(200, gin.H{"success": true, "data": albums})
}

func RenameAlbum(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		return
	}
	config.DB.Model(&model.Album{}).Where("id = ?", id).Update("name", req.Name)
	c.JSON(200, gin.H{"success": true, "message": "重命名成功！"})
}

func DeleteAlbum(c *gin.Context) {
	id := c.Param("id")
	// 极致温柔的删除：把里面的图片移到根目录(0)，绝不误删实体图！
	config.DB.Model(&model.Image{}).Where("album_id = ?", id).Update("album_id", 0)
	config.DB.Delete(&model.Album{}, id)
	c.JSON(200, gin.H{"success": true, "message": "文件夹已粉碎，内部图片已转移至根目录保护！"})
}
