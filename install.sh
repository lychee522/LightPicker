#!/bin/bash
# @author 肖肖雨歇 (tg: @肖肖雨歇)
# ✨ 拾光图床 (LightPicker) Linux 一键极速部署脚本 (全网满速终极版)

echo "================================================="
echo "   ✨ 欢迎使用 拾光图床 (LightPicker) 一键部署"
echo "================================================="

# 【重要配置】每次发布新版本，只需修改这里的版本号！
VERSION="v1.0.0"

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

# 3. 核心下载逻辑：Gitee 满速首选 + 多节点备胎
echo "📥 正在获取核心组件: $FILE_NAME ($VERSION) ..."

# 定义各个源的完整下载链接
GITEE_URL="https://gitee.com/lychee522/LightPicker/releases/download/${VERSION}/${FILE_NAME}"
GITHUB_URL="https://github.com/lychee522/LightPicker/releases/download/${VERSION}/${FILE_NAME}"

# 下载链接池：按优先级排序，Gitee 永远排第一！
URL_POOL=(
    "$GITEE_URL"
    "https://ghp.ci/$GITHUB_URL"
    "https://mirror.ghproxy.com/$GITHUB_URL"
    "https://fastgh.oso.gs/$GITHUB_URL"
    "$GITHUB_URL"
)

DOWNLOAD_SUCCESS=0

# 强行清理掉之前残留的任何同名垃圾文件
rm -f picgo-lite

for URL in "${URL_POOL[@]}"; do
    if [[ "$URL" == "$GITEE_URL" ]]; then
        echo "🚀 VIP通道: 尝试通过 Gitee 国内直链全速下载..."
    elif [[ "$URL" == "$GITHUB_URL" ]]; then
        echo "🔄 终极保底: 尝试 GitHub 官方源直连..."
    else
        echo "🔄 备用通道: 尝试 GitHub 加速节点下载..."
    fi

    # -f 拦截错误网页, -L 跟随重定向, 10秒连不上放弃, 3分钟下不完放弃
    if curl -fL --connect-timeout 10 -m 180 -o picgo-lite "$URL"; then
        if [ -s picgo-lite ]; then
            echo "✅ 下载成功！这速度绝对起飞！"
            DOWNLOAD_SUCCESS=1
            break
        fi
    fi

    echo "⚠️ 当前通道不通畅或返回错误，清理现场，自动切换下一通道..."
    rm -f picgo-lite 
done

if [ $DOWNLOAD_SUCCESS -eq 0 ]; then
    echo "❌ 所有下载通道均已失效！请检查服务器网络。"
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
