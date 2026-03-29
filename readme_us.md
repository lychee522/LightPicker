## 📄 License

MIT License

---

# 🌍 README_US（复制到 README_US.md）

````md
# ✨ LightPicker

> **Capture light. Keep it simple.**  
> A minimal self-hosted image hosting solution designed for low-resource servers (1C1G). Built with Go + Vue3. Single binary deployment. Docker image only **14.4MB**.

---

## 🖼️ Preview

### Gallery
![gallery](https://github.com/user-attachments/assets/e29e99ac-7ac5-4773-ab59-eb1797357708)

### Upload
![upload](https://github.com/user-attachments/assets/27adc465-1d59-4a0d-a60c-59b71aaf90d5)

### Settings
![settings](https://github.com/user-attachments/assets/c7da26e2-b3b7-4821-8fbc-a22a82c24b84)

---

## 🎯 Features

- 🚀 **Ultra Lightweight**  
  Single binary, no dependencies required.

- 🗜️ **Frontend WebP Compression**  
  Compress images in-browser before upload, saving **70%+** storage and bandwidth.

- 🛡️ **Hotlink Protection**  
  In-memory whitelist filtering.

- 🎲 **Random Image API**  
  Fetch random images by folder or orientation.

- 🎨 **Multiple Themes**  
  Mint / Sakura / DeepSea / Geek.

- 💾 **Backup Friendly**  
  SQLite single-file database.

---

## 📦 Deployment

### A. Script (Linux)

```bash
curl -fsSL https://raw.githubusercontent.com/lychee522/LightPicker/main/install.sh | bash
````

### B. Docker

```bash
docker run -d \
  --name lightpicker \
  -p 5894:5894 \
  -v $(pwd)/storage:/app/storage \
  --restart always \
  lycheexiaoxiao/lightpicker:latest
```

### C. Binary (Windows / macOS)

1. Download from Releases
2. Create folder: `storage/uploads`
3. Run executable
4. Open [http://127.0.0.1:5894](http://127.0.0.1:5894)

---

## ⭐ Support

* ⭐ Star
* 🍴 Fork
* 📢 Share

```

---

## ✅ 使用说明

把这三个占位符替换掉即可：

- `{{DEMO_GALLERY_URL}}`
- `{{DEMO_UPLOAD_URL}}`
- `{{DEMO_SETTINGS_URL}}`

直接粘你图床 URL 就能用了 👍

```
