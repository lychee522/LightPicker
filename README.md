# ✨ 拾光图床 (LightPicker)

![Version](https://img.shields.io/badge/version-v1.0.0-blue)
![Docker Size](https://img.shields.io/badge/docker%20size-14MB-success)
![License](https://img.shields.io/badge/license-MIT-green)

> **记录光影，极简随心。** > 专为 1C1G 小主机打造的私人极简图床，Go + Vue3 强强联合，单文件即可部署，Docker 镜像仅 14MB！

## 🎯 核心特性 | Key Features

- **🚀 极致轻量部署**：告别繁琐环境，无论是单文件二进制还是 Docker，皆可秒级启动。
- **🗜️ 独创前端 WebP 引擎**：上传前在浏览器端自动压缩图片，榨干服务器最后一滴带宽。
- **🛡️ 柔性防盗链神盾**：毫秒级内存白名单过滤，彻底告别流量被盗，且支持直接浏览器预览。
- **🎲 专属盲盒 API**：一键生成指定横竖屏、指定文件夹的随机图片 API，支持多端快捷调用。
- **🎨 极客 UI 矩阵**：内置 薄荷清新、樱花粉、深海蓝、极客黑客 4 套精美皮肤，无缝切换。
- **💾 优雅数据备份**：SQLite 单文件数据库，支持网页端一键下载备份与恢复，数据绝对私有。

---

## 📦 部署指南 | Deployment

### 方案 A：一键脚本部署 (推荐 Linux 用户)
只需在你的终端执行以下命令，脚本会自动识别你的系统架构（AMD64/ARM64）并完成下载、配置与后台运行：
```bash
curl -fsSL [https://raw.githubusercontent.com/lychee522/LightPicker/main/install.sh](https://raw.githubusercontent.com/lychee522/LightPicker/main/install.sh) | bash

### 方案 B：Docker 极速部署 (全平台通用)
镜像已推送到 Docker Hub，包含环境仅 14MB，极其省资源。

Bash
docker run -d \
  --name lightpicker \
  -p 5894:5894 \
  -v $(pwd)/storage:/app/storage \
  --restart always \
  lycheexiaoxiao/lightpicker:latest

### 方案 C：全系统独立二进制部署 (Windows / macOS)
如果您不想使用任何命令行工具，我们为您准备了开箱即用的二进制文件：

前往 Releases 页面下载对应您系统的版本。

在程序同级目录下创建 storage/uploads 文件夹。

双击运行 picgo-lite-xxx。

浏览器访问 http://127.0.0.1:5894 即可使用。

👨‍💻 创始人
TG 账号：@肖肖雨歇
