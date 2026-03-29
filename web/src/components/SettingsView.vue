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
                <code 
                  class="api-link" 
                  @click="copy(getUrl(cfg))" 
                  title="点击即可直接复制链接！"
                  style="cursor: pointer; flex: 1; margin: 0; display: block; background: #282c34; color: #98c379; padding: 10px; border-radius: 6px; font-size: 12px; overflow-x: auto; white-space: nowrap; transition: 0.2s;"
                >
                  {{ getUrl(cfg) }}
                </code>
                <button @click="copy(`![盲盒图](${getUrl(cfg)})`)" class="btn outline-btn small-btn" style="padding: 10px 15px;">📝 复制 MD</button>
              </div>
            </div>
          </div>
        </div>

        <div v-if="activeTab === 'backup'" class="panel">
          <h3>💾 数据库备份与恢复</h3>
          <div class="card-box" style="display: flex; gap: 20px; align-items: center;"><button @click="downloadBackup" class="btn primary-btn">📥 下载 data.db</button><button @click="$refs.restoreInput.click()" class="btn danger-btn">📤 上传恢复数据库</button><input type="file" ref="restoreInput" accept=".db" @change="uploadRestore" hidden /></div>
        </div>
        <div v-if="activeTab === 'about'" class="panel about-panel">
          <div class="logo-area">✨ 拾光图床</div><h3>LightPicker v1.0.0</h3>
          <ul class="about-list">
            <li><strong>👨‍💻 创始人：</strong> @肖肖雨歇</li>
            <li><strong>🎯 理念：</strong> 1C1G 零负担，数据完全私有。</li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
const emit = defineEmits(['change-theme'])

const activeTab = ref('domain')
const currentTheme = ref(localStorage.getItem('app_theme') || 'fresh')
const albums = ref([]); const configs = ref([]); 
const domains = ref(JSON.parse(localStorage.getItem('picgo_domains') || '[]'))

// 🌟 将长长的字符串变成了友好的数组
const whitelistArray = ref([]) 

const host = window.location.origin
const authHeader = () => ({ headers: { Authorization: `Bearer ${localStorage.getItem('picgo_token')}` } })

const selectTheme = (theme) => { currentTheme.value = theme; emit('change-theme', theme) }
const loadAlbums = async () => { try { const {data} = await axios.get('/api/albums', authHeader()); albums.value = data.data } catch(e){} }

const saveDomains = () => { localStorage.setItem('picgo_domains', JSON.stringify(domains.value)) }
const addDomain = () => { domains.value.push(''); saveDomains() }
const removeDomain = (idx) => { domains.value.splice(idx, 1); saveDomains() }

// 🌟 白名单数组的交互与保存逻辑
const loadWhitelist = async () => { 
  try { 
    const {data} = await axios.get('/api/whitelist', authHeader())
    if(data.data) whitelistArray.value = data.data.split(',').filter(item => item.trim() !== '')
  } catch(e){} 
}
const addWhitelist = () => { whitelistArray.value.push('') }
const removeWhitelist = (idx) => { whitelistArray.value.splice(idx, 1) }
const saveWhitelist = async () => { 
  // 发送给后端前，把数组重新拼成逗号分隔的字符串
  const str = whitelistArray.value.filter(item => item.trim() !== '').join(',')
  try { 
    await axios.post('/api/whitelist', { value: str }, authHeader()); 
    window.$toast('防盗链神盾已生效！', 'success') 
  } catch(e){ window.$toast('保存失败', 'error') } 
}

const downloadBackup = () => { window.open('/api/backup?token=' + localStorage.getItem('picgo_token')) }
const uploadRestore = async (e) => {
  const file = e.target.files[0]; if (!file) return;
  if (!confirm('警告：将覆盖当前所有分类配置，确认吗？')) return;
  const fd = new FormData(); fd.append('file', file)
  try { const { data } = await axios.post('/api/restore', fd, authHeader()); window.$toast(data.message, 'success'); setTimeout(() => window.location.reload(), 1500) } catch (err) { window.$toast('恢复失败', 'error') }
}

const getUrl = (cfg) => { 
  const p = new URLSearchParams(); if(cfg.ori!=='all') p.append('ori', cfg.ori); if(cfg.album!=='0') p.append('album_id', cfg.album); 
  const base = (cfg.domain && cfg.domain !== 'default') ? cfg.domain : host;
  return `${base}/api/random${p.toString()?'?'+p.toString():''}` 
}
const saveConfigs = () => localStorage.setItem('picgo_api_configs', JSON.stringify(configs.value))
const addConfig = () => { configs.value.unshift({ id: Date.now(), name: '新配置', album: '0', ori: 'all', domain: 'default' }); saveConfigs() }
const removeConfig = (i) => { if(confirm('删?')) { configs.value.splice(i, 1); saveConfigs() } }

const copy = text => { navigator.clipboard.writeText(text).then(() => { window.$toast('已复制链接', 'success') }) }

onMounted(() => {
  loadAlbums(); loadWhitelist();
  const saved = localStorage.getItem('picgo_api_configs'); if (saved) configs.value = JSON.parse(saved); else addConfig()
})
</script>

<style scoped>
/* 保持高级质感 CSS */
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
.theme-card:hover { border-color: var(--accent-color); }
.theme-card.active { border-color: var(--accent-color); box-shadow: 0 0 0 3px rgba(16, 185, 129, 0.2); }
.theme-preview { height: 60px; border-radius: 6px; margin-bottom: 10px; }
.card-box { background: var(--bg-color); border: 1px solid var(--border-color); padding: 15px; border-radius: 8px; }
.title-input { background: transparent; border: 1px solid transparent; color: var(--text-main); font-weight: bold; width: 60%;}
.title-input:focus { border-bottom-color: var(--accent-color); outline: none;}
.about-panel { text-align: center; padding-top: 30px;}
.logo-area { font-size: 36px; font-weight: 900; color: var(--accent-color); letter-spacing: 2px; margin-bottom: 10px;}
.about-list { list-style: none; padding: 0; text-align: left; max-width: 500px; margin: 0 auto; background: var(--bg-color); padding: 25px; border-radius: 12px; border: 1px solid var(--border-color);}
.about-list li { margin-bottom: 15px; font-size: 15px; color: var(--text-main); line-height: 1.6;}
.empty-state { text-align: center; color: var(--text-desc); padding: 20px; }

/* 代码框悬浮交互 */
.api-link:hover { opacity: 0.8; box-shadow: 0 0 8px rgba(0,0,0,0.1); }
</style>