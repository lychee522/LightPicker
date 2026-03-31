<template>
  <div class="dashboard-panel">
    <div class="header">
      <h2>🚀 上传大厅</h2>
    </div>

    <div class="config-bar">
      <div class="left-group">
        <label class="label-text">📁 存放位置：</label>
        <select v-model="selectedAlbum" class="form-select">
          <option value="0">默认根目录</option>
          <option v-for="album in albums" :key="album.id" :value="album.id">{{ album.name }}</option>
        </select>
        <button @click="createNewAlbum" class="btn outline-btn small-btn">➕ 新建</button>
      </div>
      <button @click="openScanModal" class="btn server-btn small-btn">🖥️ 导入服务器图片</button>
    </div>

    <div class="config-bar">
      <label class="check-box"><input type="checkbox" v-model="useWebp" /> 开启前端 WebP 极致瘦身</label>
      <div v-if="useWebp" class="slider-box">画质: <b>{{ quality }}%</b> <input type="range" v-model="quality" min="10" max="100" /></div>
    </div>

    <div class="drop-zone" @dragover.prevent @drop.prevent="handleDrop">
      <p class="drop-text">📥 将图片拖拽到此区域，或使用下方按钮</p>
      <div class="btn-group">
        <button @click="$refs.fileRef.click()" class="btn primary-btn">📄 网页选择图片</button>
        <button @click="$refs.folderRef.click()" class="btn outline-btn">📁 网页上传文件夹</button>
      </div>
      <input type="file" ref="fileRef" multiple accept="image/*" @change="handleFiles" hidden />
      <input type="file" ref="folderRef" webkitdirectory multiple @change="handleFiles" hidden />
    </div>

    <div class="gallery">
      <div v-for="img in images" :key="img.url" class="img-card">
        <img :src="img.url" class="thumb" />
        <div class="actions">
          <button @click="copy(img.url)">🔗 直链</button>
          <button @click="copy(`![image](${img.url})`)">📝 MD</button>
        </div>
      </div>
    </div>

    <div v-if="showScanModal" class="modal-mask" @click.self="showScanModal = false">
      <div class="modal-box">
        <div class="modal-header">
          <h3 class="scan-title">🖥️ 服务器硬盘大搜捕</h3>
          <button @click="showScanModal = false" class="close-btn">❌</button>
        </div>
        
        <div class="warning-text mb-3">
          <strong>💡 提示：</strong>后台扫描为底层 I/O 操作，将<b>直接导入原图</b>，不经过前端 WebP 压缩。<br/>
          <span v-if="isDockerEnv" style="color: #e53e3e; font-weight: bold; margin-top: 5px; display: inline-block;">
            ⚠️ 检测到您当前运行在 Docker 环境！请确保已通过 -v 参数映射了宿主机目录，否则您只能浏览容器内部的隔离空间！
          </span>
        </div>

        <div class="scan-controls mb-3">
          <input v-model="serverScanPath" @keyup.enter="loadServerDirs(serverScanPath)" type="text" placeholder="输入绝对路径并回车 (留空看盘符)" class="form-input" />
          <button @click="loadServerDirs(serverScanPath)" class="btn outline-btn small-btn">刷新</button>
        </div>

        <div class="dir-list-box mb-3">
          <div v-if="showUpButton" @click="goUpDir" class="dir-item up-dir">
            ↩️ 返回上一级
          </div>
          <div v-if="serverDirs.length === 0" class="empty-dir">当前目录为空或未找到驱动器</div>
          <div v-for="dir in serverDirs" :key="dir.path" @click="loadServerDirs(dir.path)" class="dir-item">
            📁 {{ dir.name }}
          </div>
        </div>

        <div class="modal-footer">
          <select v-model="selectedAlbum" class="form-select strategy-select" style="max-width: 120px;" title="选择导入的目标相册">
            <option value="0">📂 默认根目录</option>
            <option v-for="album in albums" :key="album.id" :value="album.id">📁 {{ album.name }}</option>
          </select>

          <select v-model="scanStrategy" class="form-select strategy-select" style="max-width: 120px;" title="选择导入策略">
            <option value="copy">复制 (安全)</option>
            <option value="move">移动 (省空间)</option>
            <option value="link">软链接 (极客)</option>
          </select>

          <button @click="triggerServerScan" class="btn scan-btn" :disabled="isScanning">
            {{ isScanning ? '🚀 狂奔扫描中...' : '🚀 立即扫描当前目录' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
// @author tg账号的肖肖雨歇
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'

const useWebp = ref(true)
const quality = ref(75)
const images = ref([])
const albums = ref([])
const selectedAlbum = ref('0')

const showScanModal = ref(false)
const serverScanPath = ref('') 
const scanStrategy = ref('copy')
const serverDirs = ref([])
const isScanning = ref(false)
const isDockerEnv = ref(false) // 🌟 记录是否为 Docker 环境

const authHeader = () => ({ headers: { Authorization: `Bearer ${localStorage.getItem('picgo_token')}` } })

// 探测服务器运行环境 (判断是否为 Docker)
const checkServerEnv = async () => {
  try {
    const { data } = await axios.get('/api/env', authHeader())
    isDockerEnv.value = data.is_docker
  } catch (err) {
    console.log("环境探测失败，忽略提示")
  }
}

const showUpButton = computed(() => {
  const p = serverScanPath.value
  if (p === '' || p === '/') return false
  return true 
})

const loadAlbums = async () => {
  try {
    const { data } = await axios.get('/api/albums', authHeader())
    albums.value = data.data
  } catch (err) {}
}

const createNewAlbum = async () => {
  const name = prompt('给新文件夹起个名字：')
  if (!name) return
  try {
    await axios.post('/api/albums', { name }, authHeader())
    loadAlbums()
  } catch (err) { alert('创建失败') }
}

const openScanModal = () => {
  showScanModal.value = true
  if (serverScanPath.value === '/') {
    serverScanPath.value = ''
  }
  loadServerDirs(serverScanPath.value)
}

const loadServerDirs = async (path) => {
  if (path === undefined || path === null) return 
  serverScanPath.value = path
  try {
    const { data } = await axios.get(`/api/fs/list?path=${encodeURIComponent(path)}`, authHeader())
    if (data.code === 200) {
      serverDirs.value = data.data || []
    } else {
      window.$toast(data.msg || '读取目录失败', 'error')
      serverDirs.value = []
    }
  } catch (err) {
    const errMsg = err.response?.data?.msg || '读取目录失败，请检查路径或权限'
    window.$toast(errMsg, 'error')
    serverDirs.value = [] 
  }
}

const goUpDir = () => {
  let current = serverScanPath.value
  let pathStr = current.replace(/\\/g, '/') 
  if (pathStr.endsWith('/')) pathStr = pathStr.slice(0, -1) 
  
  let parts = pathStr.split('/')
  parts.pop() 
  let newPath = parts.join('/')

  if (!newPath) {
    newPath = '' 
  } else if (/^[a-zA-Z]:$/.test(newPath)) {
    newPath += '/' 
  }

  loadServerDirs(newPath)
}

const triggerServerScan = async () => {
  if (!serverScanPath.value) {
    alert('❌ 老哥，扫描路径不能为空！你得点进一个具体的文件夹或盘符里啊！')
    return
  }
  
  isScanning.value = true
  
  const payload = {
    sourcePath: serverScanPath.value,
    strategy: scanStrategy.value,
    // 🌟 核心修复：强制转为字符串，治好 Go 语言的强迫症
    album: String(selectedAlbum.value) 
  }

  try {
    const { data } = await axios.post('/api/fs/import', payload, authHeader())
    
    alert(data.msg || '扫描并导入完成！') 
    
    if (data.data && data.data.urls && data.data.urls.length > 0) {
      data.data.urls.forEach(url => {
         const prefix = url.startsWith('/') ? '' : '/'
         images.value.unshift({ url: window.location.origin + prefix + url })
      })
    }

    showScanModal.value = false // 导入成功后自动关闭弹窗
  } catch (err) {
    alert(err.response?.data?.msg || '❌ 扫描失败，请检查路径是否存在或权限是否足够！')
  } finally {
    isScanning.value = false
  }
}

const handleDrop = (e) => process(e.dataTransfer.files)
const handleFiles = (e) => process(e.target.files)

const process = async (files) => {
  for (let i = 0; i < files.length; i++) {
    let file = files[i]
    if (!file.type.startsWith('image/')) continue

    let width = 0, height = 0
    await new Promise(resolve => {
      const url = URL.createObjectURL(file)
      const img = new Image()
      img.onload = () => { width = img.width; height = img.height; URL.revokeObjectURL(url); resolve() }
      img.src = url
    })

    if (useWebp.value && file.type !== 'image/webp') file = await compressToWebp(file, width, height)
    await upload(file, width, height)
  }
}

const compressToWebp = (file, width, height) => new Promise(resolve => {
  const reader = new FileReader()
  reader.onload = e => {
    const img = new Image()
    img.onload = () => {
      const cvs = document.createElement('canvas')
      cvs.width = width; cvs.height = height
      cvs.getContext('2d').drawImage(img, 0, 0)
      cvs.toBlob(blob => {
        const newName = file.name.replace(/\.[^/.]+$/, "") + ".webp"
        resolve(new File([blob], newName, { type: 'image/webp' }))
      }, 'image/webp', quality.value / 100)
    }
    img.src = e.target.result
  }
  reader.readAsDataURL(file)
})

const upload = async (file, width, height) => {
  const fd = new FormData()
  fd.append('file', file)
  fd.append('width', width)
  fd.append('height', height)
  // 🌟 这里也加上 String 保护一下，虽然正常网页传 FormData 都是字符串，但保险起见
  fd.append('album_id', String(selectedAlbum.value))

  try {
    const { data } = await axios.post('/api/upload', fd, authHeader())
    images.value.unshift({ url: window.location.origin + data.url })
  } catch (err) { alert('上传失败') }
}

const copy = text => navigator.clipboard.writeText(text).then(() => alert('🎉 复制成功！'))

onMounted(() => {
  loadAlbums()
  checkServerEnv() // 🌟 挂载时探测环境
})
</script>

<style scoped>
.dashboard-panel { width: 100%; display: flex; flex-direction: column; gap: 15px; }
.header { display: flex; justify-content: space-between; align-items: center; border-bottom: 2px solid #f0f0f0; padding-bottom: 10px;}
.config-bar { background: #f0f7ff; padding: 12px 15px; border-radius: 8px; display: flex; justify-content: space-between; align-items: center; gap: 15px; border: 1px solid #dcebfa;}
.left-group { display: flex; align-items: center; gap: 10px; flex: 1; }
.form-select { flex: 1; padding: 6px; border-radius: 4px; border: 1px solid #ccc; font-weight: bold; color: #333;}
.label-text { font-weight: bold; color: #333; }
.slider-box { display: flex; align-items: center; gap: 10px; font-size: 14px;}

.server-btn { background: #805ad5; color: white; border: none; }
.server-btn:hover { background: #6b46c1; }
.small-btn { padding: 6px 12px; font-size: 13px; }

.drop-zone { border: 2px dashed #4CAF50; border-radius: 12px; padding: 30px 10px; text-align: center; background: #fafafa; transition: 0.3s;}
.drop-zone:hover { background: #e8f5e9; }
.drop-text { color: #666; margin-bottom: 20px; font-weight: bold;}
.btn-group { display: flex; justify-content: center; gap: 10px; }
.btn { padding: 8px 16px; border: none; border-radius: 6px; cursor: pointer; transition: 0.2s; font-weight: bold;}
.primary-btn { background: #4CAF50; color: white; }
.primary-btn:hover { background: #45a049; }
.outline-btn { background: transparent; border: 2px solid #4CAF50; color: #4CAF50; }
.outline-btn:hover { background: #4CAF50; color: white; }

.gallery { display: flex; flex-direction: column; gap: 10px; max-height: 400px; overflow-y: auto;}
.img-card { display: flex; gap: 15px; background: #f9f9f9; padding: 10px; border-radius: 8px; align-items: center; border: 1px solid #eee;}
.thumb { width: 80px; height: 80px; object-fit: cover; border-radius: 6px; border: 1px solid #ddd;}
.actions { display: flex; flex-wrap: wrap; gap: 8px; }
.actions button { padding: 4px 8px; font-size: 12px; border: 1px solid #ccc; background: white; border-radius: 4px; cursor: pointer;}
.actions button:hover { background: #f0f0f0; border-color: #999;}

.modal-mask { position: fixed; top: 0; left: 0; width: 100vw; height: 100vh; background: rgba(0,0,0,0.5); display: flex; justify-content: center; align-items: center; z-index: 100; backdrop-filter: blur(2px);}
.modal-box { background: white; width: 500px; max-width: 90vw; border-radius: 12px; padding: 20px; box-shadow: 0 10px 25px rgba(0,0,0,0.2); display: flex; flex-direction: column;}
.modal-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 15px; border-bottom: 1px solid #eee; padding-bottom: 10px;}
.scan-title { margin: 0; color: #dd6b20; font-size: 18px; font-weight: bold;}
.close-btn { background: transparent; border: none; font-size: 16px; cursor: pointer; transition: 0.2s; }
.close-btn:hover { transform: scale(1.2); }
.warning-text { font-size: 12px; color: #b7791f; background: #fffff0; padding: 8px; border-radius: 6px; border: 1px dashed #f6e05e; line-height: 1.4; }
.scan-controls { display: flex; gap: 10px; align-items: center; }
.form-input { flex: 1; padding: 8px 12px; border-radius: 6px; border: 1px solid #cbd5e0; outline: none; transition: 0.2s; }
.form-input:focus { border-color: #ed8936; box-shadow: 0 0 0 2px rgba(237, 137, 54, 0.2); }

.dir-list-box { border: 1px solid #e2e8f0; border-radius: 6px; height: 200px; overflow-y: auto; background: #f8fafc; padding: 5px;}
.dir-item { padding: 8px 10px; cursor: pointer; border-radius: 4px; color: #4a5568; font-size: 14px; transition: 0.1s;}
.dir-item:hover { background: #e2e8f0; color: #2b6cb0; font-weight: bold;}
.up-dir { color: #805ad5; font-weight: bold; border-bottom: 1px dashed #cbd5e0; margin-bottom: 5px;}
.empty-dir { text-align: center; color: #a0aec0; margin-top: 20px; font-size: 13px; }

.modal-footer { display: flex; gap: 10px; margin-top: 10px;}
.strategy-select { max-width: 140px; }
.scan-btn { flex: 1; background: #ed8936; color: white; border: none; }
.scan-btn:hover { background: #dd6b20; }
.scan-btn:disabled { background: #a0aec0; cursor: not-allowed; }
</style>