<template>
  <div class="app-container" :data-theme="currentTheme">
    <div class="toast-container">
      <transition-group name="toast-slide">
        <div v-for="t in toasts" :key="t.id" class="toast-msg" :class="t.type">{{ t.icon }} {{ t.msg }}</div>
      </transition-group>
    </div>

    <div class="login-layout" v-if="appState === 'init' || appState === 'login'">
      <div class="card login-card">
        <div class="logo-area">✨</div>
        <h1 class="title">拾光图床</h1>
        <p class="subtitle" v-if="appState === 'init'">欢迎，请设置超级管理员账号</p>
        <p class="subtitle" v-else>记录光影，极简随心</p>

        <div class="form-area" v-if="appState === 'init'">
          <input v-model="form.username" type="text" placeholder="设置账号" class="input-box" />
          <input v-model="form.password" type="password" placeholder="设置密码" class="input-box" />
          <button @click="handleInit" class="btn primary-btn" style="width: 100%">立即初始化</button>
        </div>
        <div class="form-area" v-else>
          <input v-model="form.username" type="text" placeholder="账号" class="input-box" />
          <input v-model="form.password" type="password" placeholder="密码" class="input-box" />
          <button @click="handleLogin" class="btn primary-btn" style="width: 100%">安全登录</button>
        </div>
      </div>
    </div>

    <div class="admin-layout" v-else-if="appState === 'dashboard'">
      <div class="sidebar">
        <div class="sidebar-header"><div class="logo-icon">✨</div><h2 class="logo-text">拾光</h2></div>
        <ul class="menu">
          <li @click="switchTab('upload')" :class="{active: currentTab === 'upload'}">📤 上传大厅</li>
          <li @click="switchTab('gallery')" :class="{active: currentTab === 'gallery'}">🖼️ 拾光图库</li>
          <li @click="switchTab('settings')" :class="{active: currentTab === 'settings'}">⚙️ 系统配置</li>
        </ul>
        <div class="sidebar-footer"><button @click="logout" class="btn danger-btn logout-btn">🚪 安全退出</button></div>
      </div>
      <div class="main-content">
        <UploadDashboard v-if="currentTab === 'upload'" />
        <GalleryView v-if="currentTab === 'gallery'" />
        <SettingsView v-if="currentTab === 'settings'" @change-theme="updateTheme" />
      </div>
    </div>
  </div>
</template>

<script setup>
// @author 肖肖雨歇 - 修复白屏死机的完整 App.vue
import { ref, onMounted } from 'vue'
import axios from 'axios'
import UploadDashboard from './components/UploadDashboard.vue'
import GalleryView from './components/GalleryView.vue'
import SettingsView from './components/SettingsView.vue'

const toasts = ref([])
window.$toast = (msg, type = 'success') => {
  const id = Date.now(); const icon = type === 'success' ? '🎉' : type === 'error' ? '❌' : 'ℹ️'
  toasts.value.push({ id, msg, type, icon })
  setTimeout(() => { toasts.value = toasts.value.filter(t => t.id !== id) }, 3000)
}
window.alert = (msg) => window.$toast(msg, 'info')

const appState = ref('login')
// 🌟 核心：从缓存读取当前 Tab，让你刷新不迷路
const currentTab = ref(localStorage.getItem('picgo_tab') || 'gallery')
const currentTheme = ref(localStorage.getItem('app_theme') || 'fresh')
const form = ref({ username: '', password: '' })

// 🌟 核心：切换 Tab 时自动保存到本地缓存
const switchTab = (tab) => {
  currentTab.value = tab
  localStorage.setItem('picgo_tab', tab)
}

const updateTheme = (theme) => { currentTheme.value = theme; localStorage.setItem('app_theme', theme) }

onMounted(async () => {
  const token = localStorage.getItem('picgo_token')
  if (token) appState.value = 'dashboard'
  try { await axios.get('/api/ping') } catch (error) { if (error.response?.data?.action === 'redirect_to_init') appState.value = 'init' }
})

const handleInit = async () => {
  if (!form.value.username || !form.value.password) return window.$toast('账号密码不能为空！', 'error')
  try { const res = await axios.post('/api/init', form.value); window.$toast(res.data.message, 'success'); appState.value = 'login' } catch (err) { window.$toast('初始化失败', 'error') }
}
const handleLogin = async () => {
  try { const res = await axios.post('/api/login', form.value); window.$toast('登录成功！', 'success'); localStorage.setItem('picgo_token', res.data.token); appState.value = 'dashboard' } catch (err) { window.$toast(err.response?.data?.error || '登录失败', 'error') }
}
const logout = () => {
  if(!confirm('确定要退出吗？')) return; localStorage.removeItem('picgo_token'); appState.value = 'login'; window.$toast('已安全退出', 'success')
}
</script>

<style scoped>
.login-layout { width: 100vw; height: 100vh; display: flex; justify-content: center; align-items: center; background: var(--bg-color); }
.login-card { width: 340px; border-radius: 20px; box-shadow: var(--shadow-md); display: flex; flex-direction: column; align-items: center; background: var(--card-bg); padding: 40px;}
.logo-area { font-size: 48px; margin-bottom: 10px; }
.title { color: var(--text-main); margin: 0 0 10px 0; font-size: 24px; font-weight: 800;}
.subtitle { color: var(--text-desc); font-size: 14px; margin-bottom: 30px;}
.form-area { width: 100%; display: flex; flex-direction: column; gap: 15px;}
.input-box { width: 100%; padding: 10px; border: 1px solid var(--border-color); border-radius: 6px; background: var(--card-bg); color: var(--text-main); box-sizing: border-box;}

.admin-layout { display: flex; flex-direction: row; width: 100vw; height: 100vh; background: var(--bg-color); overflow: hidden; }
.sidebar { width: 240px; flex-shrink: 0; border-right: 1px solid var(--border-color); background: var(--sidebar-bg); padding: 0; display: flex; flex-direction: column; justify-content: space-between; box-shadow: var(--shadow-sm); z-index: 10;}
.sidebar-header { padding: 30px 0; text-align: center; }
.logo-icon { font-size: 32px; margin-bottom: 5px;}
.logo-text { margin: 0; color: var(--accent-color); font-weight: 900; letter-spacing: 2px;}
.menu { padding: 0 15px; flex: 1; list-style: none; margin: 0;}
.menu li { margin-bottom: 8px; border-radius: 10px; border: none; font-weight: 600; color: var(--sidebar-text); padding: 12px 20px; cursor: pointer; transition: 0.2s;}
.menu li.active { background: var(--accent-color); color: #fff; box-shadow: 0 4px 12px rgba(16, 185, 129, 0.2); }
.menu li:hover:not(.active) { background: var(--bg-color); color: var(--accent-color); }
.sidebar-footer { padding: 20px; }
.logout-btn { width: 100%; border-radius: 10px; border: 1.5px solid var(--border-color); color: var(--text-desc); background: transparent; cursor: pointer; padding: 10px; font-weight: bold; transition: 0.2s;}
.logout-btn:hover { background: #ef4444; color: #fff; border-color: #ef4444;}
.main-content { flex: 1; overflow-y: auto; padding: 30px 40px; }

.toast-container { position: fixed; top: 20px; left: 50%; transform: translateX(-50%); z-index: 9999; display: flex; flex-direction: column; gap: 10px; align-items: center; pointer-events: none;}
.toast-msg { background: var(--card-bg); color: var(--text-main); padding: 12px 24px; border-radius: 50px; font-weight: 600; font-size: 14px; box-shadow: 0 10px 30px rgba(0,0,0,0.1); display: flex; align-items: center; gap: 8px; border: 1px solid var(--border-color);}
.toast-msg.error { color: #ef4444; border-color: #fca5a5; background: #fef2f2;}
.toast-msg.success { color: var(--accent-color); border-color: #a7f3d0; background: #ecfdf5;}
.toast-slide-enter-active, .toast-slide-leave-active { transition: all 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275); }
.toast-slide-enter-from { opacity: 0; transform: translateY(-30px) scale(0.9); }
.toast-slide-leave-to { opacity: 0; transform: translateY(-20px) scale(0.9); }
</style>