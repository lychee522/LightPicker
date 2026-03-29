package config

// @author 肖肖雨歇
// @description 数据库连接初始化与全局管理，纯Go驱动拒绝CGO烦恼

import (
	"log"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite" // 纯Go实现的SQLite驱动，完美契合跨平台极简部署
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"picgo-lite/internal/model" // 注意：这里的 picgo-lite 请替换为你实际的 go.mod 模块名
)

// DB 全局数据库实例
var DB *gorm.DB

// InitDB 初始化SQLite数据库连接
func InitDB(dbPath string) {
	// 确保数据库文件所在的目录存在
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatalf("好兄弟，创建数据库目录失败了: %v", err)
	}

	var err error
	// 连接 SQLite，开启默认的日志输出方便咱们调试
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("数据库连接翻车了，快检查路径: %v", err)
	}

	// 自动迁移咱们上一步设计好的表结构
	err = model.AutoMigrate(DB)
	if err != nil {
		log.Fatalf("表结构迁移失败: %v", err)
	}

	log.Println("SQLite 数据库初始化成功，随时可以起飞！")
}