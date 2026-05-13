<template>
  <div class="create-chapter-page">
    <div class="page-header">
      <button class="back-link" @click="router.back()">← Quay lại truyện</button>
      <h2>Chỉnh sửa Chương {{ form.number }}</h2>
    </div>

    <div class="create-grid" v-if="!loading">
      <!-- Left: Form -->
      <div class="form-panel">
        <div class="form-section">
          <h3>📋 Thông tin Chương</h3>
          <div class="form-group">
            <label>Số chương <span class="req">*</span></label>
            <input type="number" v-model="form.number" placeholder="VD: 1" min="0" step="0.5" disabled />
          </div>
          <div class="form-group">
            <label>Tên chương</label>
            <input type="text" v-model="form.title" placeholder="VD: Chương 1: Bắt đầu hành trình" />
          </div>
          <div class="form-group">
            <label>Ghi chú nội bộ</label>
            <textarea v-model="form.note" rows="2" placeholder="Ghi chú không hiển thị cho người đọc..."></textarea>
          </div>
        </div>

        <!-- Server 1: Local Upload -->
        <div class="form-section" v-show="form.activeTab === 'server1'">
          <h3>☁️ Server 1 — Tải ảnh từ máy tính (Cloudflare R2)</h3>
          <p class="help-text">Nếu bạn tải ảnh mới lên, danh sách ảnh cũ của Server 1 sẽ bị thay thế.</p>

          <div class="dropzone" @dragover.prevent @drop.prevent="handleDrop" @click="fileInput?.click()">
            <input ref="fileInput" type="file" multiple hidden accept="image/*" @change="handleLocalFiles" />
            <div class="drop-info">
              <span class="drop-icon">📁</span>
              <p>Kéo thả ảnh vào đây hoặc nhấp để <strong>thay thế</strong> ảnh cũ</p>
              <small>Hỗ trợ JPG, PNG, GIF • Tối đa 20 ảnh • Tự động chuyển WebP</small>
            </div>
          </div>

          <!-- Converting indicator -->
          <div v-if="converting" class="converting-bar">
            <div class="converting-inner" :style="{ width: convertProgress + '%' }"></div>
            <span>Đang chuyển đổi WebP... {{ convertProgress }}%</span>
          </div>
        </div>

        <!-- Server 2: Google Drive -->
        <div class="form-section" v-show="form.activeTab === 'server2'">
          <h3>🚗 Server 2 — Google Drive</h3>
          <p class="help-text">Đăng nhập Google để chọn thư mục/ảnh từ Drive của bạn.</p>

          <div v-if="!driveToken" class="drive-login-box">
            <button class="btn-google" @click="initGoogleAuth">Đăng nhập với Google</button>
          </div>

          <div v-else class="drive-logged-in">
            <div class="drive-user-info">
              <span class="drive-connected">✅ Đã kết nối Google Drive</span>
              <button class="btn-text" @click="driveToken = null">Đăng xuất</button>
            </div>
            <button class="btn-drive-pick" @click="openDrivePicker">📂 Chọn ảnh từ Google Drive</button>
          </div>

          <!-- Drive Images list -->
          <div v-if="driveImages.length" class="drive-image-list">
            <div v-for="(img, i) in driveImages" :key="i" class="drive-img-item">
              <span class="order-num">{{ i + 1 }}</span>
              <span class="img-name">{{ img.name }}</span>
              <button class="remove-btn" @click="driveImages.splice(i, 1)">✕</button>
            </div>
          </div>
        </div>

        <!-- Image grid (Server 1) -->
        <div v-show="form.activeTab === 'server1' && images.length > 0" class="form-section">
          <h3>🖼️ Danh sách ảnh Server 1 ({{ images.length }}/20)</h3>
          <div class="image-grid">
            <div
              v-for="(img, i) in images" :key="i"
              class="grid-img-item"
              :class="{ 'drag-over': dragOverIndex === i }"
              draggable="true"
              @dragstart="onDragStart(i)"
              @dragover.prevent="onDragOver(i)"
              @drop.prevent="onDrop(i)"
              @dragend="dragOverIndex = -1"
            >
              <img :src="img.preview" class="grid-thumb" @click="openLightbox(i)" />
              <span class="grid-num">{{ i + 1 }}</span>
              <button class="grid-remove" @click.stop="images.splice(i, 1)">✕</button>
              <div v-if="img.blob" class="grid-size">{{ (img.blob.size / 1024).toFixed(0) }}KB</div>
              <div v-else class="grid-size">Đã lưu</div>
            </div>
          </div>
        </div>

        <!-- Lightbox -->
        <Teleport to="body">
          <div v-if="lightboxIndex >= 0" class="lightbox" @click.self="lightboxIndex = -1">
            <button class="lb-close" @click="lightboxIndex = -1">✕</button>
            <button class="lb-prev" @click="lightboxIndex = Math.max(0, lightboxIndex - 1)">‹</button>
            <img :src="images[lightboxIndex]?.preview" class="lb-img" />
            <button class="lb-next" @click="lightboxIndex = Math.min(images.length - 1, lightboxIndex + 1)">›</button>
            <div class="lb-info">{{ lightboxIndex + 1 }} / {{ images.length }} — {{ images[lightboxIndex]?.name }}</div>
          </div>
        </Teleport>
      </div>

      <!-- Right: Sidebar -->
      <aside class="sidebar">
        <div class="form-section">
          <h3>🖥️ Cấu hình Server</h3>
          <div class="server-options">
            <div class="server-option" :class="{ active: form.activeTab === 'server1' }" @click="form.activeTab = 'server1'">
              <div class="server-icon">☁️</div>
              <div style="flex: 1;">
                <strong>Server 1 (R2)</strong>
                <small>Upload từ máy tính</small>
              </div>
              <span v-if="images.length > 0" class="status-badge">✅ {{ images.length }}</span>
            </div>
            <div class="server-option" :class="{ active: form.activeTab === 'server2' }" @click="form.activeTab = 'server2'">
              <div class="server-icon">🚗</div>
              <div style="flex: 1;">
                <strong>Server 2 (Drive)</strong>
                <small>Chọn từ Google Drive</small>
              </div>
              <span v-if="driveImages.length > 0" class="status-badge">✅ {{ driveImages.length }}</span>
            </div>
          </div>
        </div>

        <div class="form-section">
          <h3>📊 Tổng kết</h3>
          <div class="summary-row"><span>S1 (Cloudflare):</span><strong>{{ images.length }} ảnh</strong></div>
          <div class="summary-row"><span>S2 (Google Drive):</span><strong>{{ driveImages.length }} ảnh</strong></div>
        </div>

        <button class="btn-publish" :disabled="isSubmitting" @click="handleSubmit">
          {{ isSubmitting ? '⏳ Đang lưu...' : '💾 Lưu Thay Đổi' }}
        </button>
      </aside>
    </div>
    <div v-else class="loading-state">
       <div class="spinner"></div>
       <p>Đang tải dữ liệu chương...</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'

definePageMeta({ layout: 'admin' })

const router = useRouter()
const route = useRoute()
const { get, put } = useApi()
const authStore = useAuthStore()

const loading = ref(true)
const isSubmitting = ref(false)
const converting = ref(false)
const convertProgress = ref(0)
const fileInput = ref<HTMLInputElement | null>(null)
const driveToken = ref<string | null>(null)
const driveImages = ref<{ name: string; url: string }[]>([])

const form = reactive({
  number: 0,
  title: '',
  note: '',
  activeTab: 'server1'
})

interface WebPImage {
  preview: string
  blob?: Blob
  name: string
  isExisting?: boolean
  url?: string
}
const images = ref<WebPImage[]>([])

// ─── Load Data ─────────────────────────────────────────────────────────────
onMounted(async () => {
  try {
    const chap = await get<any>(`/admin/chapters/${route.params.chapterId}`, {
      headers: { Authorization: `Bearer ${authStore.token}` }
    })
    form.number = chap.number
    form.title = chap.title
    form.note = chap.note

    // Parse servers
    const s1 = chap.servers?.find((s: any) => s.name.includes('Server 1'))
    if (s1) {
      const urls = JSON.parse(s1.content)
      images.value = urls.map((u: string) => ({ preview: u, isExisting: true, url: u, name: 'Ảnh cũ' }))
    }

    const s2 = chap.servers?.find((s: any) => s.name.includes('Server 2'))
    if (s2) {
      const urls = JSON.parse(s2.content)
      driveImages.value = urls.map((u: string) => ({ url: u, name: 'Drive Image' }))
    }
  } catch (err) {
    console.error(err)
    alert('Không thể tải dữ liệu chương!')
  } finally {
    loading.value = false
  }
})

// ─── WebP Converter ────────────────────────────────────────────────────────
const convertToWebP = (file: File): Promise<WebPImage> => {
  return new Promise((resolve) => {
    const img = new Image()
    const url = URL.createObjectURL(file)
    img.onload = () => {
      const canvas = document.createElement('canvas')
      canvas.width = img.width
      canvas.height = img.height
      canvas.getContext('2d')!.drawImage(img, 0, 0)
      canvas.toBlob((blob) => {
        URL.revokeObjectURL(url)
        const preview = URL.createObjectURL(blob!)
        const name = file.name.replace(/\.[^.]+$/, '.webp')
        resolve({ preview, blob: blob!, name })
      }, 'image/webp', 0.85)
    }
    img.src = url
  })
}

const processFiles = async (files: File[]) => {
  images.value = images.value.filter(i => i.isExisting) // Keep old if wanted, or clear?
  // User usually wants to REPLACE if they drop new files in edit mode
  images.value = [] 
  
  converting.value = true
  convertProgress.value = 0
  for (let i = 0; i < files.length; i++) {
    const result = await convertToWebP(files[i])
    images.value.push(result)
    convertProgress.value = Math.round(((i + 1) / files.length) * 100)
  }
  converting.value = false
}

const handleLocalFiles = (e: Event) => {
  const files = Array.from((e.target as HTMLInputElement).files || [])
  processFiles(files)
}

const handleDrop = (e: DragEvent) => {
  const files = Array.from(e.dataTransfer?.files || [])
  processFiles(files)
}

// ─── Lightbox & Drag ───────────────────────────────────────────────────────
const lightboxIndex = ref(-1)
const openLightbox = (i: number) => { lightboxIndex.value = i }
const dragSrcIndex = ref(-1)
const dragOverIndex = ref(-1)
const onDragStart = (i: number) => { dragSrcIndex.value = i }
const onDragOver = (i: number) => { dragOverIndex.value = i }
const onDrop = (targetIndex: number) => {
  const src = dragSrcIndex.value
  if (src === targetIndex || src < 0) return
  const arr = [...images.value]
  const [moved] = arr.splice(src, 1)
  arr.splice(targetIndex, 0, moved)
  images.value = arr
  dragSrcIndex.value = -1
  dragOverIndex.value = -1
}

// ─── Google Drive ──────────────────────────────────────────────────────────
const GOOGLE_CLIENT_ID = useRuntimeConfig().public.googleClientId as string
const initGoogleAuth = () => {
  const script = document.createElement('script')
  script.src = 'https://accounts.google.com/gsi/client'
  script.onload = () => {
    (window as any).google.accounts.oauth2.initTokenClient({
      client_id: GOOGLE_CLIENT_ID,
      scope: 'https://www.googleapis.com/auth/drive.readonly',
      callback: (res: any) => { if (res.access_token) driveToken.value = res.access_token }
    }).requestAccessToken()
  }
  document.head.appendChild(script)
}

const openDrivePicker = () => {
  if (!driveToken.value) return
  const script = document.createElement('script')
  script.src = 'https://apis.google.com/js/api.js'
  script.onload = () => {
    (window as any).gapi.load('picker', () => {
      const picker = new (window as any).google.picker.PickerBuilder()
        .addView(new (window as any).google.picker.DocsView().setMimeTypes('image/*'))
        .setOAuthToken(driveToken.value)
        .setCallback((data: any) => {
          if (data.action === 'picked') {
            data.docs.forEach((doc: any) => driveImages.value.push({ name: doc.name, url: `https://drive.google.com/uc?id=${doc.id}&export=view` }))
          }
        }).build()
      picker.setVisible(true)
    })
  }
  document.head.appendChild(script)
}

// ─── Submit ────────────────────────────────────────────────────────────────
const handleSubmit = async () => {
  isSubmitting.value = true
  try {
    const serversToSend = []
    
    // 1. Process Server 1
    let s1Content = ''
    const hasNewFiles = images.value.some(i => !i.isExisting)
    
    if (hasNewFiles) {
      const formData = new FormData()
      images.value.forEach(img => {
        if (img.blob) formData.append('images', img.blob, img.name)
      })
      formData.append('story_id', String(route.params.id))
      formData.append('chapter_number', String(form.number))

      const config = useRuntimeConfig()
      const res = await fetch(`${config.public.apiBase}/admin/upload-images`, {
        method: 'POST',
        headers: { Authorization: `Bearer ${authStore.token}` },
        body: formData
      })
      const data = await res.json()
      s1Content = JSON.stringify(data.urls)
    } else {
      s1Content = JSON.stringify(images.value.map(i => i.url))
    }

    if (images.value.length > 0) {
      serversToSend.push({ name: 'Server 1 (R2)', type: 'image', content: s1Content })
    }

    // 2. Server 2
    if (driveImages.value.length > 0) {
      serversToSend.push({ name: 'Server 2 (Drive)', type: 'image', content: JSON.stringify(driveImages.value.map(i => i.url)) })
    }

    await put(`/admin/chapters/${route.params.chapterId}`, {
      number: Number(form.number),
      title: form.title,
      note: form.note,
      servers: serversToSend
    }, { headers: { Authorization: `Bearer ${authStore.token}` } })

    alert('Cập nhật chương thành công!')
    router.push(`/admin/comics/${route.params.id}`)
  } catch (err: any) {
    console.error(err)
    alert('Lỗi: ' + (err.data?.error || 'Không thể lưu'))
  } finally {
    isSubmitting.value = false
  }
}
</script>

<style scoped>
.create-chapter-page { display: flex; flex-direction: column; gap: 24px; padding: 24px; padding-bottom: 40px; }
.page-header { display: flex; align-items: center; gap: 16px; }
.back-link { background: none; border: none; color: #9CA764; cursor: pointer; font-weight: 600; font-size: 0.95rem; padding: 0; }
.page-header h2 { margin: 0; color: #fff; font-family: 'Montserrat', sans-serif; }

.create-grid { display: grid; grid-template-columns: 1fr 320px; gap: 24px; align-items: start; }
@media(max-width:1024px) { .create-grid { grid-template-columns: 1fr; } }

.form-panel { display: flex; flex-direction: column; gap: 24px; }
.form-section { background: #141720; border: 1px solid rgba(255,255,255,0.05); border-radius: 16px; padding: 24px; display: flex; flex-direction: column; gap: 16px; }
.form-section h3 { margin: 0; color: #F1E8C7; font-size: 1rem; padding-bottom: 12px; border-bottom: 1px solid rgba(255,255,255,0.05); }

.form-group { display: flex; flex-direction: column; gap: 8px; }
.form-group label { font-size: 0.85rem; color: #A8A8B3; font-weight: 600; }
.form-group input, .form-group textarea { background: rgba(0,0,0,0.2); border: 1px solid rgba(255,255,255,0.1); border-radius: 8px; padding: 11px 14px; color: #fff; font-size: 0.9rem; width: 100%; box-sizing: border-box; }

.dropzone { border: 2px dashed rgba(255,255,255,0.1); border-radius: 12px; padding: 40px; text-align: center; cursor: pointer; background: rgba(0,0,0,0.1); transition: all 0.2s; }
.dropzone:hover { border-color: #9CA764; background: rgba(156,167,100,0.05); }
.drop-icon { font-size: 2.5rem; display: block; margin-bottom: 10px; }

.image-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(120px, 1fr)); gap: 12px; }
.grid-img-item { position: relative; border-radius: 10px; overflow: hidden; background: #000; border: 2px solid rgba(255,255,255,0.06); aspect-ratio: 3/4; }
.grid-thumb { width: 100%; height: 100%; object-fit: cover; cursor: zoom-in; }
.grid-num { position: absolute; top: 6px; left: 6px; background: #9CA764; color: #000; font-size: 0.7rem; font-weight: 800; padding: 2px 6px; border-radius: 4px; }
.grid-remove { position: absolute; top: 6px; right: 6px; background: rgba(239,68,68,0.85); color: #fff; border: none; width: 24px; height: 24px; border-radius: 50%; cursor: pointer; display: flex; align-items: center; justify-content: center; }
.grid-size { position: absolute; bottom: 0; left: 0; right: 0; background: rgba(0,0,0,0.6); color: #9CA764; font-size: 0.7rem; text-align: center; padding: 4px; }

.server-options { display: flex; flex-direction: column; gap: 10px; }
.server-option { display: flex; align-items: center; gap: 14px; background: rgba(0,0,0,0.2); border: 2px solid rgba(255,255,255,0.08); border-radius: 10px; padding: 14px; cursor: pointer; transition: all 0.2s; }
.server-option:hover { border-color: rgba(255,255,255,0.2); }
.server-option.active { border-color: #9CA764; background: rgba(156,167,100,0.08); }
.status-badge { background: #4ade80; color: #000; font-size: 0.7rem; font-weight: 800; padding: 2px 8px; border-radius: 10px; margin-left: auto; }

.btn-publish { background: #9CA764; color: #000; border: none; padding: 16px; border-radius: 12px; font-weight: 800; font-size: 1rem; cursor: pointer; width: 100%; margin-top: 10px; transition: all 0.2s; }
.btn-publish:hover { background: #b0c070; transform: translateY(-2px); }

.loading-state { display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 100px; color: #A8A8B3; gap: 16px; }
.spinner { width: 40px; height: 40px; border: 3px solid rgba(156,167,100,0.2); border-top-color: #9CA764; border-radius: 50%; animation: spin 0.8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }

.drive-logged-in { display: flex; flex-direction: column; gap: 12px; }
.drive-connected { color: #4ade80; font-size: 0.9rem; font-weight: 600; }
.btn-text { background: none; border: none; color: #A8A8B3; cursor: pointer; font-size: 0.85rem; text-decoration: underline; }
.btn-drive-pick { background: #4285F4; color: #fff; border: none; padding: 12px 20px; border-radius: 8px; font-weight: 700; cursor: pointer; }

.lightbox { position: fixed; inset: 0; z-index: 9999; background: rgba(0,0,0,0.92); display: flex; align-items: center; justify-content: center; }
.lb-img { max-width: 90vw; max-height: 90vh; object-fit: contain; }
.lb-close { position: absolute; top: 20px; right: 24px; background: rgba(255,255,255,0.1); border: none; color: #fff; width: 40px; height: 40px; border-radius: 50%; cursor: pointer; }
</style>
