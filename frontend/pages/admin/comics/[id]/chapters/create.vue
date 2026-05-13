<template>
  <div class="create-chapter-page">
    <div class="page-header">
      <button class="back-link" @click="router.back()">← Quay lại truyện</button>
      <h2>Thêm Chương Mới</h2>
    </div>

    <div class="create-grid">
      <!-- Left: Form -->
      <div class="form-panel">
        <!-- Chapter Info -->
        <div class="form-section">
          <h3>📋 Thông tin Chương</h3>
          <div class="form-group">
            <label>Số chương <span class="req">*</span></label>
            <input type="number" v-model="form.number" placeholder="VD: 1" min="0" step="0.5" />
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
          <p class="help-text">Ảnh sẽ tự động chuyển sang WebP để tối ưu dung lượng trước khi lưu.</p>


          <div class="dropzone" @dragover.prevent @drop.prevent="handleDrop" @click="fileInput?.click()">
            <input ref="fileInput" type="file" multiple hidden accept="image/*" @change="handleLocalFiles" />
            <div class="drop-info">
              <span class="drop-icon">📁</span>
              <p>Kéo thả ảnh vào đây hoặc nhấp để chọn</p>
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
            <div class="drive-logo">
              <svg viewBox="0 0 87.3 78" width="48"><path d="M6.6 66.85l3.85 6.65c.8 1.4 1.95 2.5 3.3 3.3l13.75-23.8H0c0 1.55.4 3.1 1.2 4.5z" fill="#0066da"/><path d="M43.65 25L29.9 1.2C28.55 2 27.4 3.1 26.6 4.5L1.2 48.4c-.8 1.4-1.2 2.95-1.2 4.5h27.5z" fill="#00ac47"/><path d="M73.55 76.8c1.35-.8 2.5-1.9 3.3-3.3l1.6-2.75 7.65-13.25c.8-1.4 1.2-2.95 1.2-4.5H59.8l5.85 11.5z" fill="#ea4335"/><path d="M43.65 25L57.4 1.2C56.05.4 54.5 0 52.9 0H34.4c-1.6 0-3.15.45-4.5 1.2z" fill="#00832d"/><path d="M59.8 52.9H27.5L13.75 76.7c1.35.8 2.9 1.3 4.5 1.3h50.8c1.6 0 3.15-.5 4.5-1.3z" fill="#2684fc"/><path d="M73.4 26.5l-12.7-22c-.8-1.4-1.95-2.5-3.3-3.3L43.65 25 59.8 52.9h27.45c0-1.55-.4-3.1-1.2-4.5z" fill="#ffba00"/></svg>
            </div>
            <p>Kết nối tài khoản Google để chọn ảnh từ Drive</p>
            <button class="btn-google" @click="initGoogleAuth">
              <svg width="18" height="18" viewBox="0 0 24 24"><path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/><path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/><path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l3.66-2.84z"/><path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/></svg>
              Đăng nhập với Google
            </button>
            <p class="drive-hint">Cần nhập <strong>Google Client ID</strong> trong file <code>.env</code>:<br><code>NUXT_PUBLIC_GOOGLE_CLIENT_ID=your_client_id</code></p>
          </div>

          <div v-else class="drive-logged-in">
            <div class="drive-user-info">
              <span class="drive-connected">✅ Đã kết nối Google Drive</span>
              <button class="btn-text" @click="driveToken = null">Đăng xuất</button>
            </div>
            <button class="btn-drive-pick" @click="openDrivePicker">
              📂 Chọn ảnh từ Google Drive
            </button>
            <p class="help-text">Sau khi chọn, danh sách URL ảnh Drive sẽ được lưu vào server này.</p>
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
              <div class="grid-size">{{ (img.blob.size / 1024).toFixed(0) }}KB</div>
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
          <p class="help-text" style="margin-bottom: 12px;">Thiết lập ảnh cho từng server. Bạn có thể nạp cả hai.</p>
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
          <div class="summary-row" style="margin-top: 8px; border-top: 1px solid rgba(255,255,255,0.1); padding-top: 12px;">
            <span>Tình trạng:</span>
            <strong v-if="images.length || driveImages.length" style="color: #4ade80;">Sẵn sàng đăng</strong>
            <strong v-else style="color: #ef4444;">Chưa có ảnh</strong>
          </div>
        </div>

        <button class="btn-publish" :disabled="isSubmitting" @click="handleSubmit">
          {{ isSubmitting ? '⏳ Đang đăng...' : '🚀 Đăng Chương' }}
        </button>
      </aside>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'

definePageMeta({ layout: 'admin', title: 'Thêm Chương Mới' })

const router = useRouter()
const route = useRoute()
const { get, post } = useApi()
const authStore = useAuthStore()

const fileInput = ref<HTMLInputElement | null>(null)
const isSubmitting = ref(false)
const converting = ref(false)
const convertProgress = ref(0)
const driveToken = ref<string | null>(null)
const driveImages = ref<{ name: string; url: string }[]>([])

onMounted(async () => {
  // Auto-calculate next chapter number
  try {
    const storyId = route.params.id
    const existingChapters = await get<any[]>(`/admin/stories/${storyId}/chapters`, {
      headers: { Authorization: `Bearer ${authStore.token}` }
    })
    
    if (existingChapters && existingChapters.length > 0) {
      const maxNum = Math.max(...existingChapters.map(c => c.number || 0))
      form.number = maxNum + 1
      form.title = `Chương ${maxNum + 1}`
    } else {
      form.number = 1
      form.title = `Chương 1`
    }
  } catch (err) {
    console.error('Error fetching existing chapters:', err)
    form.number = 1
    form.title = `Chương 1`
  }
})

const form = reactive({
  number: null as number | null,
  title: '',
  note: '',
  activeTab: 'server1'
})

interface WebPImage {
  preview: string
  blob: Blob
  name: string
}
const images = ref<WebPImage[]>([])

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

// ─── Lightbox ─────────────────────────────────────────────────────────────
const lightboxIndex = ref(-1)
const openLightbox = (i: number) => { lightboxIndex.value = i }

// ─── Drag & Drop Sort ─────────────────────────────────────────────────────
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

const processFiles = async (files: File[]) => {
  const remaining = 20 - images.value.length
  if (remaining <= 0) return alert('Đã đạt tối đa 20 ảnh!')
  const toProcess = files.slice(0, remaining)

  converting.value = true
  convertProgress.value = 0

  for (let i = 0; i < toProcess.length; i++) {
    const result = await convertToWebP(toProcess[i])
    images.value.push(result)
    convertProgress.value = Math.round(((i + 1) / toProcess.length) * 100)
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

// ─── Google Drive Picker ───────────────────────────────────────────────────
const GOOGLE_CLIENT_ID = useRuntimeConfig().public.googleClientId as string

const initGoogleAuth = () => {
  if (!GOOGLE_CLIENT_ID) {
    alert('Chưa cấu hình NUXT_PUBLIC_GOOGLE_CLIENT_ID trong .env!')
    return
  }
  // Load Google Identity Services
  const script = document.createElement('script')
  script.src = 'https://accounts.google.com/gsi/client'
  script.onload = () => {
    (window as any).google.accounts.oauth2.initTokenClient({
      client_id: GOOGLE_CLIENT_ID,
      scope: 'https://www.googleapis.com/auth/drive.readonly',
      callback: (res: any) => {
        if (res.access_token) driveToken.value = res.access_token
      }
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
        .addView(new (window as any).google.picker.ImageSearchView())
        .addView(new (window as any).google.picker.DocsView()
          .setIncludeFolders(true)
          .setMimeTypes('image/jpeg,image/png,image/webp')
          .setMode((window as any).google.picker.DocsViewMode.GRID))
        .setOAuthToken(driveToken.value)
        .setCallback((data: any) => {
          if (data.action === 'picked') {
            data.docs.forEach((doc: any) => {
              driveImages.value.push({
                name: doc.name,
                url: `https://drive.google.com/uc?id=${doc.id}&export=view`
              })
            })
          }
        })
        .build()
      picker.setVisible(true)
    })
  }
  document.head.appendChild(script)
}

// ─── Submit ────────────────────────────────────────────────────────────────
// ─── Submit ────────────────────────────────────────────────────────────────
const handleSubmit = async () => {
  if (!form.number) return alert('Vui lòng nhập số chương!')
  if (images.value.length === 0 && driveImages.value.length === 0) {
    return alert('Vui lòng cung cấp ảnh cho ít nhất 1 server!')
  }

  isSubmitting.value = true
  try {
    const serversToSend = []

    // 1. Process Server 1 (R2)
    if (images.value.length > 0) {
      const formData = new FormData()
      for (let i = 0; i < images.value.length; i++) {
        formData.append('images', images.value[i].blob, images.value[i].name)
      }
      formData.append('story_id', String(route.params.id))
      formData.append('chapter_number', String(form.number))

      let content = ''
      try {
        const config = useRuntimeConfig()
        const uploadRes = await fetch(`${config.public.apiBase}/admin/upload-images`, {
          method: 'POST',
          headers: { Authorization: `Bearer ${authStore.token}` },
          body: formData
        })
        if (uploadRes.ok) {
          const data = await uploadRes.json()
          content = JSON.stringify(data.urls)
        } else {
          content = JSON.stringify(images.value.map(img => img.name))
        }
      } catch {
        content = JSON.stringify(images.value.map(img => img.name))
      }

      serversToSend.push({
        name: 'Server 1 (R2)',
        type: 'image',
        content
      })
    }

    // 2. Process Server 2 (Drive)
    if (driveImages.value.length > 0) {
      serversToSend.push({
        name: 'Server 2 (Drive)',
        type: 'image',
        content: JSON.stringify(driveImages.value.map(i => i.url))
      })
    }

    await post('/admin/chapters', {
      story_id: Number(route.params.id),
      number: Number(form.number),
      title: form.title || `Chương ${form.number}`,
      note: form.note,
      servers: serversToSend
    }, { headers: { Authorization: `Bearer ${authStore.token}` } })

    alert(`✅ Đã đăng Chương ${form.number} thành công với ${serversToSend.length} server!`)
    router.push(`/admin/comics/${route.params.id}`)
  } catch (err: any) {
    console.error('Chapter create error:', err)
    alert(`Lỗi: ${err?.data?.error || err?.message || 'Không thể đăng chương'}`)
  } finally {
    isSubmitting.value = false
  }
}

</script>

<style scoped>
.create-chapter-page { display: flex; flex-direction: column; gap: 24px; padding-bottom: 40px; }
.page-header { display: flex; align-items: center; gap: 16px; }
.back-link { background: none; border: none; color: #9CA764; cursor: pointer; font-weight: 600; font-size: 0.95rem; padding: 0; }
.page-header h2 { margin: 0; color: #fff; font-family: 'Montserrat', sans-serif; }

.create-grid { display: grid; grid-template-columns: 1fr 300px; gap: 24px; align-items: start; }
@media(max-width:1024px) { .create-grid { grid-template-columns: 1fr; } }

.form-panel { display: flex; flex-direction: column; gap: 24px; }
.form-section { background: #141720; border: 1px solid rgba(255,255,255,0.05); border-radius: 16px; padding: 24px; display: flex; flex-direction: column; gap: 16px; }
.form-section h3 { margin: 0; color: #F1E8C7; font-size: 1rem; padding-bottom: 12px; border-bottom: 1px solid rgba(255,255,255,0.05); }
.form-group { display: flex; flex-direction: column; gap: 8px; }
.form-group label { font-size: 0.85rem; color: #A8A8B3; font-weight: 600; }
.req { color: #ef4444; }
.form-group input, .form-group textarea, .form-group select { background: rgba(0,0,0,0.2); border: 1px solid rgba(255,255,255,0.1); border-radius: 8px; padding: 11px 14px; color: #fff; font-size: 0.9rem; transition: border-color 0.2s; width: 100%; box-sizing: border-box; }
.form-group input:focus, .form-group textarea:focus { border-color: #9CA764; outline: none; }
.help-text { font-size: 0.8rem; color: #5C5C6B; margin: 0; }

/* Dropzone */
.dropzone { border: 2px dashed rgba(255,255,255,0.1); border-radius: 12px; padding: 40px; text-align: center; cursor: pointer; transition: all 0.2s; background: rgba(0,0,0,0.1); }
.dropzone:hover { border-color: #9CA764; background: rgba(156,167,100,0.04); }
.drop-icon { font-size: 2.5rem; display: block; margin-bottom: 10px; }
.drop-info p { color: #fff; font-weight: 600; margin: 0 0 6px; }
.drop-info small { color: #5C5C6B; font-size: 0.8rem; }

/* Converting progress */
.converting-bar { position: relative; background: rgba(0,0,0,0.2); border-radius: 8px; height: 36px; overflow: hidden; display: flex; align-items: center; justify-content: center; }
.converting-inner { position: absolute; left: 0; top: 0; height: 100%; background: linear-gradient(90deg, #9CA764, #b0c070); transition: width 0.2s; border-radius: 8px; }
.converting-bar span { position: relative; z-index: 1; font-size: 0.85rem; font-weight: 700; color: #000; }

/* Image grid */
.image-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(140px, 1fr)); gap: 12px; }
.grid-img-item { position: relative; border-radius: 10px; overflow: hidden; background: #000; cursor: grab; border: 2px solid rgba(255,255,255,0.06); transition: border-color 0.2s, transform 0.15s; aspect-ratio: 3/4; }
.grid-img-item:active { cursor: grabbing; }
.grid-img-item.drag-over { border-color: #9CA764; box-shadow: 0 0 0 3px rgba(156,167,100,0.3); transform: scale(0.97); }
.grid-thumb { width: 100%; height: 100%; object-fit: cover; cursor: zoom-in; transition: transform 0.2s; display: block; }
.grid-thumb:hover { transform: scale(1.04); }
.grid-num { position: absolute; top: 6px; left: 6px; background: #9CA764; color: #000; font-size: 0.7rem; font-weight: 800; padding: 2px 6px; border-radius: 4px; }
.grid-remove { position: absolute; top: 6px; right: 6px; background: rgba(239,68,68,0.85); color: #fff; border: none; width: 24px; height: 24px; border-radius: 50%; cursor: pointer; font-size: 0.75rem; display: flex; align-items: center; justify-content: center; opacity: 0; transition: opacity 0.2s; }
.grid-img-item:hover .grid-remove { opacity: 1; }
.grid-size { position: absolute; bottom: 0; left: 0; right: 0; background: rgba(0,0,0,0.6); color: #9CA764; font-size: 0.7rem; text-align: center; padding: 4px; }

/* Lightbox */
.lightbox { position: fixed; inset: 0; z-index: 9999; background: rgba(0,0,0,0.92); display: flex; align-items: center; justify-content: center; }
.lb-img { max-width: 90vw; max-height: 90vh; object-fit: contain; border-radius: 8px; box-shadow: 0 8px 60px rgba(0,0,0,0.8); }
.lb-close { position: absolute; top: 20px; right: 24px; background: rgba(255,255,255,0.1); border: none; color: #fff; width: 40px; height: 40px; border-radius: 50%; font-size: 1.2rem; cursor: pointer; transition: background 0.2s; }
.lb-close:hover { background: rgba(239,68,68,0.7); }
.lb-prev, .lb-next { position: absolute; top: 50%; transform: translateY(-50%); background: rgba(255,255,255,0.1); border: none; color: #fff; width: 48px; height: 48px; border-radius: 50%; font-size: 1.8rem; cursor: pointer; transition: background 0.2s; display: flex; align-items: center; justify-content: center; }
.lb-prev { left: 20px; }
.lb-next { right: 20px; }
.lb-prev:hover, .lb-next:hover { background: rgba(156,167,100,0.5); }
.lb-info { position: absolute; bottom: 20px; left: 50%; transform: translateX(-50%); background: rgba(0,0,0,0.6); color: #fff; font-size: 0.85rem; padding: 6px 16px; border-radius: 20px; white-space: nowrap; max-width: 80vw; overflow: hidden; text-overflow: ellipsis; }

/* Google Drive */
.drive-login-box { display: flex; flex-direction: column; align-items: center; gap: 16px; padding: 24px; border: 1px dashed rgba(255,255,255,0.1); border-radius: 12px; text-align: center; }
.drive-logo { display: flex; align-items: center; justify-content: center; }
.drive-login-box p { color: #A8A8B3; margin: 0; }
.btn-google { display: flex; align-items: center; gap: 10px; background: #fff; color: #333; border: none; padding: 10px 20px; border-radius: 8px; font-weight: 700; cursor: pointer; font-size: 0.95rem; transition: box-shadow 0.2s; }
.btn-google:hover { box-shadow: 0 2px 12px rgba(255,255,255,0.2); }
.drive-hint { font-size: 0.75rem; color: #5C5C6B; line-height: 1.6; }
.drive-hint code { background: rgba(255,255,255,0.05); padding: 2px 6px; border-radius: 4px; font-size: 0.7rem; }
.drive-logged-in { display: flex; flex-direction: column; gap: 12px; }
.drive-user-info { display: flex; align-items: center; justify-content: space-between; }
.drive-connected { color: #4ade80; font-size: 0.9rem; font-weight: 600; }
.btn-text { background: none; border: none; color: #A8A8B3; cursor: pointer; font-size: 0.85rem; text-decoration: underline; }
.btn-drive-pick { background: #4285F4; color: #fff; border: none; padding: 12px 20px; border-radius: 8px; font-weight: 700; cursor: pointer; transition: all 0.2s; font-size: 0.95rem; }
.btn-drive-pick:hover { background: #3367d6; transform: translateY(-1px); }
.drive-image-list { display: flex; flex-direction: column; gap: 8px; margin-top: 4px; }
.drive-img-item { display: flex; align-items: center; gap: 10px; background: rgba(0,0,0,0.2); border-radius: 8px; padding: 10px 14px; }
.drive-img-item .img-name { flex: 1; font-size: 0.85rem; color: #fff; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

/* Server options */
.server-options { display: flex; flex-direction: column; gap: 10px; }
.server-option { display: flex; align-items: center; gap: 14px; background: rgba(0,0,0,0.2); border: 2px solid rgba(255,255,255,0.08); border-radius: 10px; padding: 14px; cursor: pointer; transition: all 0.2s; }
.server-option input { display: none; }
.server-option.active { border-color: #9CA764; background: rgba(156,167,100,0.08); }
.server-icon { font-size: 1.5rem; }
.server-option strong { display: block; color: #fff; font-size: 0.95rem; }
.server-option small { color: #5C5C6B; font-size: 0.8rem; }
.status-badge { background: #4ade80; color: #000; font-size: 0.7rem; font-weight: 800; padding: 2px 8px; border-radius: 10px; margin-left: 8px; }

/* Summary */
.summary-row { display: flex; justify-content: space-between; align-items: center; padding: 8px 0; border-bottom: 1px solid rgba(255,255,255,0.04); font-size: 0.9rem; }
.summary-row span { color: #5C5C6B; }
.summary-row strong { color: #F1E8C7; }

/* Sidebar */
.sidebar { display: flex; flex-direction: column; gap: 24px; position: sticky; top: 24px; }
.btn-publish { background: #9CA764; color: #000; border: none; padding: 16px; border-radius: 12px; font-weight: 800; font-size: 1rem; cursor: pointer; transition: all 0.2s; width: 100%; }
.btn-publish:hover:not(:disabled) { background: #b0c070; transform: translateY(-2px); box-shadow: 0 6px 20px rgba(156,167,100,0.4); }
.btn-publish:disabled { opacity: 0.6; cursor: not-allowed; }
</style>
