✨ 拾光图床 (LightPicker)
**记录光影，极简随心。**专为 1C1G 小主机打造的私人极简图床。Go + Vue3 强强联合，单文件即可部署，Docker 镜像仅 14.4MB！

🎯 核心特性 | Key Features
🚀 极致轻量部署：单文件二进制直接运行，无需安装任何依赖环境。

🗜️ 前端 WebP 引擎：上传前在浏览器端自动压缩，节省带宽。

🛡️ 柔性防盗链神盾：毫秒级内存白名单过滤，防盗图的同时支持直接预览。

🎲 专属盲盒 API：支持按文件夹、横竖屏抽取的随机图片 API。

🎨 极客 UI 矩阵：内置薄荷、樱花、深海、极客 4 套主题，随心切换。

💾 优雅数据备份：SQLite 单文件数据库，支持网页端一键下载备份。

📦 部署指南 | Deployment
方案 A：一键脚本部署 (推荐 Linux 用户)
curl -fsSL https://raw.githubusercontent.com/lychee522/LightPicker/main/install.sh | bash

方案 B：Docker 极速部署
docker run -d

--name lightpicker

-p 5894:5894

-v $(pwd)/storage:/app/storage

--restart always

lycheexiaoxiao/lightpicker:latest

方案 C：二进制部署
请前往 Releases 页面下载对应版本，解压并直接运行即可。

👨‍💻 创始人
TG 账号: @肖肖雨歇

GitHub: lychee522

📄 开源协议
MIT License
