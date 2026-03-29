package model

// @author 肖肖雨歇
// @description 极简图床核心数据库模型，默认采用 SQLite 驱动

import (
	"time"

	"gorm.io/gorm"
)

// User 用户表：记录管理员账号信息 (用于鉴权)
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Username  string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"` // 密码Hash，JSON中忽略
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Image 图片表：记录图片元数据，支持哈希去重和按比例筛选
type Image struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	Filename     string         `gorm:"type:varchar(255);not null" json:"filename"` // 原始文件名
	StoragePath  string         `gorm:"type:varchar(255);not null" json:"url"`      // 物理存储相对路径
	Hash         string         `gorm:"type:varchar(64);index;not null" json:"hash"`// 文件内容哈希(MD5/SHA1)，用于秒传和查重
	MimeType     string         `gorm:"type:varchar(50)" json:"mime_type"`          // 文件类型，如 image/webp
	Size         int64          `gorm:"not null" json:"size"`                       // 文件大小(字节)
	Width        int            `gorm:"not null" json:"width"`                      // 图片宽度 (用于随机图比例筛选)
	Height       int            `gorm:"not null" json:"height"`                     // 图片高度 (用于随机图比例筛选)
	AlbumID      uint           `gorm:"index" json:"album_id"`                      // 归属相册ID
	IsWebp       bool           `gorm:"default:false" json:"is_webp"`               // 是否已转换为WebP
	CreatedAt    time.Time      `gorm:"index" json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// Album 相册/目录表：用于图片的分类组织
type Album struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`
	Description string         `gorm:"type:varchar(255)" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// Setting 系统配置表：存储域名、水印、白名单等全局配置
type Setting struct {
	ID          uint   `gorm:"primarykey" json:"id"`
	Key         string `gorm:"type:varchar(50);uniqueIndex;not null" json:"key"` // 配置键名，如 "domain", "watermark_text"
	Value       string `gorm:"type:text" json:"value"`                           // 配置内容
	Description string `gorm:"type:varchar(255)" json:"description"`             // 配置说明
}

// AutoMigrate 自动迁移数据库结构
func AutoMigrate(db *gorm.DB) error {
	// GORM 会自动检查结构体并创建/更新对应的 SQLite 表
	return db.AutoMigrate(
		&User{},
		&Image{},
		&Album{},
		&Setting{},
	)
}