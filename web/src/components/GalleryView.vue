<template>
  <div class="gallery-panel" @click="closeContextMenu">
    <div class="toolbar">
      <div class="left-tools">
        <h2 class="title">🖼️ 拾光图库</h2>
        
        <button 
          @click="toggleMultiSelectMode" 
          :class="['btn small-btn', isMultiSelectMode ? 'danger-btn' : 'outline-btn']"
        >
          {{ isMultiSelectMode ? '❌ 退出多选' : '☑️ 开启多选' }}
        </button>
        
        <div v-if="isMultiSelectMode" class="bulk-action-bar fade-in">
          <span class="selected-text">已选 <b>{{ selectedIds.length }}</b> 项</span>
          <button @click="selectAllOnPage" class="btn small-btn outline-btn">☑️ 全选本页</button>
          <button @click="clearSelection" class="btn small-btn outline-btn">取消</button>
          
          <div class="divider"></div>
          
          <select v-model="bulkTargetAlbum" class="form-select small-select">
            <option value="" disabled>选择目标文件夹...</option>
            <option value="0">默认相册</option>
            <option v-for="album in albums" :key="album.id" :value="album.id">{{ album.name }}</option>
          </select>
          <button @click="bulkMove" class="btn small-btn primary-btn" :disabled="isProcessing || !bulkTargetAlbum || selectedIds.length === 0">
            📁 批量移动
          </button>
          
          <div class="divider"></div>

          <button @click="bulkDelete" class="btn small-btn danger-btn" :disabled="isProcessing || selectedIds.length === 0">
            {{ isProcessing ? '清理中...' : '🗑️ 彻底删除' }}
          </button>
        </div>
      </div>

      <div class="right-tools">
        <select v-model="selectedAlbum" @change="resetAndLoad" class="form-select">
          <option value="0">默认相册</option>
          <option v-for="album in albums" :key="album.id" :value="album.id">{{ album.name }}</option>
        </select>
        
        <select v-model="viewMode" @change="handleViewModeChange" class="form-select">
          <option value="large">🔲 大图模式 (24/页)</option>
          <option value="small">🔳 小图模式 (48/页)</option>
          <option value="list">📄 列表模式 (100/页)</option>
        </select>
        
        <button @click="resetAndLoad" class="btn primary-btn small-btn">🔄 刷新</button>
      </div>
    </div>

    <div v-if="images.length === 0" class="empty-state">
      图库空空如也，快去上传大厅传几张吧！
    </div>

    <div v-else :class="['image-grid', viewMode]">
      <div 
        v-for="img in images" 
        :key="img.id" 
        :class="['img-card', { 'is-selected': selectedIds.includes(img.id), 'multi-mode': isMultiSelectMode }]"
        @click="isMultiSelectMode ? toggleSelect(img.id) : openPreview(img.url)"
        @contextmenu.prevent="handleContextMenu($event, img)" 
      >
        <div v-if="isMultiSelectMode" class="checkbox-wrap">
          <input type="checkbox" :value="img.id" v-model="selectedIds" @click.stop />
        </div>

        <div class="img-wrapper">
          <img :src="img.url" class="thumb" loading="lazy" />
        </div>
        
        <div class="info-bar">
          <span class="filename-text" v-if="viewMode === 'list'">📄 {{ img.filename }}</span>
          <span class="size-text">💾 {{ formatSize(img.size) }}</span>
        </div>
        
        <div class="actions" @click.stop>
          <button @click="copy(img.url)">🔗 直链</button>
          <button @click="copy(`![image](${img.url})`)">📝 MD</button>
          <button @click="copy(`<img src=\x22${img.url}\x22 alt=\x22image\x22 />`)">🌐 HTML</button>
          <button @click="copy(`[img]${img.url}[/img]`)">💬 BBCode</button>
          <button @click="copy(`[url=${img.url}][img]${img.url}[/img][/url]`)">🚀 论坛代码</button>
        </div>
      </div>
    </div>

    <div class="pagination-bar" v-if="totalPages > 1">
      <button @click="goToPage(1)" :disabled="page === 1" class="page-btn">首页</button>
      <button @click="goToPage(page - 1)" :disabled="page === 1" class="page-btn">上一页</button>
      
      <div class="page-numbers">
        <button 
          v-for="p in pageNumbers" 
          :key="p" 
          @click="goToPage(p)"
          :class="['num-btn', { 'active': p === page }]"
        >
          {{ p }}
        </button>
      </div>

      <button @click="goToPage(page + 1)" :disabled="page === totalPages" class="page-btn">下一页</button>
      <button @click="goToPage(totalPages)" :disabled="page === totalPages" class="page-btn">尾页</button>
      
      <div class="jump-box">
        共 {{ totalPages }} 页，跳至
        <input type="number" v-model.number="jumpPageNum" @keyup.enter="jumpToPage" min="1" :max="totalPages" class="jump-input" />
        页
        <button @click="jumpToPage" class="btn primary-btn small-btn jump-btn">GO</button>
      </div>
    </div>

    <div 
      v-if="contextMenu.show" 
      class="custom-context-menu" 
      :style="{ top: contextMenu.y + 'px', left: contextMenu.x + 'px' }" 
      @click.stop
    >
      <div class="menu-header">操作图片</div>
      
      <div class="menu-item move-item">
        📁 移动至: 
        <select v-model="contextMenuTargetAlbum" class="context-select" @change="moveSingle(contextMenu.img, contextMenuTargetAlbum)">
          <option value="" disabled>选择文件夹</option>
          <option value="0">默认相册</option>
          <option v-for="album in albums" :key="album.id" :value="album.id">{{ album.name }}</option>
        </select>
      </div>
      
      <div class="menu-item" @click="renameSingle(contextMenu.img)">📝 重命名</div>
      <div class="menu-separator"></div>
      <div class="menu-item danger" @click="deleteSingle(contextMenu.img)">🗑️ 删除此图</div>
    </div>

    <div v-if="previewImageUrl" class="preview-mask" @click="closePreview">
      <img :src="previewImageUrl" class="preview-img" @click.stop />
      <button class="close-preview-btn" @click="closePreview" title="关闭预览">❌</button>
    </div>

  </div>
</template>

<script setup>
// @author tg账号的肖肖雨歇
import { ref, onMounted, computed, onUnmounted } from 'vue'
import axios from 'axios'

const images = ref([])
const albums = ref([])
const selectedAlbum = ref('0')
const viewMode = ref('large')

// 多选与批量操作
const isMultiSelectMode = ref(false)
const selectedIds = ref([])
const isProcessing = ref(false)
const bulkTargetAlbum = ref('')

// 分页相关
const page = ref(1)
const pageSize = ref(24)
const totalPages = ref(1)
const jumpPageNum = ref(1)

// 右键菜单状态
const contextMenu = ref({ show: false, x: 0, y: 0, img: null })
const contextMenuTargetAlbum = ref('')

// 全屏预览状态
const previewImageUrl = ref(null)

const authHeader = () => ({ headers: { Authorization: `Bearer ${localStorage.getItem('picgo_token')}` } })

// ---------------- 全屏预览控制 ----------------
const openPreview = (url) => {
  previewImageUrl.value = url
  document.body.style.overflow = 'hidden'
}

const closePreview = () => {
  previewImageUrl.value = null
  document.body.style.overflow = ''
}

// ---------------- 基础图库与相册加载 ----------------
const loadAlbums = async () => {
  try {
    const { data } = await axios.get('/api/albums', authHeader())
    albums.value = data.data || []
  } catch (err) {}
}

const loadImages = async () => {
  try {
    const params = new URLSearchParams()
    params.append('album_id', selectedAlbum.value)
    params.append('page', page.value)
    params.append('size', pageSize.value)

    const { data } = await axios.get(`/api/images?${params.toString()}`, authHeader())
    
    totalPages.value = Math.ceil((data.total || 0) / pageSize.value) || 1
    
    let fetchedImages = []
    if (data.data && Array.isArray(data.data)) {
      fetchedImages = data.data
    }

    images.value = fetchedImages.map(img => {
      if (!img.url.startsWith('http')) {
        const prefix = img.url.startsWith('/') ? '' : '/'
        img.url = window.location.origin + prefix + img.url
      }
      return img
    })
    
    if (isMultiSelectMode.value) {
      selectedIds.value = []
    }
  } catch (err) {
    alert('加载图库失败，请检查网络或后端状态')
  }
}

// ---------------- 视图与分页控制 ----------------
const handleViewModeChange = () => {
  if (viewMode.value === 'large') pageSize.value = 24
  else if (viewMode.value === 'small') pageSize.value = 48
  else if (viewMode.value === 'list') pageSize.value = 100
  resetAndLoad()
}

const resetAndLoad = () => {
  page.value = 1
  loadImages()
}

const goToPage = (p) => {
  if (p >= 1 && p <= totalPages.value) {
    page.value = p
    loadImages()
  }
}

const jumpToPage = () => {
  let target = parseInt(jumpPageNum.value)
  if (isNaN(target) || target < 1) target = 1
  if (target > totalPages.value) target = totalPages.value
  jumpPageNum.value = target
  goToPage(target)
}

const pageNumbers = computed(() => {
  let start = Math.max(1, page.value - 2)
  let end = Math.min(totalPages.value, page.value + 2)
  if (end - start < 4) {
    if (start === 1) end = Math.min(totalPages.value, start + 4)
    else if (end === totalPages.value) start = Math.max(1, end - 4)
  }
  const arr = []
  for (let i = start; i <= end; i++) arr.push(i)
  return arr
})

// ---------------- 批量多选逻辑 ----------------
const toggleMultiSelectMode = () => {
  isMultiSelectMode.value = !isMultiSelectMode.value
  if (!isMultiSelectMode.value) {
    selectedIds.value = []
  }
}

const toggleSelect = (id) => {
  const index = selectedIds.value.indexOf(id)
  if (index === -1) selectedIds.value.push(id)
  else selectedIds.value.splice(index, 1)
}

const selectAllOnPage = () => {
  selectedIds.value = images.value.map(img => img.id)
}

const clearSelection = () => {
  selectedIds.value = []
}

const bulkDelete = async () => {
  if (selectedIds.value.length === 0) return
  if (!confirm(`⚠️ 确定要永远删除这 ${selectedIds.length} 张图片吗？`)) return

  isProcessing.value = true
  let successCount = 0, failCount = 0

  const deletePromises = selectedIds.value.map(id => 
    axios.delete(`/api/images/${id}`, authHeader())
      .then(() => successCount++)
      .catch(() => failCount++)
  )

  await Promise.allSettled(deletePromises)
  isProcessing.value = false
  alert(`批量清理完成！\n✅ 成功: ${successCount}\n❌ 失败: ${failCount}`)
  selectedIds.value = []
  loadImages()
}

const bulkMove = async () => {
  if (selectedIds.value.length === 0 || !bulkTargetAlbum.value) return
  
  isProcessing.value = true
  let successCount = 0, failCount = 0

  const movePromises = selectedIds.value.map(id => 
    axios.put(`/api/images/${id}/move`, { album_id: bulkTargetAlbum.value }, authHeader())
      .then(() => successCount++)
      .catch(() => failCount++)
  )

  await Promise.allSettled(movePromises)
  isProcessing.value = false
  alert(`批量移动完成！\n✅ 成功: ${successCount}\n❌ 失败: ${failCount}`)
  selectedIds.value = []
  loadImages()
}

// ---------------- 右键菜单核心 ----------------
const handleContextMenu = (e, img) => {
  if (isMultiSelectMode.value) return 
  
  contextMenu.value = {
    show: true,
    x: e.clientX,
    y: e.clientY,
    img: img
  }
  contextMenuTargetAlbum.value = '' 
}

const closeContextMenu = () => {
  if (contextMenu.value.show) {
    contextMenu.value.show = false
  }
}

const deleteSingle = async (img) => {
  closeContextMenu()
  if (!confirm('确定要删除这张图片吗？')) return
  try {
    await axios.delete(`/api/images/${img.id}`, authHeader())
    loadImages()
  } catch (err) { alert('删除失败') }
}

const moveSingle = async (img, targetAlbumId) => {
  closeContextMenu()
  if (!targetAlbumId) return
  try {
    await axios.put(`/api/images/${img.id}/move`, { album_id: targetAlbumId }, authHeader())
    loadImages()
  } catch (err) { alert('移动失败') }
}

const renameSingle = async (img) => {
  closeContextMenu()
  const newName = prompt('📝 请输入新的文件名:', img.filename || 'new_name')
  if (!newName) return
  try {
    await axios.put(`/api/images/${img.id}/rename`, { filename: newName }, authHeader())
    loadImages()
  } catch (err) { 
    alert(err.response?.data?.error || '重命名失败 (可能需更新后端API支持)') 
  }
}

// ---------------- 辅助工具 ----------------
const formatSize = (bytes) => {
  if (bytes === 0 || !bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

const copy = text => navigator.clipboard.writeText(text).then(() => alert('🎉 复制成功！'))

onMounted(() => {
  loadAlbums()
  loadImages()
  window.addEventListener('click', closeContextMenu)
  window.addEventListener('scroll', closeContextMenu, true)
  
  window.addEventListener('keydown', (e) => {
    if (e.key === 'Escape' && previewImageUrl.value) {
      closePreview()
    }
  })
})

onUnmounted(() => {
  window.removeEventListener('click', closeContextMenu)
  window.removeEventListener('scroll', closeContextMenu, true)
})
</script>

<style scoped>
.gallery-panel { width: 100%; display: flex; flex-direction: column; gap: 20px; padding-bottom: 20px; position: relative;}

/* 顶部工具栏 */
.toolbar { display: flex; justify-content: space-between; align-items: center; background: white; padding: 15px 20px; border-radius: 12px; box-shadow: 0 2px 10px rgba(0,0,0,0.05); border: 1px solid #f0f0f0;}
.left-tools { display: flex; align-items: center; gap: 15px; flex-wrap: wrap;}
.title { margin: 0; font-size: 20px; color: #2d3748;}
.right-tools { display: flex; align-items: center; gap: 10px; flex-wrap: wrap;}

.form-select { padding: 8px 12px; border-radius: 6px; border: 1px solid #cbd5e0; background: #f8fafc; font-weight: bold; color: #4a5568; outline: none;}
.form-select:focus { border-color: #4CAF50;}
.small-select { padding: 6px 10px; font-size: 13px;}

/* 动态操作栏 */
.bulk-action-bar { display: flex; align-items: center; gap: 8px; background: #f0fdf4; padding: 6px 15px; border-radius: 20px; border: 1px solid #c6f6d5;}
.selected-text { font-size: 13px; color: #276749;}
.divider { width: 1px; height: 20px; background: #cbd5e0; margin: 0 5px;}
.fade-in { animation: fadeIn 0.3s ease-in-out; }
@keyframes fadeIn { from { opacity: 0; transform: translateX(-10px); } to { opacity: 1; transform: translateX(0); } }

/* 按钮样式 */
.btn { padding: 8px 16px; border: none; border-radius: 6px; cursor: pointer; transition: 0.2s; font-weight: bold;}
.small-btn { padding: 6px 12px; font-size: 13px; }
.primary-btn { background: #4CAF50; color: white; }
.primary-btn:hover { background: #45a049; }
.primary-btn:disabled { background: #a0aec0; cursor: not-allowed; }
.outline-btn { background: white; border: 1px solid #cbd5e0; color: #4a5568; }
.outline-btn:hover { background: #edf2f7; }
.danger-btn { background: #e53e3e; color: white; }
.danger-btn:hover { background: #c53030; }
.danger-btn:disabled { background: #fc8181; cursor: not-allowed; }

/* 缺省页 */
.empty-state { text-align: center; padding: 50px; color: #a0aec0; background: #f8fafc; border-radius: 12px; border: 2px dashed #e2e8f0; font-size: 16px; font-weight: bold;}

/* 🌟 核心修复：三大布局流派归位 */
.image-grid { display: grid; gap: 20px; }

/* 1. 大图模式：重塑霸气尺寸，最小320px */
.image-grid.large { grid-template-columns: repeat(auto-fill, minmax(600px, 1fr)); }
.image-grid.large .img-wrapper { aspect-ratio: 16 / 9; }

/* 2. 小图模式：正方形紧凑排版，最小160px */
.image-grid.small { grid-template-columns: repeat(auto-fill, minmax(320px, 1fr)); }
.image-grid.small .img-wrapper { aspect-ratio: 1 / 1; }

/* 3. 列表模式：抛弃网格，回归纯粹的 Flex 单行左右流式布局 */
.image-grid.list { display: flex; flex-direction: column; gap: 10px; }
.image-grid.list .img-card { flex-direction: row; align-items: center; justify-content: space-between; padding: 10px 15px; height: auto; }
.image-grid.list .checkbox-wrap { position: static; background: transparent; box-shadow: none; margin-right: 15px; padding: 0; }
.image-grid.list .img-wrapper { width: 80px; height: 80px; aspect-ratio: auto; flex-shrink: 0; border-radius: 6px; }
.image-grid.list .info-bar { border-top: none; background: transparent; flex: 1; text-align: left; padding: 0 20px; display: flex; flex-direction: column; gap: 5px; }
.image-grid.list .filename-text { font-size: 14px; font-weight: bold; color: #2d3748; }
.image-grid.list .actions { padding: 0; background: transparent; flex-shrink: 0; gap: 8px; }

/* 图片卡片通用样式 */
.img-card { position: relative; background: white; border-radius: 10px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.08); border: 2px solid transparent; transition: all 0.2s; display: flex; flex-direction: column; cursor: pointer;}
.img-card:hover { transform: translateY(-3px); box-shadow: 0 6px 15px rgba(0,0,0,0.12); }
.img-card.is-selected { border-color: #e53e3e; background: #fff5f5; }

/* 悬浮复选框 (非列表模式) */
.image-grid:not(.list) .checkbox-wrap { position: absolute; top: 10px; left: 10px; z-index: 10; background: rgba(255,255,255,0.9); border-radius: 4px; padding: 4px; display: flex; align-items: center; justify-content: center; box-shadow: 0 2px 4px rgba(0,0,0,0.1);}
.checkbox-wrap input[type="checkbox"] { width: 18px; height: 18px; cursor: pointer; accent-color: #e53e3e;}

/* 缩略图通用 */
.img-wrapper { width: 100%; background: #f0f0f0; overflow: hidden; display: flex; align-items: center; justify-content: center;}
.thumb { width: 100%; height: 100%; object-fit: contain; transition: 0.3s;} /* 🌟 修复：从 cover 改为 contain，保留原比例！ */
.img-card:hover .thumb { transform: scale(1.05); }

/* 信息栏与操作按钮通用 */
.info-bar { padding: 10px 12px; font-size: 13px; color: #718096; background: #f8fafc; border-top: 1px solid #edf2f7; text-align: center;}
.actions { display: flex; flex-wrap: wrap; gap: 5px; padding: 12px; justify-content: center; background: white;}
.actions button { padding: 5px 8px; font-size: 12px; font-weight: bold; border: 1px solid #e2e8f0; background: white; border-radius: 4px; cursor: pointer; color: #4a5568; transition: 0.1s;}
.actions button:hover { background: #edf2f7; color: #4CAF50; border-color: #cbd5e0;}

/* 豪华分页器样式 */
.pagination-bar { display: flex; justify-content: center; align-items: center; gap: 10px; padding: 20px; background: white; border-radius: 12px; box-shadow: 0 2px 10px rgba(0,0,0,0.05); flex-wrap: wrap;}
.page-btn { padding: 6px 12px; border: 1px solid #cbd5e0; background: white; color: #4a5568; border-radius: 6px; cursor: pointer; font-weight: bold; transition: 0.2s;}
.page-btn:hover:not(:disabled) { background: #edf2f7; border-color: #a0aec0;}
.page-btn:disabled { color: #cbd5e0; cursor: not-allowed; background: #f8fafc;}

.page-numbers { display: flex; gap: 5px;}
.num-btn { width: 32px; height: 32px; border: 1px solid #cbd5e0; background: white; border-radius: 6px; cursor: pointer; font-weight: bold; color: #4a5568; transition: 0.2s;}
.num-btn:hover { background: #edf2f7;}
.num-btn.active { background: #4CAF50; color: white; border-color: #4CAF50;}

.jump-box { display: flex; align-items: center; gap: 8px; font-size: 14px; color: #4a5568; margin-left: 10px;}
.jump-input { width: 50px; padding: 4px 8px; border: 1px solid #cbd5e0; border-radius: 4px; text-align: center; outline: none;}
.jump-input:focus { border-color: #4CAF50;}
.jump-btn { padding: 4px 10px; margin-left: 5px;}

/* 右键菜单样式 */
.custom-context-menu { position: fixed; z-index: 1000; background: white; border-radius: 8px; box-shadow: 0 4px 20px rgba(0,0,0,0.15); border: 1px solid #e2e8f0; padding: 5px 0; min-width: 180px; animation: menuFadeIn 0.2s ease-out;}
@keyframes menuFadeIn { from { opacity: 0; transform: scale(0.95); } to { opacity: 1; transform: scale(1); } }
.menu-header { padding: 8px 15px; font-size: 12px; color: #a0aec0; border-bottom: 1px solid #edf2f7; margin-bottom: 5px;}
.menu-item { padding: 10px 15px; font-size: 14px; color: #2d3748; cursor: pointer; transition: 0.2s; display: flex; align-items: center; gap: 8px;}
.menu-item:hover { background: #edf2f7; color: #4CAF50;}
.menu-item.danger { color: #e53e3e;}
.menu-item.danger:hover { background: #fff5f5;}
.menu-separator { height: 1px; background: #edf2f7; margin: 5px 0;}
.move-item { flex-direction: column; align-items: flex-start; gap: 5px;}
.context-select { width: 100%; padding: 4px; font-size: 12px; border: 1px solid #cbd5e0; border-radius: 4px; margin-top: 2px; outline: none;}

/* 全屏大图预览灯箱样式 */
.preview-mask { position: fixed; top: 0; left: 0; width: 100vw; height: 100vh; background: rgba(0, 0, 0, 0.85); z-index: 2000; display: flex; justify-content: center; align-items: center; backdrop-filter: blur(5px);}
.preview-img { max-width: 90vw; max-height: 90vh; object-fit: contain; border-radius: 8px; box-shadow: 0 10px 30px rgba(0,0,0,0.5); animation: zoomIn 0.3s cubic-bezier(0.18, 0.89, 0.32, 1.28);}
.close-preview-btn { position: absolute; top: 20px; right: 30px; background: rgba(255, 255, 255, 0.1); border: 2px solid rgba(255, 255, 255, 0.3); color: white; font-size: 20px; width: 44px; height: 44px; border-radius: 50%; cursor: pointer; transition: all 0.2s; display: flex; justify-content: center; align-items: center;}
.close-preview-btn:hover { background: rgba(255, 255, 255, 0.3); transform: scale(1.1); border-color: white;}
@keyframes zoomIn { from { transform: scale(0.8); opacity: 0; } to { transform: scale(1); opacity: 1; } }
</style>