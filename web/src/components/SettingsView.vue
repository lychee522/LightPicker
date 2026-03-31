<template>
  <div class="settings-container" @click="closeMenu">
    <div class="header"><h2>⚙️ 系统配置中心</h2></div>
    <div class="settings-layout">
      <div class="settings-sidebar">
        <div class="menu-item" :class="{ active: activeTab === 'theme' }" @click="activeTab = 'theme'">🎨 界面配置</div>
        <div class="menu-item" :class="{ active: activeTab === 'domain' }" @click="activeTab = 'domain'">🌐 域名与防盗链</div>
        <div class="menu-item" :class="{ active: activeTab === 'api' }" @click="activeTab = 'api'">🎲 专属盲盒 API</div>
        <div class="menu-item" :class="{ active: activeTab === 'backup' }" @click="activeTab = 'backup'">💾 数据备份</div>
        <div class="menu-item" :class="{ active: activeTab === 'about' }" @click="activeTab = 'about'">ℹ️ 关于拾光</div>
      </div>

      <div class="settings-content">
        <div v-if="activeTab === 'theme'" class="panel">
          <h3>🎨 皮肤样式选择</h3>
          <div class="theme-grid">
            <div class="theme-card" :class="{ active: currentTheme === 'fresh' }" @click="selectTheme('fresh')"><div class="theme-preview" style="background: #f4f7f6; border: 1px solid #10b981;"></div><p>薄荷清新 (默认)</p></div>
            <div class="theme-card" :class="{ active: currentTheme === 'sakura' }" @click="selectTheme('sakura')"><div class="theme-preview" style="background: #fff0f3; border: 1px solid #fb6f92;"></div><p>樱花粉</p></div>
            <div class="theme-card" :class="{ active: currentTheme === 'ocean' }" @click="selectTheme('ocean')"><div class="theme-preview" style="background: #f0f4f8; border: 1px solid #3b82f6;"></div><p>深海蓝</p></div>
            <div class="theme-card" :class="{ active: currentTheme === 'geek' }" @click="selectTheme('geek')"><div class="theme-preview" style="background: #0d1117; border: 1px solid #00ff41;"></div><p>极客黑客</p></div>
          </div>
        </div>

        <div v-if="activeTab === 'domain'" class="panel">
          <div class="panel-header" style="display: flex; justify-content: space-between;">
            <h3>🌐 自定义外网域名池</h3><button @click="addDomain" class="btn primary-btn small-btn">➕ 新增域名</button>
          </div>
          <div class="matrix-list" style="margin-bottom: 30px;">
            <div v-for="(dom, index) in domains" :key="index" class="card-box" style="display: flex; gap: 10px; align-items: center; margin-bottom: 10px;">
              <span style="font-size: 20px;">🔗</span>
              <input v-model="domains[index]" @change="saveDomains" class="form-input" style="flex: 1; margin: 0;" placeholder="如: https://img.yourdomain.com" />
              <button @click="removeDomain(index)" class="btn danger-btn small-btn">删除</button>
            </div>
          </div>

          <div class="panel-header" style="display: flex; justify-content: space-between; align-items: center;">
            <h3 style="margin:0;">🛡️ 防盗链白名单</h3>
            <div>
              <button @click="addWhitelist" class="btn outline-btn small-btn" style="margin-right: 10px;">➕ 新增规则</button>
              <button @click="saveWhitelist" class="btn primary-btn small-btn">💾 保存神盾配置</button>
            </div>
          </div>
          <p class="desc">设置允许直接引用本站图片的域名。列表为空则不设防盗链限制。</p>
          <div class="matrix-list">
            <div v-for="(rule, index) in whitelistArray" :key="index" class="card-box" style="display: flex; gap: 10px; align-items: center; margin-bottom: 10px; padding: 10px 15px;">
              <span style="font-size: 16px;">🛡️</span>
              <input v-model="whitelistArray[index]" class="form-input" style="flex: 1; margin: 0;" placeholder="如: github.com" />
              <button @click="removeWhitelist(index)" class="btn danger-btn small-btn">删除</button>
            </div>
            <div v-if="whitelistArray.length === 0" class="empty-state">当前白名单为空，防盗链处于关闭状态。</div>
          </div>
        </div>

        <div v-if="activeTab === 'api'" class="panel">
          <div class="panel-header" style="display: flex; justify-content: space-between;"><h3>🎲 API 矩阵</h3><button @click="addConfig" class="btn primary-btn small-btn">➕ 新建</button></div>
          <div class="matrix-list">
            <div v-for="(cfg, index) in configs" :key="cfg.id" class="card-box" style="margin-bottom: 15px;">
              <div style="display: flex; justify-content: space-between; margin-bottom: 10px;">
                <input v-model="cfg.name" @change="saveConfigs" class="title-input" placeholder="输入备注..."/>
                <button @click="removeConfig(index)" class="btn danger-btn small-btn">删除</button>
              </div>
              <div style="display: flex; gap: 10px; margin-bottom: 10px;">
                <select v-model="cfg.album" @change="saveConfigs" class="form-select"><option value="0">📂 全部</option><option v-for="a in albums" :value="a.id">📁 {{ a.name }}</option></select>
                <select v-model="cfg.ori" @change="saveConfigs" class="form-select"><option value="all">不限</option><option value="landscape">仅横屏</option><option value="portrait">仅竖屏</option></select>
                <select v-model="cfg.domain" @change="saveConfigs" class="form-select"><option value="default">🌐 默认来源</option><option v-for="d in domains" :value="d">🚀 {{ d }}</option></select>
              </div>
              
              <div style="display: flex; gap: 10px; align-items: center;">
                <code class="api-link" @click="copy(getUrl(cfg))" title="点击直接复制链接" style="cursor: pointer; flex: 1; margin: 0; display: block; background: #282c34; color: #98c379; padding: 10px; border-radius: 6px; font-size: 12px; overflow-x: auto; white-space: nowrap; transition: 0.2s;">
                  {{ getUrl(cfg) }}
                </code>
                <button @click="copy(`![盲盒图](${getUrl(cfg)})`)" class="btn outline-btn small-btn" style="padding: 10px 15px;">📝 复制 MD</button>
              </div>
            </div>
          </div>
        </div>

        <div v-if="activeTab === 'backup'" class="panel">
          <h3>💾 数据库备份与恢复</h3>
          <div class="card-box" style="display: flex; gap: 20px; align-items: center;">
            <button @click="downloadBackup" class="btn primary-btn">📥 下载 data.db</button>
            <button @click="$refs.restoreInput.click()" class="btn danger-btn">📤 上传恢复数据库</button>
            <input type="file" ref="restoreInput" accept=".db" @change="uploadRestore" hidden />
          </div>
        </div>

        <div v-if="activeTab === 'about'" class="panel about-panel">
          <div class="logo-area">✨ 拾光图床</div>
          <h3>LightPicker {{ currentLocalVersion }}</h3>
          
          <div class="update-section card-box" style="margin-bottom: 20px; background: rgba(16, 185, 129, 0.05);">
            <div v-if="updateStatus === 'idle'">
              <button @click="manualCheckUpdate" class="btn primary-btn" :disabled="isChecking">
                {{ isChecking ? '🛰️ 正在巡检云端...' : '🔍 检查在线更新' }}
              </button>
            </div>
            <div v-else-if="updateStatus === 'found'" style="display: flex; flex-direction: column; gap: 10px; align-items: center;">
              <span style="color: #10b981; font-weight: bold;">🎉 发现新版本 {{ remoteVersion }}！</span>
              <button @click="showModal = true" class="btn success-btn">🚀 查看更新详情</button>
            </div>
            <div v-else-if="updateStatus === 'upgrading'" style="width: 100%;">
              <p style="margin-bottom: 10px;">⚡ 补丁拉取中: {{ upgradeProgress }}%</p>
              <div class="progress-bar-container">
                <div class="progress-bar-fill animated" :style="{ width: upgradeProgress + '%' }"></div>
              </div>
            </div>
            <div v-else-if="updateStatus === 'ready_to_restart'" style="display: flex; flex-direction: column; gap: 10px; align-items: center;">
              <span style="color: #3b82f6; font-weight: bold;">✅ 补丁已就绪！</span>
              <button @click="handleRestart" class="btn primary-btn" style="background: #3b82f6;">🔄 立即重启应用新功能</button>
            </div>
          </div>

          <ul class="about-list">
            <li><strong>👨‍💻 创始人：</strong> @肖肖雨歇</li>
            <li><strong>🎯 理念：</strong> 1C1G 零负担，数据完全私有。</li>
            <li><strong>🔑 提示：</strong> 忘记密码？可在终端执行 <code>./picgo-lite admin reset 123</code></li>
          </ul>
        </div>
      </div>
    </div>

    <div v-if="showModal" class="update-modal-mask">
      <div class="update-modal-card">
        <div class="modal-header">
          <span class="version-badge">New v{{ remoteVersion }}</span>
          <h3>🚀 发现拾光系统补丁！</h3>
        </div>
        <div class="modal-body">
          <p class="section-title">📦 更新日志：</p>
          <div class="changelog-container">{{ remoteChangelog }}</div>
          <p class="note">注意：升级会自动重启服务，请避开业务高峰。数据已做安全保护，不会丢失。</p>
        </div>
        <div class="modal-footer">
          <button @click="startUpgradeProcess" class="btn primary-btn">立即更新</button>
          <button @click="skipCurrentVersion" class="btn outline-btn">跳过此版本</button>
          <button @click="showModal = false" class="btn ghost-btn">下次再说</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import axios from 'axios'
const emit = defineEmits(['change-theme'])

// --- 基础状态 ---
const activeTab = ref('about')
const currentTheme = ref(localStorage.getItem('app_theme') || 'fresh')
const albums = ref([]); const configs = ref([]); 
const domains = ref(JSON.parse(localStorage.getItem('picgo_domains') || '[]'))
const whitelistArray = ref([]) 
const host = window.location.origin
const authHeader = () => ({ headers: { Authorization: `Bearer ${localStorage.getItem('picgo_token')}` } })

// --- 🌟 OTA 升级核心状态 (修复比对逻辑) ---
const currentLocalVersion = ref('v1.2.0') // 🌟 必须定义这个变量！
const updateStatus = ref('idle') // idle, found, upgrading, ready_to_restart
const upgradeProgress = ref(0)
const isChecking = ref(false)
const showModal = ref(false)
const remoteVersion = ref('')
const remoteChangelog = ref('')
let checkTimer = null
let progressTimer = null

// 🌟 语义化版本比对：返回 true 代表 remote 更大
const isNewer = (local, remote) => {
  const l = local.replace('v', '').split('.').map(Number)
  const r = remote.replace('v', '').split('.').map(Number)
  for (let i = 0; i < Math.max(l.length, r.length); i++) {
    const lV = l[i] || 0
    const rV = r[i] || 0
    if (rV > lV) return true
    if (rV < lV) return false
  }
  return false
}

// --- 逻辑：巡检云端版本 ---
const performUpdateCheck = async (isManual = false) => {
  if (isChecking.value) return
  isChecking.value = true
  try {
    const { data } = await axios.get('/api/system/update-check', authHeader())
    // 🌟 只有后端说有新版本，且前端比对确实更大时才弹窗
    if (data.has_new && isNewer(currentLocalVersion.value, data.version)) {
      const skipped = localStorage.getItem('ota_skipped_version')
      if (data.version !== skipped || isManual) {
        remoteVersion.value = data.version
        remoteChangelog.value = data.changelog
        updateStatus.value = 'found'
        if (data.version !== skipped) showModal.value = true
      }
    } else if (isManual) {
      window.$toast('当前已是最新版，真稳！', 'success')
    }
  } catch (e) {
    if (isManual) window.$toast('连接云端失败', 'error')
  } finally {
    isChecking.value = false
  }
}

// --- 逻辑：一键立即更新 ---
const startUpgradeProcess = async () => {
  showModal.value = false
  updateStatus.value = 'upgrading'
  upgradeProgress.value = 0
  try {
    await axios.post('/api/system/upgrade-exec', {}, authHeader())
    startPollingProgress()
  } catch (e) {
    window.$toast('升级指令下发失败', 'error')
    updateStatus.value = 'found'
  }
}

// --- 逻辑：轮询进度 ---
const startPollingProgress = () => {
  progressTimer = setInterval(async () => {
    try {
      const { data } = await axios.get('/api/system/upgrade-progress', authHeader())
      upgradeProgress.value = data.progress
      if (data.status === 'ready_to_restart') {
        clearInterval(progressTimer)
        updateStatus.value = 'ready_to_restart'
        window.$toast('升级包下载完成，准备重启！', 'success')
      }
    } catch (e) { /* 轮询容错 */ }
  }, 1000)
}

// --- 逻辑：跳过此版本 ---
const skipCurrentVersion = () => {
  localStorage.setItem('ota_skipped_version', remoteVersion.value)
  showModal.value = false
  updateStatus.value = 'idle'
  window.$toast('已跳过此版本', 'info')
}

// --- 逻辑：重启 ---
const handleRestart = () => {
  if (confirm('确认重启服务？')) {
    window.$toast('重启中...', 'info')
    setTimeout(() => window.location.reload(), 3000)
  }
}

const manualCheckUpdate = () => performUpdateCheck(true)

// --- 基础功能逻辑 (全量保留) ---
const selectTheme = (theme) => { currentTheme.value = theme; emit('change-theme', theme) }
const loadAlbums = async () => { try { const {data} = await axios.get('/api/albums', authHeader()); albums.value = data.data } catch(e){} }
const saveDomains = () => { localStorage.setItem('picgo_domains', JSON.stringify(domains.value)) }
const addDomain = () => { domains.value.push(''); saveDomains() }
const removeDomain = (idx) => { domains.value.splice(idx, 1); saveDomains() }
const loadWhitelist = async () => { 
  try { 
    const {data} = await axios.get('/api/whitelist', authHeader())
    if(data.data) whitelistArray.value = data.data.split(',').filter(i => i.trim() !== '') 
  } catch(e){} 
}
const addWhitelist = () => { whitelistArray.value.push('') }
const removeWhitelist = (idx) => { whitelistArray.value.splice(idx, 1) }
const saveWhitelist = async () => { 
  const str = whitelistArray.value.filter(i => i.trim() !== '').join(',')
  try { await axios.post('/api/whitelist', { value: str }, authHeader()); window.$toast('已生效！', 'success') } catch(e){ window.$toast('失败', 'error') } 
}
const downloadBackup = () => { window.open('/api/backup?token=' + localStorage.getItem('picgo_token')) }
const uploadRestore = async (e) => {
  const file = e.target.files[0]; if (!file || !confirm('确认覆盖？')) return;
  const fd = new FormData(); fd.append('file', file)
  try { await axios.post('/api/restore', fd, authHeader()); window.$toast('成功', 'success'); setTimeout(() => window.location.reload(), 1000) } catch (err) { window.$toast('失败', 'error') }
}
const getUrl = (cfg) => { 
  const p = new URLSearchParams(); if(cfg.ori!=='all') p.append('ori', cfg.ori); if(cfg.album!=='0') p.append('album_id', cfg.album); 
  const base = (cfg.domain && cfg.domain !== 'default') ? cfg.domain : host;
  return `${base}/api/random${p.toString()?'?'+p.toString():''}` 
}
const saveConfigs = () => localStorage.setItem('picgo_api_configs', JSON.stringify(configs.value))
const addConfig = () => { configs.value.unshift({ id: Date.now(), name: '新配置', album: '0', ori: 'all', domain: 'default' }); saveConfigs() }
const removeConfig = (i) => { if(confirm('删?')) { configs.value.splice(i, 1); saveConfigs() } }
const copy = text => { navigator.clipboard.writeText(text).then(() => { window.$toast('已复制', 'success') }) }

onMounted(() => {
  loadAlbums(); loadWhitelist();
  const saved = localStorage.getItem('picgo_api_configs'); if (saved) configs.value = JSON.parse(saved); else addConfig()
  performUpdateCheck()
  checkTimer = setInterval(() => performUpdateCheck(), 3600000)
})

onUnmounted(() => {
  if (checkTimer) clearInterval(checkTimer)
  if (progressTimer) clearInterval(progressTimer)
})
</script>

<style scoped>
/* 全量样式保留 */
.settings-container { background: var(--card-bg); padding: 25px; border-radius: 12px; box-shadow: var(--shadow-sm); min-height: 500px;}
.header { margin-bottom: 20px; border-bottom: 2px solid var(--border-color); padding-bottom: 10px; }
.desc { color: var(--text-desc); font-size: 14px; margin-bottom: 15px; }
.settings-layout { display: flex; gap: 30px; }
.settings-sidebar { width: 160px; border-right: 1px solid var(--border-color); padding-right: 15px; }
.menu-item { padding: 12px 15px; border-radius: 8px; cursor: pointer; margin-bottom: 5px; color: var(--text-desc); font-weight: 600; transition: 0.2s;}
.menu-item.active { background: var(--accent-color); color: #fff; }
.menu-item:hover:not(.active) { background: var(--bg-color); color: var(--accent-color); }
.settings-content { flex: 1; }
.theme-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; }
.theme-card { border: 2px solid var(--border-color); border-radius: 12px; padding: 15px; text-align: center; cursor: pointer; transition: 0.3s; background: var(--card-bg); font-weight: bold;}
.theme-card.active { border-color: var(--accent-color); box-shadow: 0 0 0 3px rgba(16, 185, 129, 0.2); }
.theme-preview { height: 60px; border-radius: 6px; margin-bottom: 10px; }
.card-box { background: var(--bg-color); border: 1px solid var(--border-color); padding: 15px; border-radius: 8px; }
.about-panel { text-align: center; padding-top: 30px;}
.logo-area { font-size: 36px; font-weight: 900; color: var(--accent-color); letter-spacing: 2px; margin-bottom: 10px;}
.about-list { list-style: none; padding: 0; text-align: left; max-width: 500px; margin: 0 auto; background: var(--bg-color); padding: 25px; border-radius: 12px; border: 1px solid var(--border-color);}
.about-list li { margin-bottom: 15px; font-size: 15px; color: var(--text-main); line-height: 1.6;}
.progress-bar-container { width: 100%; height: 10px; background: rgba(0,0,0,0.05); border-radius: 5px; overflow: hidden; margin: 10px 0; }
.progress-bar-fill { height: 100%; background: #10b981; transition: width 0.3s ease; }
.progress-bar-fill.animated { 
  background: linear-gradient(90deg, #10b981 0%, #34d399 50%, #10b981 100%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite linear;
}
@keyframes shimmer { from { background-position: 200% 0; } to { background-position: 0 0; } }
.update-modal-mask { position: fixed; top:0; left:0; width:100%; height:100%; background: rgba(0,0,0,0.6); display: flex; align-items: center; justify-content: center; z-index: 1000; backdrop-filter: blur(4px); }
.update-modal-card { background: var(--card-bg); width: 450px; border-radius: 16px; padding: 30px; box-shadow: 0 20px 40px rgba(0,0,0,0.3); border: 1px solid var(--border-color); }
.modal-header { text-align: center; margin-bottom: 20px; }
.version-badge { background: #10b981; color: #fff; padding: 4px 12px; border-radius: 20px; font-size: 12px; font-weight: bold; }
.changelog-container { background: var(--bg-color); padding: 15px; border-radius: 8px; font-size: 14px; color: var(--text-desc); max-height: 150px; overflow-y: auto; line-height: 1.5; white-space: pre-wrap; }
.modal-footer { display: flex; flex-direction: column; gap: 10px; }
.btn { border: none; padding: 12px; border-radius: 8px; cursor: pointer; font-weight: 600; transition: 0.3s; width: 100%; }
.primary-btn { background: var(--accent-color); color: #fff; }
.outline-btn { background: transparent; border: 1px solid var(--border-color); color: var(--text-main); }
.ghost-btn { background: transparent; color: var(--text-desc); font-size: 13px; }
.note { font-size: 12px; color: #ff9800; margin-top: 15px; font-style: italic; }
</style>