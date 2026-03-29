<template>
  <div class="gallery-container" @click="closeMenu">
    <div class="header">
      <h2>🖼️ 拾光图库</h2>
      <div class="toolbar">
        <select v-model="currentDomain" class="view-select" style="border-color: var(--accent-color); color: var(--accent-color);">
          <option value="default">🌐 默认线路</option>
          <option v-for="d in domains" :key="d" :value="d">🚀 {{ d }}</option>
        </select>
        
        <select v-model="browseMode" @change="switchMode" class="view-select"><option value="all">🖼️ 全部图</option><option value="folder">🗂️ 文件夹</option></select>
        <select v-model="viewMode" class="view-select"><option value="large">🔳 大图</option><option value="small">🔲 小图</option><option value="list">📄 列表</option></select>
        <button @click="refreshData" class="btn primary-btn refresh-btn">🔄 刷新</button>
      </div>
    </div>

    <div v-if="browseMode === 'folder' && currentAlbumId !== '0'" class="breadcrumb">
      <button @click="goBack" class="btn outline-btn small-btn">🔙 返回上级</button>
      <span class="current-path">当前：{{ currentAlbumName }}</span>
    </div>

    <div v-if="browseMode === 'folder' && currentAlbumId === '0'" class="folder-grid">
      <div v-for="album in albums" :key="album.id" class="folder-card" @click="enterFolder(album)" @contextmenu.prevent="openMenu($event, 'folder', album)">
        <div class="folder-icon">📁</div><div class="folder-name">{{ album.name }}</div>
      </div>
    </div>

    <div v-else :class="['gallery-wrapper', viewMode]">
      <div v-for="img in images" :key="img.id" class="image-card" @contextmenu.prevent="openMenu($event, 'image', img)">
        <div class="img-wrapper" @click="openPreview(img.url)">
          <img :src="`/${img.url}`" class="preview-img" loading="lazy" />
          <div class="img-overlay">🔍 右键管理</div>
        </div>
        <div class="card-body">
          <div class="card-info"><span class="file-size">{{ (img.size / 1024).toFixed(1) }} KB</span></div>
          <div class="actions">
            <button @click="copy(`${getImgBase()}/${img.url}`)">🔗 直链</button>
            <button @click="copy(`![image](${getImgBase()}/${img.url})`)">📝 MD</button>
            <button @click="copy(`<img src='${getImgBase()}/${img.url}' />`)">🌐 HTML</button>
            <button @click="copy(`[img]${getImgBase()}/${img.url}[/img]`)">🏷️ BBCode</button>
          </div>
        </div>
      </div>
    </div>

    <div v-show="ctxMenu.show" :style="{ top: ctxMenu.y + 'px', left: ctxMenu.x + 'px' }" class="context-menu" @click.stop>
      <template v-if="ctxMenu.type === 'folder'">
        <div class="context-menu-item" @click="handleRenameFolder">✏️ 重命名</div>
        <div class="context-menu-item danger" @click="handleDeleteFolder">🗑️ 删除(保图)</div>
      </template>
      <template v-if="ctxMenu.type === 'image'">
        <div class="context-menu-item" @click="copy(`${getImgBase()}/${ctxMenu.target.url}`)">🔗 复制外链</div>
        <div class="context-menu-item" style="cursor: default;">
          🚚 移至: 
          <select class="move-select" v-model="moveTargetId" @change="confirmMoveImg">
            <option value="0">默认根目录</option><option v-for="a in albums" :key="a.id" :value="a.id.toString()">📁 {{ a.name }}</option>
          </select>
        </div>
        <div class="context-menu-item danger" @click="deleteImg(ctxMenu.target.id)">🗑️ 彻底删除</div>
      </template>
    </div>

    <div v-if="previewUrl" class="lightbox" @click="previewUrl = null">
      <button class="close-btn">❌ 关闭</button>
      <img :src="previewUrl" class="lightbox-img" />
    </div>
  </div>
</template>

<script setup>
// @author 肖肖雨歇 - 加入全局外网域名切换引擎
import { ref, onMounted } from 'vue'
import axios from 'axios'

const host = window.location.origin
const domains = ref(JSON.parse(localStorage.getItem('picgo_domains') || '[]'))
const currentDomain = ref('default') // 🌟 当前选中的域名线路

// 动态计算图片的绝对基础路径
const getImgBase = () => {
  return currentDomain.value === 'default' ? host : currentDomain.value
}

const browseMode = ref('all'); const currentAlbumId = ref('0'); const currentAlbumName = ref('')
const viewMode = ref('large'); const albums = ref([]); const images = ref([]); const previewUrl = ref(null)
const ctxMenu = ref({ show: false, x: 0, y: 0, type: '', target: null }); const moveTargetId = ref('0')

const authHeader = () => ({ headers: { Authorization: `Bearer ${localStorage.getItem('picgo_token')}` } })

const loadAlbums = async () => { try { const { data } = await axios.get('/api/albums', authHeader()); albums.value = data.data } catch (err) {} }
const loadImages = async () => {
  let url = '/api/images'; if (browseMode.value === 'folder' && currentAlbumId.value !== '0') url += `?album_id=${currentAlbumId.value}`
  try { const { data } = await axios.get(url, authHeader()); images.value = data.data } catch (err) {}
}

const switchMode = () => { currentAlbumId.value = '0'; if (browseMode.value === 'all') loadImages() }
const enterFolder = (a) => { currentAlbumId.value = a.id.toString(); currentAlbumName.value = a.name; loadImages() }
const goBack = () => { currentAlbumId.value = '0'; images.value = [] }
const refreshData = () => { loadAlbums(); if (browseMode.value === 'all' || currentAlbumId.value !== '0') loadImages() }

const openMenu = (e, type, target) => { e.preventDefault(); moveTargetId.value = type === 'image' ? target.album_id.toString() : '0'; ctxMenu.value = { show: true, x: e.clientX, y: e.clientY, type, target } }
const closeMenu = () => { ctxMenu.value.show = false }

const handleRenameFolder = async () => {
  const newName = prompt('新名字：', ctxMenu.value.target.name)
  if (!newName || newName === ctxMenu.value.target.name) return closeMenu()
  try { await axios.put(`/api/albums/${ctxMenu.value.target.id}`, { name: newName }, authHeader()); window.$toast('重命名成功'); refreshData() } catch(e){}
  closeMenu()
}
const handleDeleteFolder = async () => {
  if (!confirm('确定删除吗？图片会被移到根目录。')) return closeMenu()
  try { await axios.delete(`/api/albums/${ctxMenu.value.target.id}`, authHeader()); window.$toast('已删除'); refreshData() } catch(e){}
  closeMenu()
}
const confirmMoveImg = async () => {
  try { await axios.put(`/api/images/${ctxMenu.value.target.id}/move`, { album_id: moveTargetId.value }, authHeader()); window.$toast('移动成功'); refreshData() } catch(e){}
  closeMenu()
}

const openPreview = (url) => { previewUrl.value = `${host}/${url}` } // 预览始终用本地线路最快
const deleteImg = async (id) => {
  if (!confirm('确定彻底粉碎？')) return closeMenu()
  try { await axios.delete(`/api/images/${id}`, authHeader()); window.$toast('已删除'); refreshData() } catch(e){}
  closeMenu()
}
const copy = text => { navigator.clipboard.writeText(text).then(() => { window.$toast('已复制外链'); closeMenu() }) }

onMounted(() => { refreshData() })
</script>

<style scoped>
/* 精简复用 CSS，包含高级右键菜单 */
.gallery-container { background: var(--card-bg); padding: 25px; border-radius: 12px; box-shadow: var(--shadow-sm); min-height: 100%;}
.header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; border-bottom: 2px solid var(--border-color); padding-bottom: 10px; }
.toolbar { display: flex; gap: 10px; }
.view-select { padding: 8px; border-radius: 6px; border: 1px solid var(--border-color); font-weight: bold; background: var(--bg-color); color: var(--text-main);}
.breadcrumb { padding: 10px; background: var(--bg-color); border-radius: 8px; border-left: 4px solid var(--accent-color); margin-bottom: 20px;}
.current-path { font-weight: bold; margin-left: 10px;}
.folder-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(120px, 1fr)); gap: 20px; }
.folder-card { background: var(--card-bg); border: 1px solid var(--border-color); border-radius: 12px; padding: 20px 10px; text-align: center; cursor: pointer; transition: 0.2s;}
.folder-card:hover { transform: translateY(-3px); border-color: var(--accent-color); }
.folder-icon { font-size: 40px; margin-bottom: 10px; }
.folder-name { font-weight: bold; font-size: 14px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;}
.gallery-wrapper { display: grid; gap: 20px; }
.image-card { border: 1px solid var(--border-color); border-radius: 8px; overflow: hidden; background: var(--card-bg); transition: 0.2s; display: flex; flex-direction: column;}
.image-card:hover { border-color: var(--accent-color); }
.img-wrapper { position: relative; width: 100%; cursor: zoom-in; background: var(--bg-color); display: flex; justify-content: center; align-items: center; }
.preview-img { width: 100%; height: 100%; object-fit: contain; }
.img-overlay { position: absolute; inset: 0; background: rgba(0,0,0,0.5); display: flex; justify-content: center; align-items: center; opacity: 0; transition: 0.3s; color: white; font-weight: bold;}
.img-wrapper:hover .img-overlay { opacity: 1; }
.card-body { display: flex; flex-direction: column; flex: 1; }
.card-info { padding: 10px; border-bottom: 1px dashed var(--border-color); background: var(--bg-color); }
.actions { display: flex; flex-wrap: wrap; justify-content: center; gap: 5px; padding: 10px; background: var(--bg-color); }
.actions button { padding: 4px; font-size: 12px; border: 1px solid var(--border-color); background: var(--card-bg); color: var(--text-main); border-radius: 4px; cursor: pointer; flex: 1; min-width: 45%;}
.small { grid-template-columns: repeat(auto-fill, minmax(160px, 1fr)); } .small .img-wrapper { height: 120px; }
.large { grid-template-columns: repeat(auto-fill, minmax(260px, 1fr)); } .large .img-wrapper { height: 200px; }
.list { display: flex; flex-direction: column; gap: 10px; }
.list .image-card { flex-direction: row; align-items: center; padding: 10px; gap: 20px; }
.list .img-wrapper { width: 60px; height: 60px; }
.list .card-body { flex-direction: row; justify-content: space-between; align-items: center; padding: 0; background: transparent; flex: 1;}
.list .actions { padding: 0; background: transparent; gap: 10px; flex-wrap: nowrap; justify-content: flex-end;}
.list .actions button { min-width: 50px; flex: none;}
.context-menu { position: fixed; background: var(--card-bg); border: 1px solid var(--border-color); box-shadow: var(--shadow-md); border-radius: 8px; padding: 5px; min-width: 180px; z-index: 10000; display: flex; flex-direction: column; gap: 2px;}
.context-menu-item { padding: 10px; font-size: 14px; font-weight: 600; color: var(--text-main); cursor: pointer; border-radius: 6px; transition: 0.2s; display: flex; justify-content: space-between; align-items: center;}
.context-menu-item:hover { background: var(--bg-color); color: var(--accent-color); }
.context-menu-item.danger:hover { background: #fef2f2; color: #ef4444; }
.move-select { padding: 4px; border-radius: 4px; border: 1px solid var(--border-color); font-weight: bold; background: var(--card-bg); color: var(--text-main); max-width: 100px;}
.lightbox { position: fixed; top: 0; left: 0; width: 100vw; height: 100vh; background: rgba(0, 0, 0, 0.85); display: flex; justify-content: center; align-items: center; z-index: 999999; cursor: zoom-out; }
.lightbox-img { max-width: 90%; max-height: 90%; border-radius: 8px; box-shadow: 0 0 20px rgba(0, 0, 0, 0.5); cursor: default; pointer-events: none; }
.close-btn { position: absolute; top: 20px; right: 30px; background: #ef4444; color: white; border: none; padding: 10px 20px; border-radius: 6px; font-weight: bold; cursor: pointer; z-index: 1000000; }
</style>