#!/bin/bash
# @author 肖肖雨歇 (tg: @肖肖雨歇)
# ✨ 拾光图床 (LightPicker) Linux 一键极速部署脚本 (Systemd 守护版)

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

# 4. 授权并配置 systemd 守护进程
chmod +x picgo-lite

echo "⚙️ 正在配置 Systemd 守护进程 (开机自启 & 崩溃重启)..."

# 获取当前绝对路径和运行用户，防止写死路径报错
WORK_DIR=$(pwd)
CURRENT_USER=$(whoami)

# 先停止并清理可能存在的旧服务
sudo systemctl stop lightpicker 2>/dev/null
sudo systemctl disable lightpicker 2>/dev/null

# 动态生成 Systemd 服务配置文件
cat <<EOF | sudo tee /etc/systemd/system/lightpicker.service > /dev/null
[Unit]
Description=LightPicker Service - Powered by @肖肖雨歇
After=network.target

[Service]
Type=simple
User=$CURRENT_USER
WorkingDirectory=$WORK_DIR
ExecStart=$WORK_DIR/picgo-lite
Restart=always
RestartSec=5
# 将日志重定向到咱们之前设定的 run.log 中
StandardOutput=append:$WORK_DIR/run.log
StandardError=append:$WORK_DIR/run.log

[Install]
WantedBy=multi-user.target
EOF

echo "🚀 服务启动中..."

# 重新加载 systemd 配置，设置开机自启并立刻启动服务
sudo systemctl daemon-reload
sudo systemctl enable lightpicker
sudo systemctl start lightpicker

sleep 2

# 5. 结果反馈 (改用 systemctl 状态判断)
if sudo systemctl is-active --quiet lightpicker; then
    echo "================================================="
    echo "🎉 部署大功告成！程序已配置为开机自启与自动重启！"
    echo "👉 访问地址: http://你的服务器IP:5894"
    echo "👉 管理日志: tail -f ~/lightpicker/run.log"
    echo "👉 服务状态: sudo systemctl status lightpicker"
    echo "================================================="
else
    echo "❌ 启动失败，请使用命令检查具体报错: sudo systemctl status lightpicker"
fi
