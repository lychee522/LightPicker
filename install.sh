#!/bin/bash
# @author 肖肖雨歇 (tg: @肖肖雨歇)
# ✨ 拾光图床 (LightPicker) Linux 一键极速部署脚本

echo "================================================="
echo "   ✨ 欢迎使用 拾光图床 (LightPicker) 一键部署"
echo "================================================="

# 1. 环境检查：检测系统架构
ARCH=$(uname -m)
if [ "$ARCH" = "x86_64" ]; then
    FILE_NAME="picgo-lite-linux-amd64"
elif [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
    FILE_NAME="picgo-lite-linux-arm64"
else
    echo "❌ 抱歉，暂不支持您的系统架构: $ARCH"
    exit 1
fi

# 2. 创建运行目录
echo "📁 正在准备运行环境..."
mkdir -p ~/lightpicker/storage/uploads
cd ~/lightpicker || exit

# 3. 从 GitHub 获取最新版下载地址
echo "📥 正在下载核心组件: $FILE_NAME ..."
# 直接定位到你的 Release 链接
DOWNLOAD_URL="https://github.com/lychee522/LightPicker/releases/latest/download/${FILE_NAME}"

curl -L -o picgo-lite "$DOWNLOAD_URL"

if [ ! -s picgo-lite ]; then
    echo "❌ 下载失败，请检查网络是否能访问 GitHub Release"
    exit 1
fi

# 4. 授权并启动
chmod +x picgo-lite

echo "🚀 服务启动中..."
# 先杀掉可能存在的旧进程
pkill -f "./picgo-lite" 2>/dev/null

# 使用 nohup 保证后台持续运行
nohup ./picgo-lite > run.log 2>&1 &

sleep 2

# 5. 结果反馈
if pgrep -f "./picgo-lite" > /dev/null; then
    echo "================================================="
    echo "🎉 部署大功告成！"
    echo "👉 访问地址: http://你的服务器IP:5894"
    echo "👉 管理日志: tail -f ~/lightpicker/run.log"
    echo "================================================="
else
    echo "❌ 启动失败，请检查端口 5894 是否被占用。"
fi