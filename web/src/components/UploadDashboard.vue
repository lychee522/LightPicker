<template>
  <div class="dashboard-panel">
    <div class="header">
      <h2>🚀 上传大厅</h2>
      </div>

    <div class="config-bar">
      <label class="label-text">📁 存放位置：</label>
      <select v-model="selectedAlbum" class="form-select">
        <option value="0">默认根目录</option>
        <option v-for="album in albums" :key="album.id" :value="album.id">{{ album.name }}</option>
      </select>
      <button @click="createNewAlbum" class="btn outline-btn small-btn">➕ 新建文件夹</button>
    </div>

    <div class="config-bar">
      <label class="check-box"><input type="checkbox" v-model="useWebp" /> 开启前端 WebP 极致瘦身</label>
      <div v-if="useWebp" class="slider-box">画质: <b>{{ quality }}%</b> <input type="range" v-model="quality" min="10" max="100" /></div>
    </div>

    <div class="drop-zone" @dragover.prevent @drop.prevent="handleDrop">
      <p class="drop-text">📥 将图片拖拽到此区域，或使用下方按钮</p>
      <div class="btn-group">
        <button @click="$refs.fileRef.click()" class="btn primary-btn">📄 选择图片</button>
        <button @click="$refs.folderRef.click()" class="btn outline-btn">📁 扫描本地文件夹</button>
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
  </div>
</template>

<script setup>
// @author 肖肖雨歇
import { ref, onMounted } from 'vue'
import axios from 'axios'

const useWebp = ref(true)
const quality = ref(75)
const images = ref([])
const albums = ref([])
const selectedAlbum = ref('0')

const authHeader = () => ({ headers: { Authorization: `Bearer ${localStorage.getItem('picgo_token')}` } })

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
  fd.append('album_id', selectedAlbum.value)

  try {
    const { data } = await axios.post('/api/upload', fd, authHeader())
    images.value.unshift({ url: window.location.origin + data.url })
  } catch (err) { alert('上传失败') }
}

const copy = text => navigator.clipboard.writeText(text).then(() => alert('🎉 复制成功！'))

onMounted(() => loadAlbums())
</script>

<style scoped>
.dashboard-panel { width: 100%; display: flex; flex-direction: column; gap: 15px; }
.header { display: flex; justify-content: space-between; align-items: center; border-bottom: 2px solid #f0f0f0; padding-bottom: 10px;}
.config-bar { background: #f0f7ff; padding: 12px 15px; border-radius: 8px; display: flex; justify-content: space-between; align-items: center; gap: 15px; border: 1px solid #dcebfa;}
.form-select { flex: 1; padding: 6px; border-radius: 4px; border: 1px solid #ccc; font-weight: bold; color: #333;}
.label-text { font-weight: bold; color: #333; }
.slider-box { display: flex; align-items: center; gap: 10px; font-size: 14px;}
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
</style>