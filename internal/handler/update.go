package handler

// @author tg账号的肖肖雨歇
// @description OTA 沙盒升级系统：防伪校验、防降级护盾、平滑替换与自动重启

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// OTAUpdate 处理二进制包的上传与沙盒升级逻辑
func OTAUpdate(c *gin.Context) {
	// 获取是否强制降级标志
	forceDowngrade := c.PostForm("force") == "true"

	// 获取上传的二进制文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "没有找到上传的升级包文件"})
		return
	}

	// 1. 建立沙盒环境 (放置到系统临时目录)
	tempDir := os.TempDir()
	tempExecPath := filepath.Join(tempDir, "picgo_lite_update_temp")

	if err := c.SaveUploadedFile(file, tempExecPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "沙盒环境创建失败: " + err.Error()})
		return
	}
	defer os.Remove(tempExecPath) // 无论成败，运行结束清理沙盒

	// 赋予沙盒程序执行权限
	if err := os.Chmod(tempExecPath, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "沙盒程序赋权失败: " + err.Error()})
		return
	}

	// 2. 防伪校验：在沙盒中执行 ./temp --version
	cmd := exec.Command(tempExecPath, "--version")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "非法的升级包，无法在沙盒中执行环境校验！"})
		return
	}

	output := strings.TrimSpace(out.String())
	magicString := "LightPicker Core Version: "
	if !strings.Contains(output, magicString) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "安全警告：缺少专属标识，拒绝升级！你传的怕不是个病毒吧？"})
		return
	}

	// 提取沙盒包的目标版本号
	newVersion := strings.TrimPrefix(output, magicString)

	// 3. 动态获取当前正在运行的程序版本 (免除硬编码)
	currentExec, err := os.Executable()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "无法定位当前主程序路径"})
		return
	}

	var currentVersion string
	curCmd := exec.Command(currentExec, "--version")
	var curOut bytes.Buffer
	curCmd.Stdout = &curOut
	if curCmd.Run() == nil {
		currentVersion = strings.TrimPrefix(strings.TrimSpace(curOut.String()), magicString)
	}

	// 4. 防降级护盾 (字符串比对)
	if currentVersion != "" && newVersion < currentVersion && !forceDowngrade {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "🛡️ 触发防降级护盾！当前版本 " + currentVersion + "，目标版本 " + newVersion + "。若执意降级可能导致数据库格式损坏，请勾选【强制覆盖】重试！",
		})
		return
	}

	// 5. 备份核心资产 data.db -> data.db.bak
	dbPath := filepath.Join("storage", "data.db")
	bakPath := filepath.Join("storage", fmt.Sprintf("data_%d.db.bak", time.Now().Unix()))
	if err := copyFileInternal(dbPath, bakPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库备份失败，终止升级保平安: " + err.Error()})
		return
	}
	log.Printf("✅ 升级前置：数据库已安全备份至 %s", bakPath)

	// 6. 平滑替换黑魔法 (规避 Text file busy)
	oldExecPath := currentExec + ".old"
	os.Rename(currentExec, oldExecPath) // 移花接木：把正在运行的自己改名挪开
	defer os.Remove(oldExecPath)        // 稍后清理旧躯壳

	// 把沙盒里的新程序复制到主干道
	if err := copyFileInternal(tempExecPath, currentExec); err != nil {
		// 覆盖失败，赶紧把旧躯壳搬回来抢救一下
		os.Rename(oldExecPath, currentExec)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "主程序覆盖失败，已自动回滚: " + err.Error()})
		return
	}
	os.Chmod(currentExec, 0755) // 重新赋予灵魂(权限)

	// 7. 发送成功响应，触发系统自裁重启
	c.JSON(http.StatusOK, gin.H{"message": "🎉 升级包校验通过，数据已备份！系统将在 2 秒后自动重启更新，请稍后刷新页面。"})

	go func() {
		log.Printf("🚀 接收到来自版本 %s 的升级指令，正在准备涅槃重生...", newVersion)
		time.Sleep(2 * time.Second) // 留点时间让前端把 200 响应收完
		os.Exit(0)                  // 拔管自杀，依赖 Systemd 或 Docker Restart=always 将新版本拉起
	}()
}

// copyFileInternal 内部物理复制文件辅助函数
func copyFileInternal(src, dst string) error {
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
