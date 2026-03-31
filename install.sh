#!/bin/bash
# @author 肖肖雨歇 (tg: @肖肖雨歇)
# ✨ 拾光图床 (LightPicker) Linux 一键极速部署脚本 (多节点自动故障转移版)

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

# 3. 核心下载逻辑：多节点自动切换
echo "📥 正在获取核心组件: $FILE_NAME ..."

# 基础下载直链
RAW_URL="https://github.com/lychee522/LightPicker/releases/latest/download/${FILE_NAME}"

# 定义多个真实可用的加速节点池，最后一个留空代表使用官方直连
PROXIES=(
    "https://ghp.ci/"
    "https://github.moeyy.xyz/"
    "https://fastgh.oso.gs/"
    "https://mirror.ghproxy.com/"
    "" 
)

DOWNLOAD_SUCCESS=0

# 遍历尝试所有节点
for PROXY in "${PROXIES[@]}"; do
    if [ -z "$PROXY" ]; then
        echo "🔄 尝试直连 GitHub 官方源下载..."
        DOWNLOAD_URL="$RAW_URL"
    else
        echo "🔄 尝试使用加速节点下载: $PROXY ..."
        DOWNLOAD_URL="${PROXY}${RAW_URL}"
    fi

    # 使用 curl 下载:
    # --connect-timeout 10 : 10秒连不上直接放弃换下一个
    # -m 180 : 最多允许下载 3 分钟，防止龟速卡死
    curl -L --connect-timeout 10 -m 180 -o picgo-lite "$DOWNLOAD_URL"

    # 检查文件是否下载成功且大小不为0
    if [ -s picgo-lite ]; then
        echo "✅ 下载成功！"
        DOWNLOAD_SUCCESS=1
        break
    else
        echo "⚠️ 当前节点不可用或超时，正在清理残缺文件，准备切换下一个..."
        rm -f picgo-lite 
    fi
done

if [ $DOWNLOAD_SUCCESS -eq 0 ]; then
    echo "❌ 所有下载节点均已失效或超时！"
    echo "建议：检查服务器网络，或尝试手动下载核心组件并放入 ~/lightpicker 目录。"
    exit 1
fi

# 4. 授权并配置 systemd 守护进程
chmod +x picgo-lite
echo "⚙️ 正在配置 Systemd 守护进程 (开机自启 & 崩溃重启)..."

WORK_DIR=$(pwd)
CURRENT_USER=$(whoami)

sudo systemctl stop lightpicker 2>/dev/null
sudo systemctl disable lightpicker 2>/dev/null

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
StandardOutput=append:$WORK_DIR/run.log
StandardError=append:$WORK_DIR/run.log

[Install]
WantedBy=multi-user.target
EOF

echo "🚀 服务启动中..."
sudo systemctl daemon-reload
sudo systemctl enable lightpicker
sudo systemctl start lightpicker
sleep 2

# 5. 结果反馈
if sudo systemctl is-active --quiet lightpicker; then
    echo "================================================="
    echo "🎉 部署大功告成！程序已配置为开机自启与自动重启！"
    echo "👉 访问地址: http://你的服务器IP:5894"
    echo "👉 管理日志: tail -f ~/lightpicker/run.log"
    echo "================================================="
else
    echo "❌ 启动失败，请使用命令检查具体报错: sudo systemctl status lightpicker"
fi
