# ✨ 拾光图床 (LightPicker)

![Version](https://img.shields.io/badge/version-v1.0.0-blue)
![Docker Size](https://img.shields.io/badge/docker%20size-14.4MB-success)
![License](https://img.shields.io/badge/license-MIT-green)

> **记录光影，极简随心。**
> 专为 1C1G 小主机打造的私人极简图床。Go + Vue3 强强联合，单文件即可部署，Docker 镜像仅 **14.4MB**！

---

## 🖼️ 演示截图

### 📚 图库界面

![gallery](https://github.com/user-attachments/assets/e29e99ac-7ac5-4773-ab59-eb1797357708)

### 📤 上传界面

![upload](https://github.com/user-attachments/assets/27adc465-1d59-4a0d-a60c-59b71aaf90d5)

### ⚙️ 系统配置

![settings](https://github.com/user-attachments/assets/c7da26e2-b3b7-4821-8fbc-a22a82c24b84)

---
## 🎯 核心特性 | Key Features

* 🚀 **极致轻量部署**
  单文件二进制直接运行，无需安装任何依赖环境。

* 🗜️ **前端 WebP 引擎**
  上传前在浏览器端自动压缩，节省 **70%+** 的存储空间与带宽。

* 🛡️ **柔性防盗链神盾**
  毫秒级内存白名单过滤，防盗图的同时支持直接预览。

* 🎲 **专属盲盒 API**
  支持按文件夹、横竖屏抽取的随机图片 API。

* 🎨 **极客 UI 矩阵**
  内置「薄荷 / 樱花 / 深海 / 极客」四套主题，随心切换。

* 💾 **优雅数据备份**
  SQLite 单文件数据库，支持网页端一键下载与恢复。

---

## 📦 部署指南 | Deployment

### 🅰️ 方案 A：一键脚本部署（推荐 Linux 用户）

```bash
curl -fsSL https://raw.githubusercontent.com/lychee522/LightPicker/main/install.sh | bash
```

---

### 🅱️ 方案 B：Docker 极速部署（全平台通用）

```bash
docker run -d \
  --name lightpicker \
  -p 5894:5894 \
  -v $(pwd)/storage:/app/storage \
  --restart always \
  lycheexiaoxiao/lightpicker:latest
```

---

### 🅲 方案 C：独立二进制部署（Windows / macOS）

1. 前往 **Releases** 页面下载对应系统版本。
2. 在程序同级目录创建：

   ```
   storage/uploads
   ```
3. 双击运行程序（如：`picgo-lite-xxx`）。
4. 浏览器访问：

   ```
   http://127.0.0.1:5894
   ```

---

## 👨‍💻 创始人

* Telegram：[@肖肖雨歇](https://t.me/x9426y9464)
* GitHub：[lychee522](https://github.com/lychee522)

---

## 📄 开源协议

本项目基于 **MIT License** 开源。

---

## ⭐ 支持项目

如果这个项目对你有帮助，欢迎：

* ⭐ 点个 Star
* 🍴 Fork 使用
* 📢 分享给更多人

你的支持就是作者持续更新的最大动力！ 🚀
