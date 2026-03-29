# ✨ 拾光图床 (LightPicker)

![Version](https://img.shields.io/badge/version-v1.0.0-blue)
![Docker Size](https://img.shields.io/badge/docker%20size-14.4MB-success)
![Go](https://img.shields.io/badge/language-Go-00ADD8)
![Vue](https://img.shields.io/badge/framework-Vue3-4FC08D)
![License](https://img.shields.io/badge/license-MIT-green)

> **记录光影，极简随心。** > 专为 1C1G 小主机打造的私人极简图床。Go + Vue3 强强联合，单文件即可部署，Docker 镜像仅 **14.4MB**！

---

## 🎯 核心特性 | Key Features

- **🚀 极致轻量部署**：单文件二进制直接运行，无需安装任何依赖环境。
- **🗜️ 前端 WebP 引擎**：上传前自动压缩，节省 70% 以上的存储空间与带宽。
- **🛡️ 柔性防盗链神盾**：毫秒级内存白名单过滤，拦截恶意盗链的同时支持直接预览。
- **🎲 专属盲盒 API**：支持按文件夹、横竖屏抽取的随机图片 API，一键生成代码。
- **🎨 极客 UI 矩阵**：内置 **薄荷、樱花、深海、极客** 4 套精美主题，随心切换。
- **💾 优雅数据备份**：SQLite 单文件数据库，支持网页端一键下载与恢复，数据完全私有。

---

## 📦 部署指南 | Deployment

### 方案 A：一键脚本部署 (推荐 Linux 用户)
只需在你的终端执行以下命令，脚本会自动识别架构并完成部署：

```bash
curl -fsSL [https://raw.githubusercontent.com/lychee522/LightPicker/main/install.sh](https://raw.githubusercontent.com/lychee522/LightPicker/main/install.sh) | bash


docker run -d \
  --name lightpicker \
  -p 5894:5894 \
  -v $(pwd)/storage:/app/storage \
  --restart always \
  lycheexiaoxiao/lightpicker:latest
