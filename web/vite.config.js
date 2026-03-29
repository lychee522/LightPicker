// @author 肖肖雨歇
// @description 前端构建与代理配置

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    host: '0.0.0.0', // 允许局域网访问
    proxy: {
      // 凡是 /api 开头的请求，全部转发给咱们的 Go 后端
      '/api': {
        target: 'http://127.0.0.1:5894',
        changeOrigin: true,
      },
      // 图片资源的请求也转发给后端
      '/uploads': {
        target: 'http://127.0.0.1:5894',
        changeOrigin: true,
      }
    }
  }
})